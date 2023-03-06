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
	"database/sql"
	"errors"
	"icepay-svc/handler/request"
	"icepay-svc/handler/response"
	"icepay-svc/model"
	"icepay-svc/runtime"
	"icepay-svc/service"
	"icepay-svc/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/nats-io/nats.go"
)

type Payment struct {
	svcTransaction *service.Transaction
	svcCredential  *service.Credential
	svcAuth        *service.Auth
	svcCard        *service.Card
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
	paymentG.Put("/:id", h.update).Name("PaymentPut")
	paymentG.Get("/list", h.list).Name("PaymentGetList")
	paymentG.Get("/status", h.status).Name("PaymentGetStatus")
	paymentG.Get("/:id", h.get).Name("PaymentGet")
	paymentG.Put("/:id", h.update).Name("PaymentPut")

	h.svcTransaction = service.NewTransaction()
	h.svcCredential = service.NewCredential()
	h.svcAuth = service.NewAuth()
	h.svcCard = service.NewCard()

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

	transaction, err := h.svcTransaction.Create(c.Context(), &model.Transaction{
		Client:   clientID,
		Tenant:   id,
		Amount:   req.Amount,
		Currency: req.Currency,
		Detail:   req.Detail,
	})
	if err != nil {
		runtime.Logger.Warnf("create payment failed : %s", err)
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodePaymentCreateFailed
		resp.Message = response.MsgPaymentCreateFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	// Create notification
	err = h.svcTransaction.Notify(c.Context(), transaction)
	if err != nil {
		runtime.Logger.Errorf("payment notify failed : %s", err)
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodePaymentNotifyFailed
		resp.Message = response.MsgPaymentNotifyFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.PaymentPost{
		TransactionID: transaction.ID,
	})
	resp.Status = fiber.StatusCreated

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// update: Update payment status
func (h *Payment) update(c *fiber.Ctx) error {
	var req request.PaymentPut
	err := c.BodyParser(&req)
	if err != nil {
		runtime.Logger.Warnf("parse request body failed : %s", err)
		c.SendStatus(fiber.StatusBadRequest)

		return err
	}

	req.Status = strings.ToUpper(req.Status)
	if req.Status != service.TransactionStatusComfirmed && req.Status != service.TransactionStatusAborted {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeInvalidParameter
		resp.Message = response.MsgInvalidParameter
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || t != "client" {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	hint := &model.Transaction{
		ID:     c.Params("id"),
		Status: req.Status,
	}
	if req.Status == service.TransactionStatusComfirmed {
		// Check payment password
		checked, _ := h.svcAuth.CheckPaymentPassword(id, req.PaymentPassword)
		if !checked {
			// Password mismatch
			runtime.Logger.Warnf("client [%s] try to confirm transaction [%s] with wrong password", id, hint.ID)
			resp := utils.WrapResponse(nil)
			resp.Code = response.CodePaymentUpdateFailed
			resp.Message = response.MsgPaymentUpdateFailed
			resp.Status = fiber.StatusUnauthorized

			return c.Status(fiber.StatusUnauthorized).JSON(resp)
		}

		// Check card
		card, _ := h.svcCard.Get(c.Context(), &model.Card{
			ID:        req.Card,
			OwnerID:   id,
			OwnerType: t,
		})
		if card == nil {
			// Not your card
			runtime.Logger.Warnf("client [%s] try to update transaction [%s] with invalid card [%s]", id, hint.ID, req.Card)
			resp := utils.WrapResponse(nil)
			resp.Code = response.CodePaymentUpdateFailed
			resp.Message = response.MsgPaymentUpdateFailed
			resp.Status = fiber.StatusBadRequest

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		hint.Card = card.ID
	}

	transaction, err := h.svcTransaction.Update(c.Context(), hint)
	if err != nil {
		runtime.Logger.Errorf("update payment failed : %s", err)
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodePaymentUpdateFailed
		resp.Message = response.MsgPaymentUpdateFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	// Create notification
	err = h.svcTransaction.Notify(c.Context(), transaction)
	if err != nil {
		runtime.Logger.Errorf("payment notify failed : %s", err)
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodePaymentNotifyFailed
		resp.Message = response.MsgPaymentNotifyFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.PaymentPut{
		TransactionID:     transaction.ID,
		TransactionStatus: transaction.Status,
	})

	return c.JSON(resp)
}

// get: Get payment by id
func (h *Payment) get(c *fiber.Ctx) error {
	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || (t != "client" && t != "tenant") {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	input := &model.Transaction{
		ID: c.Params("id"),
	}
	if t == "client" {
		input.Client = id
	} else {
		input.Tenant = id
	}

	transaction, err := h.svcTransaction.Get(c.Context(), input)
	if err != nil {
		resp := utils.WrapResponse(nil)
		if errors.Is(err, sql.ErrNoRows) {
			resp.Code = response.CodeTargetNotFound
			resp.Message = response.MsgTargetNotFound
			resp.Status = fiber.StatusNotFound

			return c.Status(fiber.StatusNotFound).JSON(resp)
		}

		runtime.Logger.Errorf("get payment failed : %s", err)
		resp.Code = response.CodePaymentGetFailed
		resp.Message = response.MsgPaymentGetFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.PaymentGet{
		ID:       transaction.ID,
		Client:   transaction.Client,
		Tenant:   transaction.Tenant,
		Amount:   transaction.Amount,
		Currency: transaction.Currency,
		Status:   transaction.Status,
		Detail:   transaction.Detail,
	})

	return c.JSON(resp)
}

// list: List payment for client / tenant
func (h *Payment) list(c *fiber.Ctx) error {
	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || (t != "client" && t != "tenant") {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	input := &model.Transaction{
		ID: c.Params("id"),
	}
	if t == "client" {
		input.Client = id
	} else {
		input.Tenant = id
	}

	ret, err := h.svcTransaction.List(c.Context(), input)
	if err != nil {
		runtime.Logger.Errorf("get payment list failed : %s", err)
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodePaymentListFailed
		resp.Message = response.MsgPaymentListFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	payments := &response.PaymentGetList{
		Total: len(ret),
		List:  make([]*response.PaymentGet, len(ret)),
	}
	for idx, payment := range ret {
		payments.List[idx] = &response.PaymentGet{
			ID:       payment.ID,
			Client:   payment.Client,
			Tenant:   payment.Tenant,
			Amount:   payment.Amount,
			Currency: payment.Currency,
			Status:   payment.Status,
			Detail:   payment.Detail,
		}
	}

	return c.JSON(utils.WrapResponse(payments))
}

// status: long-pull status event
func (h *Payment) status(c *fiber.Ctx) error {
	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || (t != "client" && t != "tenant") {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	transaction, err := h.svcTransaction.Wait(c.Context(), id, t)
	if err != nil {
		resp := utils.WrapResponse(nil)
		if errors.Is(err, nats.ErrTimeout) {
			runtime.Logger.Debugf("wait payment status timeout")
			resp.Code = response.CodeTimeout
			resp.Message = response.MsgTimeout
			resp.Status = fiber.StatusRequestTimeout

			return c.Status(fiber.StatusRequestTimeout).JSON(resp)
		}

		runtime.Logger.Warnf("wait payment status error : %s", err)
		resp.Code = response.CodePaymentWaitFailed
		resp.Message = response.MsgPaymentWaitFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.PaymentGet{
		ID:       transaction.ID,
		Client:   transaction.Client,
		Tenant:   transaction.Tenant,
		Amount:   transaction.Amount,
		Currency: transaction.Currency,
		Status:   transaction.Status,
		Detail:   transaction.Detail,
	})

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
