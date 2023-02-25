/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file misc.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 02/25/2023
 */

package handler

import (
	"icepay-svc/runtime"
	"icepay-svc/utils"

	"github.com/gofiber/fiber/v2"
)

type Misc struct {
}

func InitMisc() *Misc {
	h := new(Misc)

	// Load routers
	runtime.Server.All("/", h.Index).Name("Index")
	runtime.Server.Get("/routers", h.Routers).Name("GetRouters")

	return h
}

func (h *Misc) Index(c *fiber.Ctx) error {
	return c.JSON(utils.WrapResponse(nil))
}

func (h *Misc) Routers(c *fiber.Ctx) error {
	return c.JSON(utils.WrapResponse(runtime.Server.Stack()))
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
