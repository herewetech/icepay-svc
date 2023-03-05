/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file card.go
 * @package service
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package service

import (
	"context"
	"errors"
	"icepay-svc/model"
	"icepay-svc/utils"
	"strings"
)

type Card struct{}

func NewCard() *Card {
	s := new(Card)

	return s
}

/* {{{ [Methods] */

// Create
func (s *Card) Create(ctx context.Context, ownerID, ownerType, number string) (*model.Card, error) {
	// Valid card number
	number = strings.TrimSpace(number)
	number = strings.ReplaceAll(number, " ", "")
	c := utils.ValidCardNumber(number)
	if c == utils.CardInvalid {
		return nil, errors.New("Invalid card number")
	}

	card := &model.Card{
		OwnerID:   ownerID,
		OwnerType: ownerType,
		Number:    number,
		CardType:  c,
	}

	err := card.Create(ctx)
	if err != nil {
		return nil, err
	}

	return card, nil
}

// Delete
func (s *Card) Delete(ctx context.Context, ownerID, ownerType, id string) error {
	card := &model.Card{
		OwnerID:   ownerID,
		OwnerType: ownerType,
		ID:        id,
	}

	return card.Delete(ctx)
}

// Get
func (s *Card) Get(ctx context.Context, ownerID, ownerType, id string) (*model.Card, error) {
	card := &model.Card{
		OwnerID:   ownerID,
		OwnerType: ownerType,
		ID:        id,
	}

	err := card.Get(ctx)
	if err != nil {
		return nil, err
	}

	return card, nil
}

// List
func (s *Card) List(ctx context.Context, ownerID, ownerType string) ([]*model.Card, error) {
	card := &model.Card{
		OwnerID:   ownerID,
		OwnerType: ownerType,
	}

	list, err := card.List(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
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
