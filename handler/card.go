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
	"database/sql"
	"errors"
	"icepay-svc/handler/request"
	"icepay-svc/handler/response"
	"icepay-svc/model"
	"icepay-svc/runtime"
	"icepay-svc/service"
	"icepay-svc/utils"

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
	CardG.Post("/", h.add).Name("CardPost")
	CardG.Delete("/:id", h.delete).Name("CardDelete")
	CardG.Get("/list", h.list).Name("CardGetList")
	CardG.Get("/:id", h.get).Name("CardGet")
	CardG.Put("/:id", h.update).Name("CardUpdate")

	h.svcCard = service.NewCard()

	return h
}

/* {{{ [Routers] - Definitions */

// add: Add card to client

// @Tags Card
// @Summary Add bank card
// @Description 添加银行卡，自动识别信用卡或借记卡(debit)。当返回的卡片类型不是借记卡时，可进一步调用update设置信用卡信息（持卡人、有效期、CVV）。
// @ID CardPost
// @Produce json
// @Param data body request.CardPost true "Input information"
// @Success 201 {object} response.CardPost
// @Failure 422 string message
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Router /card [post]
func (h *Card) add(c *fiber.Ctx) error {
	var req request.CardPost
	err := c.BodyParser(&req)
	if err != nil {
		runtime.Logger.Warnf("parse request body failed : %s", err)
		c.SendStatus(fiber.StatusBadRequest)

		return err
	}

	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || (t != "client" && t != "tenant") {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	card, err := h.svcCard.Create(c.Context(), &model.Card{
		OwnerID:   id,
		OwnerType: t,
		Number:    req.Number,
	})
	if err != nil {
		runtime.Logger.Warnf("create card failed : %s", err)
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeCardCreateFailed
		resp.Message = response.MsgCardCreateFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.CardPost{
		ID:       card.ID,
		CardType: card.CardType,
		Number:   card.Number,
	})
	resp.Status = fiber.StatusCreated

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// delete: Remove card from client

// @Tags Card
// @Summary Delete bank card
// @Description 移除银行卡，只对当前验证者自己的卡生效
// @ID CardDelete
// @Produce json
// @Success 200 {object} nil
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Failure 404 {object} nil 卡片不存在
// @Router /card/{:id} [delete]
func (h *Card) delete(c *fiber.Ctx) error {
	cardID := c.Params("id")
	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || (t != "client" && t != "tenant") {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	err := h.svcCard.Delete(c.Context(), &model.Card{
		OwnerID:   id,
		OwnerType: t,
		ID:        cardID,
	})
	if errors.Is(err, model.ErrCardDoesNotExists) {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeTargetNotFound
		resp.Message = response.MsgTargetNotFound
		resp.Status = fiber.StatusNotFound

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	if err != nil {
		runtime.Logger.Warnf("delete card failed : %s", err)
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeCardDeleteFailed
		resp.Message = response.MsgCardDeleteFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(nil)

	return c.JSON(resp)
}

// get: Get card by given id

// @Tags Card
// @Summary Get bank card information
// @Description 获取银行卡信息，只对当前验证者自己的卡有效
// @ID CardGet
// @Produce json
// @Success 200 {object} response.CardGet
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Failure 404 {object} nil 卡片不存在
// @Router /card/{:id} [get]
func (h *Card) get(c *fiber.Ctx) error {
	cardID := c.Params("id")
	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || (t != "client" && t != "tenant") {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	ret, err := h.svcCard.Get(c.Context(), &model.Card{
		OwnerID:   id,
		OwnerType: t,
		ID:        cardID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeTargetNotFound
		resp.Message = response.MsgTargetNotFound
		resp.Status = fiber.StatusNotFound

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	if ret == nil || err != nil {
		runtime.Logger.Warnf("get card failed : %s", err)
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeCardGetFailed
		resp.Message = response.MsgCardGetFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := utils.WrapResponse(&response.CardGet{
		ID:       cardID,
		Number:   ret.Number,
		CardType: ret.CardType,
	})

	return c.JSON(resp)
}

// list: Get cards (list) of current login client

// @Tags Card
// @Summary Get bank cards list
// @Description 获取银行卡列表，限制范围是当前验证者所属
// @ID CardGetList
// @Produce json
// @Success 200 {object} []response.CardGet
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Router /card/list [get]
func (h *Card) list(c *fiber.Ctx) error {
	id, _ := c.Locals("AuthID").(string)
	t, _ := c.Locals("AuthType").(string)
	if id == "" || (t != "client" && t != "tenant") {
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeAuthInformationMissing
		resp.Message = response.MsgAuthInformationMissing
		resp.Status = fiber.StatusBadRequest

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	ret, err := h.svcCard.List(c.Context(), &model.Card{
		OwnerID:   id,
		OwnerType: t,
	})
	if err != nil {
		runtime.Logger.Warnf("get card failed : %s", err)
		resp := utils.WrapResponse(nil)
		resp.Code = response.CodeCardGetFailed
		resp.Message = response.MsgCardGetFailed
		resp.Status = fiber.StatusInternalServerError

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	cards := &response.CardGetList{
		Total: len(ret),
		List:  make([]*response.CardGet, len(ret)),
	}
	for idx, card := range ret {
		cards.List[idx] = &response.CardGet{
			ID:       card.ID,
			Number:   card.Number,
			CardType: card.CardType,
		}
	}

	return c.JSON(utils.WrapResponse(cards))
}

// update: Update card information

// @Tags Card
// @Summary Update bank (credit) card
// @Description 更新银行卡（信用卡）信息
// @ID CardUpdate
// @Produce json
// @Param data body request.CardUpdate true "Input information"
// @Success 200 {object} nil
// @Failure 422 string message
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Failure 404 {object} nil 卡片不存在
// @Router /card/{:id} [put]
func (h *Card) update(c *fiber.Ctx) error {
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
