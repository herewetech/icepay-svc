/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file client.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 02/25/2023
 */

package handler

import (
	"database/sql"
	"errors"
	"icepay-svc/handler/request"
	"icepay-svc/handler/response"
	"icepay-svc/model"
	"icepay-svc/runtime"
	"icepay-svc/service"
	"icepay-svc/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/skip2/go-qrcode"
)

type Client struct {
	svcAuth       *service.Auth
	svcCredential *service.Credential
}

func InitClient() *Client {
	h := new(Client)

	clientG := runtime.Server.Group("/client")
	clientG.Post("/token", h.token).Name("ClientPostToken")
	clientG.Post("/refresh", h.refresh).Name("ClientPostRefresh")
	clientG.Use(jwtware.New(jwtware.Config{
		SigningKey:     []byte(runtime.Config.Auth.JWTAccessSecret),
		SuccessHandler: jwtSuccessHandler,
		ErrorHandler:   jwtErrorHandler,
	}))
	clientG.Put("/", h.update).Name("ClientPut")
	//clientG.Put("/payment-password", h.changePaymentPassword).Name("ClientPutPaymentPassword")
	clientG.Get("/me", h.me).Name("ClientGetMe")

	clientG.Get("/credential", h.credential).Name("ClientGetCredential")

	h.svcAuth = service.NewAuth()
	h.svcCredential = service.NewCredential()

	return h
}

/* {{{ [Routers] - Definitions */

// token: Get JWT token

// @Tags Client
// @Summary Get authorize token of client
// @Description 通过client的身份（登录）获取认证信息，包括access_token和refresh_token两个jwt
// @ID ClientPostToken
// @Produce json
// @Param data body request.ClientPostToken true "Input information"
// @Success 201 {object} response.ClientPostToken
// @Success 226 {object} response.ClientPostToken
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Failure 401 {object} nil
// @Router /client/token [post]
func (h *Client) token(c *fiber.Ctx) error {
	var (
		errResp *utils.Envelope
		clt     *model.Client
	)
	idToken := c.Get("Authorization")
	if idToken != "" {
		// Firebase authorization
		errResp, clt = h.tokenFirebase(c, idToken)
	} else {
		// Normal login
		errResp, clt = h.tokenNormal(c)
	}

	if errResp != nil {
		return c.Status(errResp.Status).JSON(errResp)
	}

	// AccessToken
	jwtAccess, errAccess := h.svcAuth.JWTSign(&service.Sign{
		Sub:       clt.ID,
		Name:      clt.Email,
		Type:      "client",
		ExpiresIn: time.Duration(runtime.Config.Auth.JWTAccessExpiry) * time.Minute,
	})
	// RefreshToken
	jwtRefresh, errRefresh := h.svcAuth.JWTSign(&service.Sign{
		Sub:       clt.ID,
		Name:      clt.Email,
		Type:      "client::refresh",
		ExpiresIn: time.Duration(runtime.Config.Auth.JWTRefreshExpiry) * time.Minute,
	})
	if errAccess != nil || errRefresh != nil {
		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeAuthInternal
		resp.Message = response.MsgAuthInternal

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.ClientPostToken{
		AccessToken:   jwtAccess.Token,
		RefreshToken:  jwtRefresh.Token,
		AccessExpiry:  jwtAccess.Expiry,
		RefreshExpiry: jwtRefresh.Expiry,
		TokenType:     "bearer",
	})

	if clt.Phone == "" {
		resp.Status = fiber.StatusIMUsed
		c = c.Status(fiber.StatusIMUsed)
	} else {
		resp.Status = fiber.StatusCreated
		c = c.Status(fiber.StatusCreated)
	}

	return c.JSON(resp)
}

func (h *Client) tokenNormal(c *fiber.Ctx) (*utils.Envelope, *model.Client) {
	resp := utils.WrapResponse(nil)
	var req request.ClientPostToken
	err := c.BodyParser(&req)
	if err != nil {
		runtime.Logger.Warnf("parse request body failed: %s", err)

		resp.Status = fiber.StatusBadRequest
		resp.Message = err.Error()

		return resp, nil
	}

	if req.Email == "" || req.Password == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Code = response.CodeInvalidEmailOrPassword
		resp.Message = response.MsgInvalidEmailOrPassword

		return resp, nil
	}

	clt := &model.Client{
		Email: req.Email,
	}
	err = clt.Get(c.Context())
	if err != nil {
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeClientGetError
		resp.Message = response.MsgClientGetError

		return resp, nil
	}

	if clt.ID == "" {
		// Client not exists
		runtime.Logger.Warnf("try to fetch a nonexistent client of email [%s]", clt.Email)

		resp.Status = fiber.StatusUnauthorized
		resp.Code = response.CodeClientDoesNotExists
		resp.Message = response.MsgClientDoesNotExists

		return resp, nil
	}

	// Check password
	check := utils.EncryptPassword(req.Password, clt.Salt, clt.Email)
	if check != clt.Password {
		runtime.Logger.Warnf("wrong password given for cleitn [%s]", clt.Email)

		resp.Status = fiber.StatusUnauthorized
		resp.Code = response.CodeClientWrongPassword
		resp.Message = response.MsgClientWrongPassword

		return resp, nil
	}

	return nil, clt
}

