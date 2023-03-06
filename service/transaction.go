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
	"encoding/json"
	"errors"
	"icepay-svc/model"
	"icepay-svc/runtime"
	"time"
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

// Update
func (s *Transaction) Update(ctx context.Context, input *model.Transaction) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:     input.ID,
		Client: input.Client,
		Tenant: input.Tenant,
		Status: input.Status,
		Card:   input.Card,
	}

	if transaction.Status == "" {
		transaction.Status = TransactionStatusClosed
	}

	err := transaction.Update(ctx)
	if err != nil {
		return nil, nil
	}

	return transaction, nil
}

// Get
func (s *Transaction) Get(ctx context.Context, input *model.Transaction) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:     input.ID,
		Client: input.Client,
		Tenant: input.Tenant,
		Status: input.Status,
	}

	err := transaction.Get(ctx)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// List
func (s *Transaction) List(ctx context.Context, input *model.Transaction) ([]*model.Transaction, error) {
	transaction := &model.Transaction{
		Client: input.Client,
		Tenant: input.Tenant,
		Status: input.Status,
	}

	list, err := transaction.List(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Notify
func (s *Transaction) Notify(ctx context.Context, input *model.Transaction) error {
	sub := ""
	switch input.Status {
	case TransactionStatusPreCreate, TransactionStatusCreated:
		sub = "pay::client::" + input.Client
	case TransactionStatusComfirmed, TransactionStatusAborted:
		sub = "pay::client::" + input.Tenant
	}

	if sub == "" {
		// Do nothing
		return errors.New("No notification needed")
	}

	b, _ := json.Marshal(input)

	return runtime.Nats.Publish(sub, b)
}

// Wait
func (s *Transaction) Wait(ctx context.Context, subscriber, subscriberType string) (*model.Transaction, error) {
	sub := "pay::" + subscriberType + "::" + subscriber
	suber, err := runtime.Nats.SubscribeSync(sub)
	if err != nil {
		return nil, err
	}

	msg, err := suber.NextMsg(time.Duration(runtime.Config.HTTP.LongPollingTimeout) * time.Second)
	if err != nil {
		return nil, err
	}

	output := new(model.Transaction)
	json.Unmarshal(msg.Data, &output)

	return output, nil
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
