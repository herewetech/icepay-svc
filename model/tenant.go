/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file tenant.go
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

	"github.com/uptrace/bun"
)

type Tenant struct {
	bun.BaseModel `bun:"table:tenant"`
	ID            string `bun:"id,pk" json:"id"`
	Name          string `bun:"name,notnull" json:"name"`
	Email         string `bun:"email,notnull" json:"email"`
	Phone         string `bun:"phone" json:"phone"`
	Password      string `bun:"password" json:"password"`
	Salt          string `bun:"salt" json:"salt"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero" json:"-"`
}

/* {{{ [Actions] - Definitions */

// Create: creates client
func (m *Tenant) Create(ctx context.Context) error {
	_, err := runtime.DB.NewInsert().Model(m).Exec(ctx)
	if err == nil {
		runtime.Logger.Infof("tenent [%s] created", m.ID)
	} else {
		runtime.Logger.Errorf("create tenent failed : %s", err)
	}

	return err
}

// Get: gets client
func (m *Tenant) Get(ctx context.Context) error {
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
	if err != nil {
		runtime.Logger.Errorf("get tenent failed : %s", err)
	}

	return err
}

// Debug
func (m *Tenant) Debug() string {
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
