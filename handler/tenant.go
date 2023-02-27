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
	"crypto/sha512"
	"fmt"
	"icepay-svc/handler/request"
	"icepay-svc/handler/response"
	"icepay-svc/model"
	"icepay-svc/runtime"
	"icepay-svc/service"
	"icepay-svc/utils"

	"github.com/gofiber/fiber/v2"
)

type Tenant struct {
	svcAuth *service.Auth
}

func InitTenant() *Tenant {
	h := new(Tenant)

	tenantG := runtime.Server.Group("/tenant")
	tenantG.Post("/token", h.token).Name("TenantPostToken")
	tenantG.Post("/refresh", h.refresh).Name("TenantPostRefresh")
	tenantG.Put("/password", h.changePassword).Name("TenantPutPassword")

	h.svcAuth = service.NewAuth()

	return h
}

/* {{{ [Routers] - Definitions */

// token : Get JWT token
func (h *Tenant) token(c *fiber.Ctx) error {
	var req request.TenantPostToken
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

	tnt := &model.Tenant{
		Email: req.Email,
	}
	err = tnt.Get(c.Context())
	if err != nil {
		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeTenantGetError
		resp.Message = response.MsgTenantGetError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	if tnt.ID == "" {
		// Client not exists
		runtime.Logger.Warnf("try to fetch a nonexistent tenant of email [%s]", tnt.Email)

		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusUnauthorized
		resp.Code = response.CodeTenantDoesNotExists
		resp.Message = response.MsgTenantDoesNotExists

		return c.Status(fiber.StatusUnauthorized).JSON(resp)
	}

	// Check password
	hash := sha512.New()
	hash.Write([]byte(req.Password))
	hash.Write([]byte(tnt.Salt))
	hash.Write([]byte(tnt.Email))
	check := fmt.Sprintf("%02x", hash.Sum(nil))
	if check != tnt.Password {
		runtime.Logger.Warnf("wrong password given for tenant [%s]", tnt.Email)

		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusUnauthorized
		resp.Code = response.CodeClientWrongPassword
		resp.Message = response.MsgClientWrongPassword

		return c.Status(fiber.StatusUnauthorized).JSON(resp)
	}

	jwt, exp, err := h.svcAuth.JWTSign(tnt.ID, tnt.Email)
	if err != nil {
		resp := utils.WrapResponse(nil)
		resp.Status = fiber.StatusInternalServerError
		resp.Code = response.CodeAuthInternal
		resp.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.TenantPostToken{
		Token:  jwt,
		Expiry: exp,
	})

	return c.JSON(resp)
}

// refresh : Refresh JWT token
func (h *Tenant) refresh(c *fiber.Ctx) error {
	return nil
}

// changePassword: Change password
func (h *Tenant) changePassword(c *fiber.Ctx) error {
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
