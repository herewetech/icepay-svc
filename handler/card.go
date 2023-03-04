/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file card.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package handler

import (
	"icepay-svc/runtime"
	"icepay-svc/service"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type Card struct {
	svcCard *service.Card
}

func InitCreditCard() *Card {
	h := new(Card)

	CardG := runtime.Server.Group("/card")
	CardG.Use(jwtware.New(jwtware.Config{
		SigningKey:     []byte(runtime.Config.Auth.JWTAccessSecret),
		SuccessHandler: jwtSuccessHandler,
		ErrorHandler:   jwtErrorHandler,
	}))
	CardG.Post("/", h.add).Name("CreditCardPost")
	CardG.Delete("/:id", h.delete).Name("CreditCardDelete")
	CardG.Get("/:id", h.get).Name("CreditCardGet")
	CardG.Get("/list", h.list).Name("CreditCardGetList")

	h.svcCard = service.NewCard()

	return h
}

/* {{{ [Routers] - Definitions */

// add: Add card to client
func (h *Card) add(c *fiber.Ctx) error {
	//h.svcCard.Create()
	return nil
}

// delete: Remove card from client
func (h *Card) delete(c *fiber.Ctx) error {
	return nil
}

// get: Get card by given id
func (h *Card) get(c *fiber.Ctx) error {
	return nil
}

// list: Get cards (list) of current login client
func (h *Card) list(c *fiber.Ctx) error {
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
