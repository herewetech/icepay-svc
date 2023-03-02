/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file credit_card.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package handler

import (
	"icepay-svc/runtime"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type CreditCard struct{}

func InitCreditCard() *CreditCard {
	h := new(CreditCard)

	creditCardG := runtime.Server.Group("/credit_card")
	creditCardG.Use(jwtware.New(jwtware.Config{
		SigningKey:     []byte(runtime.Config.Auth.JWTAccessSecret),
		SuccessHandler: jwtSuccessHandler,
		ErrorHandler:   jwtErrorHandler,
	}))
	creditCardG.Post("/", h.add).Name("CreditCardPost")
	creditCardG.Delete("/:id", h.delete).Name("CreditCardDelete")
	creditCardG.Get("/:id", h.get).Name("CreditCardGet")
	creditCardG.Get("/list", h.list).Name("CreditCardGetList")

	return h
}

/* {{{ [Routers] - Definitions */

// add: Add card to client
func (h *CreditCard) add(c *fiber.Ctx) error {
	return nil
}

// delete: Remove card from client
func (h *CreditCard) delete(c *fiber.Ctx) error {
	return nil
}

// get: Get card by given id
func (h *CreditCard) get(c *fiber.Ctx) error {
	return nil
}

// list: Get cards (list) of current login client
func (h *CreditCard) list(c *fiber.Ctx) error {
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
