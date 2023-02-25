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
	"icepay-svc/runtime"

	"github.com/gofiber/fiber/v2"
)

type Tenant struct{}

func InitTenant() *Tenant {
	h := new(Tenant)

	tenantG := runtime.Server.Group("/tenant")
	tenantG.Post("/token", h.token).Name("TenantPostToken")
	tenantG.Post("/refresh", h.refresh).Name("TenantPostRefresh")
	tenantG.Put("/password", h.changePassword).Name("TenantPutPassword")

	return h
}

/* {{{ [Routers] - Definitions */

// token : Get JWT token
func (h *Tenant) token(c *fiber.Ctx) error {
	return nil
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
