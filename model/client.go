/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file client.go
 * @package model
 * @author Dr.NP <np@herewe.tech>
 * @since 02/25/2022
 */

package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"icepay-svc/runtime"
	"icepay-svc/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Client struct {
	bun.BaseModel   `bun:"table:client"`
	ID              string `bun:"id,pk" json:"id"`
	Name            string `bun:"name,notnull" json:"name"`
	Email           string `bun:"email,notnull" json:"email"`
	Phone           string `bun:"phone" json:"phone"`
	Password        string `bun:"password" json:"password"`
	PaymentPassword string `bun:"payment_password" json:"payment_password"`
	Salt            string `bun:"salt" json:"salt"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero" json:"-"`
}

/* {{{ [Actions] - Definitions */

// Create: creates client
func (m *Client) Create(ctx context.Context) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}

	if m.Salt == "" {
		// Generate random salt
		m.Salt = utils.RandomString(32)
	}

	if m.Password == "" {
		m.Password = "__NOT_SET__"
	} else {
		// Encrypt
		m.Password = utils.EncryptPassword(m.Password, m.Salt, m.Email)
	}

	if m.PaymentPassword == "" {
		m.PaymentPassword = "__NOT_SET__"
	} else {
		// Encrypt
		m.PaymentPassword = utils.EncryptPaymentPassword(m.PaymentPassword, m.Salt, m.Email)
	}

	_, err := runtime.DB.NewInsert().Model(m).Returning("").Exec(ctx)
	if err == nil {
		runtime.Logger.Infof("client [%s] created", m.ID)
	} else {
		runtime.Logger.Errorf("create client failed : %s", err)
	}

	return err
}

// Get: gets client
func (m *Client) Get(ctx context.Context) error {
	sq := runtime.DB.NewSelect().Model(m)
	if m.ID != "" {
		sq = sq.Where("id = ?", m.ID)
	}

	if m.Email != "" {
		sq = sq.Where("email = ?", m.Email)
	}

	if m.Phone != "" {
		sq = sq.Where("phone = ?", m.Phone)
	}

	err := sq.Limit(1).Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		runtime.Logger.Errorf("get client failed : %s", err)
	}

	return err
}

// Update: updates client
func (m *Client) Update(ctx context.Context) error {
	uq := runtime.DB.NewUpdate().Model(m).Where("id = ?", m.ID)
	if m.Phone != "" {
		uq = uq.Set("phone = ?", m.Phone)
	}

	if m.Name != "" {
		uq = uq.Set("name = ?", m.Name)
	}

	if m.Password != "" {
		uq = uq.Set("password = ?", utils.EncryptPassword(m.Password, m.Salt, m.Email))
	}

	if m.PaymentPassword != "" {
		uq = uq.Set("payment_password = ?", utils.EncryptPaymentPassword(m.PaymentPassword, m.Salt, m.Email))
	}

	_, err := uq.Returning("").Exec(ctx)
	if err == nil {
		runtime.Logger.Infof("client [%s] updated", m.ID)
	} else {
		runtime.Logger.Errorf("update client failed : %s", err)
	}

	return err
}

// Debug
func (m *Client) Debug() string {
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
