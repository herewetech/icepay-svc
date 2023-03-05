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
	"fmt"
	"icepay-svc/handler/request"
	"icepay-svc/handler/response"
	"icepay-svc/runtime"
	"icepay-svc/service"
	"icepay-svc/utils"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type Payment struct {
	svcPayment    *service.Payment
	svcCredential *service.Credential
}

func InitPayment() *Payment {
	h := new(Payment)

	paymentG := runtime.Server.Group("/payment")
	paymentG.Use(jwtware.New(jwtware.Config{
		SigningKey:     []byte(runtime.Config.Auth.JWTAccessSecret),
		SuccessHandler: jwtSuccessHandler,
		ErrorHandler:   jwtErrorHandler,
	}))
	paymentG.Post("/", h.add).Name("PaymentPost")

	h.svcPayment = service.NewPayment()
	h.svcCredential = service.NewCredential()

	return h
}

/* {{{ [Routers] - Definitions */

// add: Create new payment
func (h *Payment) add(c *fiber.Ctx) error {
	var req request.PaymentPost
	err := c.BodyParser(&req)
	if err != nil {
		runtime.Logger.Warnf("parse request body failed : %s", err)
		c.SendStatus(fiber.StatusBadRequest)

		return err
	}

	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || t != "tenant" {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	clientID, err := h.svcCredential.Decode(req.Credential)
	if err != nil {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeDecodeFailed
		resp.Message = response.MsgDecodeFailed
		resp.Status = fiber.StatusInternalServerError

		runtime.Logger.Warnf("AES decrypt failed : %s", err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	fmt.Println(clientID)

	/*
		transaction, err := h.svcPayment.Create(c.Context(), id, req.Credential, req.Amount, req.Currency, req.Detail)
		if err != nil {
			runtime.Logger.Warnf("create payment failed : %s", err)
			resp := utils.WrapResponse(nil)
			resp.Code = response.CodePaymentCreateFailed
			resp.Message = response.MsgPaymentCreateFailed
			resp.Status = fiber.StatusInternalServerError

			return c.Status(fiber.StatusInternalServerError).JSON(resp)
		}
	*/

	resp := utils.WrapResponse(&response.PaymentPost{
		//ID: transaction.ID,
	})
	resp.Status = fiber.StatusCreated

	return c.Status(fiber.StatusCreated).JSON(resp)
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
