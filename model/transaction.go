/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file transaction.go
 * @package model
 * @author Dr.NP <np@herewe.tech>
 * @since 02/25/2022
 */

package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"icepay-svc/runtime"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transaction"`
	ID            string `bun:"id,pk" json:"id"`
	Client        string `bun:"client,notnull" json:"client"`
	Tenant        string `bun:"tenant,notnull" json:"tenant"`
	Amount        int64  `bun:"amount,notnull" json:"amount"`
	Currency      string `bun:"currency,notnull" json:"currency"`
	Status        string `bun:"status" json:"status"`
	Card          string `bun:"card" json:"card"`
	Detail        string `bun:"detail" json:"detail"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero" json:"-"`
}

/* {{{ [Actions] - Definitions */

// Create
func (m *Transaction) Create(ctx context.Context) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}

	_, err := runtime.DB.NewInsert().Model(m).Returning("").Exec(ctx)
	if err == nil {
		runtime.Logger.Infof("Transaction [%s] created", m.ID)
	}

	return err
}

// Update: updates transcation
func (m *Transaction) Update(ctx context.Context) error {
	uq := runtime.DB.NewUpdate().Model(m).Set("status = ?", m.Status).Set("updated_at = CURRENT_TIMESTAMP")
	if m.Card != "" {
		uq = uq.Set("card = ?", m.Card)
	}

	if m.ID != "" {
		uq = uq.Where("id = ?", m.ID)
	}

	if m.Client != "" {
		uq = uq.Where("client = ?", m.Client)
	}

	if m.Tenant != "" {
		uq = uq.Where("tenent = ?", m.Tenant)
	}

	_, err := uq.Returning("").Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("Update transaction failed : %s", err)
	}

	return err
}

// Get
func (m *Transaction) Get(ctx context.Context) error {
	sq := runtime.DB.NewSelect().Model(m)
	if m.ID != "" {
		sq = sq.Where("id = ?", m.ID)
	}

	if m.Client != "" {
		sq = sq.Where("client = ?", m.Client)
	}

	if m.Tenant != "" {
		sq = sq.Where("tenent = ?", m.Tenant)
	}

	if m.Status != "" {
		sq = sq.Where("status = ?", m.Status)
	}

	err := sq.Limit(1).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			runtime.Logger.Warnf("transaction does not exists")
		} else {
			runtime.Logger.Errorf("get transaction failed : %s", err)
		}
	}

	return err
}

// List: list transaction by given conditions
func (m *Transaction) List(ctx context.Context) ([]*Transaction, error) {
	var transactions []*Transaction
	sq := runtime.DB.NewSelect().Model(&transactions)
	if m.Client != "" {
		sq = sq.Where("client = ?", m.Client)
	}

	if m.Tenant != "" {
		sq = sq.Where("tenant = ?", m.Tenant)
	}

	if m.Status != "" {
		sq = sq.Where("status = ?", m.Status)
	}

	err := sq.Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transactions, nil
		}

		return nil, err
	}

	return transactions, nil
}

// Debug
func (m *Transaction) Debug() string {
	b, _ := json.MarshalIndent(m, "", "  ")

	return string(b)
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
