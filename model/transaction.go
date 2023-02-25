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
	"encoding/json"
	"icepay-svc/runtime"
	"time"
)

type Transaction struct {
	ID       string `bun:"id,pk" json:"id"`
	Client   string `bun:"client,notnull" json:"client"`
	Tenant   string `bun:"tenant,notnull" json:"tenant"`
	Amount   int64  `bun:"amount,notnull" json:"amount"`
	Currency string `bun:"currency,notnull" json:"currency"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero" json:"-"`
}

/* {{{ [Actions] - Definitions */

// Create
func (m *Transaction) Create(ctx context.Context) error {
	_, err := runtime.DB.NewInsert().Model(m).Exec(ctx)
	if err == nil {
		runtime.Logger.Infof("Transaction [%s] created", m.ID)
	}

	return err
}

// Get
func (m *Transaction) Get(ctx context.Context) error {
	return runtime.DB.NewSelect().Model(m).Scan(ctx)
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
