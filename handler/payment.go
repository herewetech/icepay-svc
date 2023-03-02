/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file payment.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 02/25/2023
 */

package handler

import (
	"icepay-svc/runtime"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type Payment struct{}

func InitPayment() *Payment {
	h := new(Payment)

	paymentG := runtime.Server.Group("/payment")
	paymentG.Use(jwtware.New(jwtware.Config{
		SigningKey:     []byte(runtime.Config.Auth.JWTAccessSecret),
		SuccessHandler: jwtSuccessHandler,
		ErrorHandler:   jwtErrorHandler,
	}))
	paymentG.Post("/", h.add).Name("PaymentPost")

	return h
}

/* {{{ [Routers] - Definitions */

// add: Create new payment
func (h *Payment) add(c *fiber.Ctx) error {
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