func (h *Client) tokenFirebase(c *fiber.Ctx, idToken string) (*utils.Envelope, *model.Client) {
	ctx := c.Context()
	resp := utils.WrapResponse(nil)
	claims, err := h.svcAuth.FirebaseAuth(ctx, idToken)
	if err != nil {
		runtime.Logger.Errorf("Authenticate with firebase failed : %s", err)

		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeFirebaseFailed
		resp.Message = response.MsgFirebaseFailed

		return resp, nil
	}

	// Check existing
	clt := &model.Client{}
	clt.Email, _ = claims["email"].(string)
	err = clt.Get(ctx)
	if err != nil && err != sql.ErrNoRows {
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeClientGetError
		resp.Message = response.MsgClientGetError

		return resp, nil
	}

	if clt.ID == "" {
		// Create new client
		clt.Name, _ = claims["name"].(string)
		err = clt.Create(ctx)
		if err != nil {
			resp.Status = fiber.StatusInternalServerError
			resp.Code = response.CodeClientCreateError
			resp.Message = response.MsgClientCreateError

			return resp, nil
		}
	}

	return nil, clt
}

// refresh: Refresh JWT token

// @Tags Client
// @Summary Refresh access_token via refresh_token
// @Description 使用refresh_token获取新的access_token，避免客户端重复登录
// @ID ClientPostRefresh
// @Produce json
// @Success 201 {object} response.ClientPostRefresh
// @Failure 401 {object} nil
// @Failure 500 {object} nil
// @Router /client/refresh [post]
func (h *Client) refresh(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if len(auth) < 8 || strings.ToLower(auth[0:7]) != "bearer " {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeClientInvalidAuthorization
		resp.Message = response.MsgClientInvalidAuthorization
		resp.Status = fiber.StatusUnauthorized

		return c.Status(fiber.StatusUnauthorized).JSON(resp)
	}

	refreshToken := auth[7:]
	claims, err := h.svcAuth.JWTValid(refreshToken)
	if err != nil {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInternal
		resp.Message = response.MsgAuthInternal
		resp.Status = fiber.StatusUnauthorized

		return c.Status(fiber.StatusUnauthorized).JSON(resp)
	}

	// Re-sign token
	id, _ := claims["sub"].(string)
	name, _ := claims["name"].(string)
	jwt, err := h.svcAuth.JWTSign(&service.Sign{
		Sub:       id,
		Name:      name,
		Type:      "client",
		ExpiresIn: time.Duration(runtime.Config.Auth.JWTAccessExpiry) * time.Minute,
	})
	if err != nil {
		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeAuthInternal
		resp.Message = response.MsgAuthInternal

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.ClientPostRefresh{
		AccessToken:  jwt.Token,
		AccessExpiry: jwt.Expiry,
		TokenType:    "bearer",
	})
	resp.Status = fiber.StatusCreated

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// update: Update client

// @Tags Client
// @Summary Update client
// @Description 更新client
// @ID ClientPut
// @Produce json
// @Param data body request.ClientPut true "input information"
// @Success 200 {object} response.ClientPut
// @Failure 422 string message
// @Failure 400 {object} nil
// @Failure 401 {object} nil
// @Failure 500 {object} nil
// @Router /client/password [put]
func (h *Client) update(c *fiber.Ctx) error {
	return nil
}

/*
// changePaymentPassword: Change payment password

// @Tags Client
// @Summary Change payment password
// @Description 修改支付密码（6位数字）
// @ID ClientPutPaymentPassword
// @Produce json
// @Param data body request.ClientPutPaymentPassword true "input information"
// @Success 200 {object} response.ClientPutPaymentPassword
// @Failure 422 string message
// @Failure 400 {object} nil
// @Failure 401 {object} nil
// @Failure 500 {object} nil
// @Router /client/payment-password [put]
func (h *Client) changePaymentPassword(c *fiber.Ctx) error {
	return nil
}
*/

// me: Get myself

// @Tags Client
// @Summary Show me
// @Description 解析access_token，返回当前验证者信息（脱敏）
// @ID ClientGetMe
// @Produce json
// @Success 200 {object} nil
// @Failure 500 {object} nil
// @Router /client/me [get]
func (h *Client) me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return errors.New("parse userdata from context failed")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("wrong userdata format")
	}

	resp := utils.WrapResponse(claims)

	return c.JSON(resp)
}

// credential: Get information for QR render

// @Tags Client
// @Summary Get credential
// @Description 获取客户签名（二维码、付款码），默认5分钟有效期，签名内信息AES加密。URL后缀?img=true，生成512x512像素的png图片（二维码）
// @ID ClientGetCredential
// @Produce json
// @Success 200 {object} nil
// @Success 200 string png
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Router /client/credential [get]
func (h *Client) credential(c *fiber.Ctx) error {
	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || t != "client" {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	credential, err := h.svcCredential.Encode(id)
	if err != nil {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeEncodeFailed
		resp.Message = response.MsgEncodeFailed
		resp.Status = fiber.StatusInternalServerError

		runtime.Logger.Errorf("AES crypt failed : %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	if c.Query("img") != "" {
		// Render image
		png, err := qrcode.Encode(credential, qrcode.High, 512)
		if err != nil {
			return err
		}

		c.Set("Content-Type", "image/png")
		c.Write(png)

		return nil
	}

	resp := utils.WrapResponse(&response.ClientGetCredential{
		Credential: credential,
	})

	return c.JSON(resp)
}

/* }}} */

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
