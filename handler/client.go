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
	"icepay-svc/handler/request"
	"icepay-svc/runtime"
	"icepay-svc/service"
	"icepay-svc/utils"

	"github.com/gofiber/fiber/v2"
)

type Client struct {
	svcAuth *service.Auth
}

func InitClient() *Client {
	h := new(Client)

	clientG := runtime.Server.Group("/client")
	clientG.Post("/token", h.token).Name("ClientPostToken")
	clientG.Post("/refresh", h.refresh).Name("ClientPostRefresh")
	clientG.Put("/password", h.changePassword).Name("ClientPutPassword")
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

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	return nil
}

// refresh: Refresh JWT token
func (h *Client) refresh(c *fiber.Ctx) error {
	return nil
}

// changePassword: Change password
func (h *Client) changePassword(c *fiber.Ctx) error {
	return nil
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
