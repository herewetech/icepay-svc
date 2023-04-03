/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file tenant.go
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
)

type Tenant struct {
	svcAuth *service.Auth
}

func InitTenant() *Tenant {
	h := new(Tenant)

	tenantG := runtime.Server.Group("/tenant")
	tenantG.Post("/token", h.token).Name("TenantPostToken")
	tenantG.Post("/refresh", h.refresh).Name("TenantPostRefresh")
	tenantG.Use(jwtware.New(jwtware.Config{
		SigningKey:     []byte(runtime.Config.Auth.JWTRefreshSecret),
		SuccessHandler: jwtSuccessHandler,
		ErrorHandler:   jwtErrorHandler,
	}))
	tenantG.Put("/password", h.changePassword).Name("TenantPutPassword")
	tenantG.Get("/me", h.me).Name("TenantGetMe")

	h.svcAuth = service.NewAuth()

	return h
}

/* {{{ [Routers] - Definitions */

// token : Get JWT token

// @Tags Tenant
// @Summary Get authorize token of tenant
// @Description 通过client的身份（登录）获取认证信息，包括access_token和refresh_token两个jwt
// @ID TenantPostToken
// @Produce json
// @Param data body request.ClientPostToken true "Input information"
// @Success 201 {object} response.TenantPostToken
// @Failure 422 string message
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Failure 401 {object} nil
// @Router /tenant/token [post]
func (h *Tenant) token(c *fiber.Ctx) error {
	var (
		errResp *utils.Envelope
		tnt     *model.Tenant
	)
	idToken := c.Get("Authorization")
	if idToken != "" {
		// Firebase authorization
		errResp, tnt = h.tokenFirebase(c, idToken)
	} else {
		// Normal login
		errResp, tnt = h.tokenNormal(c)
	}

	if errResp != nil {
		return c.Status(errResp.Status).JSON(errResp)
	}

	// AccessToken
	jwtAccess, errAccess := h.svcAuth.JWTSign(&service.Sign{
		Sub:       tnt.ID,
		Name:      tnt.Email,
		Type:      "tenant",
		ExpiresIn: time.Duration(runtime.Config.Auth.JWTAccessExpiry) * time.Minute,
	})
	// RefreshToken
	jwtRefresh, errRefresh := h.svcAuth.JWTSign(&service.Sign{
		Sub:       tnt.ID,
		Name:      tnt.Email,
		Type:      "tenant::refresh",
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

	if tnt.Phone == "" {
		resp.Status = fiber.StatusIMUsed
		c = c.Status(fiber.StatusIMUsed)
	} else {
		resp.Status = fiber.StatusCreated
		c = c.Status(fiber.StatusCreated)
	}

	return c.JSON(resp)
}

func (h *Tenant) tokenNormal(c *fiber.Ctx) (*utils.Envelope, *model.Tenant) {
	resp := utils.WrapResponse(nil)
	var req request.TenantPostToken
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

	tnt := &model.Tenant{
		Email: req.Email,
	}
	err = tnt.Get(c.Context())
	if err != nil {
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeTenantGetError
		resp.Message = response.MsgTenantGetError

		return resp, nil
	}

	if tnt.ID == "" {
		// Tenant not exists
		runtime.Logger.Warnf("try to fetch a nonexistent tenant of email [%s]", tnt.Email)

		resp.Status = fiber.StatusUnauthorized
		resp.Code = response.CodeTenantDoesNotExists
		resp.Message = response.MsgTenantDoesNotExists

		return resp, nil
	}

	// Check password
	check := utils.EncryptPassword(req.Password, tnt.Salt, tnt.Email)
	if check != tnt.Password {
		runtime.Logger.Warnf("wrong password given for cleitn [%s]", tnt.Email)

		resp.Status = fiber.StatusUnauthorized
		resp.Code = response.CodeTenantWrongPassword
		resp.Message = response.MsgTenantWrongPassword

		return resp, nil
	}

	return nil, tnt
}

func (h *Tenant) tokenFirebase(c *fiber.Ctx, idToken string) (*utils.Envelope, *model.Tenant) {
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
	tnt := &model.Tenant{}
	tnt.Email, _ = claims["email"].(string)
	err = tnt.Get(ctx)
	if err != nil && err != sql.ErrNoRows {
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeTenantGetError
		resp.Message = response.MsgTenantGetError

		return resp, nil
	}

	if tnt.ID == "" {
		// Create new client
		tnt.Name, _ = claims["name"].(string)
		err = tnt.Create(ctx)
		if err != nil {
			resp.Status = fiber.StatusInternalServerError
			resp.Code = response.CodeTenantCreateError
			resp.Message = response.MsgTenantCreateError

			return resp, nil
		}
	}

	return nil, tnt
}

// refresh: Refresh JWT token

// @Tags Tenant
// @Summary Refresh access_token via refresh_token
// @Description 使用refresh_token获取新的access_token，避免客户端重复登录
// @ID TenantPostRefresh
// @Produce json
// @Success 201 {object} response.TenantPostRefresh
// @Failure 401 {object} nil
// @Failure 500 {object} nil
// @Router /tenant/refresh [post]
func (h *Tenant) refresh(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if len(auth) < 8 || strings.ToLower(auth[0:7]) != "bearer " {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeTenantInvalidAuthorization
		resp.Message = response.MsgTenantInvalidAuthorization
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
		Type:      "tenant",
		ExpiresIn: time.Duration(runtime.Config.Auth.JWTAccessExpiry) * time.Minute,
	})
	if err != nil {
		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeAuthInternal
		resp.Message = response.MsgAuthInternal

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.TenantPostRefresh{
		AccessToken:  jwt.Token,
		AccessExpiry: jwt.Expiry,
		TokenType:    "bearer",
	})

	return c.JSON(resp)
}

// changePassword: Change password

// @Tags Tenant
// @Summary Change password
// @Description 修改登录密码
// @ID TenantPutPassword
// @Produce json
// @Param data body request.TenantPutPassword true "input information"
// @Success 200 {object} response.TenantPutPassword
// @Failure 422 string message
// @Failure 400 {object} nil
// @Failure 401 {object} nil
// @Failure 500 {object} nil
// @Router /tenant/password [put]
func (h *Tenant) changePassword(c *fiber.Ctx) error {
	return nil
}

// me: Get myself

// @Tags Tenant
// @Summary Show me
// @Description 解析access_token，返回当前验证者信息（脱敏）
// @ID TenantGetMe
// @Produce json
// @Success 200 {object} nil
// @Failure 500 {object} nil
// @Router /tenant/me [get]
func (h *Tenant) me(c *fiber.Ctx) error {
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

/* }}} */

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
