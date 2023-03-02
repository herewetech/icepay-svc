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
	"crypto/sha512"
	"errors"
	"fmt"
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
)

type Client struct {
	svcAuth *service.Auth
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
	clientG.Put("/password", h.changePassword).Name("ClientPutPassword")
	clientG.Get("/me", h.me).Name("ClientGetMe")

	clientG.Get("/credential", h.credential).Name("ClientGetCredential")

	h.svcAuth = service.NewAuth()

	return h
}

/* {{{ [Routers] - Definitions */

// token: Get JWT token
func (h *Client) token(c *fiber.Ctx) error {
	var req request.ClientPostToken
	err := c.BodyParser(&req)
	if err != nil {
		runtime.Logger.Warnf("parse request body failed: %s", err)

		return err
	}

	if req.Email == "" || req.Password == "" {
		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusBadRequest
		resp.Code = response.CodeInvalidEmailOrPassword
		resp.Message = response.MsgInvalidEmailOrPassword

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	clt := &model.Client{
		Email: req.Email,
	}
	err = clt.Get(c.Context())
	if err != nil {
		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeClientGetError
		resp.Message = response.MsgClientGetError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	if clt.ID == "" {
		// Client not exists
		runtime.Logger.Warnf("try to fetch a nonexistent client of email [%s]", clt.Email)

		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusUnauthorized
		resp.Code = response.CodeClientDoesNotExists
		resp.Message = response.MsgClientDoesNotExists

		return c.Status(fiber.StatusUnauthorized).JSON(resp)
	}

	// Check password
	hash := sha512.New()
	hash.Write([]byte(req.Password))
	hash.Write([]byte(clt.Salt))
	hash.Write([]byte(clt.Email))
	check := fmt.Sprintf("%02x", hash.Sum(nil))
	if check != clt.Password {
		runtime.Logger.Warnf("wrong password given for cleitn [%s]", clt.Email)

		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusUnauthorized
		resp.Code = response.CodeClientWrongPassword
		resp.Message = response.MsgClientWrongPassword

		return c.Status(fiber.StatusUnauthorized).JSON(resp)
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

	return c.JSON(resp)
}

// refresh: Refresh JWT token
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

	return c.JSON(resp)
}

// changePassword: Change password
func (h *Client) changePassword(c *fiber.Ctx) error {
	return nil
}

// me: Get myself
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
func (h *Client) credential(c *fiber.Ctx) error {
	return nil
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
