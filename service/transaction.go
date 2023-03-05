/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file transaction.go
 * @package service
 * @author Dr.NP <np@herewe.tech>
 * @since 03/05/2023
 */

package service

import (
	"context"
	"icepay-svc/model"
)

const (
	TransactionStatusPreCreate = "PRE"
	TransactionStatusCreated   = "CREATED"
	TransactionStatusComfirmed = "CONFIRMED"
	TransactionStatusAborted   = "ABORTED"
	TransactionStatusClosed    = "CLOSED"
	TransactionStatusInvalid   = "INVALID"
)

type Transaction struct{}

func NewTransaction() *Transaction {
	s := new(Transaction)

	return s
}

/* {{{ [Methods] */

// Create
func (s *Transaction) Create(ctx context.Context, input *model.Transaction) (*model.Transaction, error) {
	transaction := &model.Transaction{
		Client:   input.Client,
		Tenant:   input.Tenant,
		Amount:   input.Amount,
		Currency: input.Currency,
		Status:   TransactionStatusCreated,
		Detail:   input.Detail,
	}

	err := transaction.Create(ctx)
	if err != nil {
		return nil, err
	}

	return transaction, nil
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
