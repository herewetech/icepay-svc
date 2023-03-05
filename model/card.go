/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file card.go
 * @package model
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2022
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

type Card struct {
	bun.BaseModel `bun:"table:card"`
	ID            string `bun:"id,pk" json:"id"`
	OwnerID       string `bun:"owner_id" json:"owner_id"`
	OwnerType     string `bun:"owner_type" json:"owner_type"`
	Number        string `bun:"number" json:"nunber"`
	CardType      string `bun:"card_type" json:"card_type"`
	Holder        string `bun:"holder" json:"holder"`
	Expiration    string `bun:"expiration" json:"expiration"`
	CVV           string `bun:"cvv" json:"cvv"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero" json:"-"`
}

var (
	ErrCardDoesNotExists = errors.New("Card does not exists")
)

/* {{{ [Actions] - Definitions */

// Create: creats card
func (m *Card) Create(ctx context.Context) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}

	_, err := runtime.DB.NewInsert().Model(m).Returning("").Exec(ctx)
	if err == nil {
		runtime.Logger.Infof("card [%s] created", m.ID)
	} else {
		runtime.Logger.Errorf("create card failed : %s", err.Error())
	}

	return err
}

// Get: gets card
func (m *Card) Get(ctx context.Context) error {
	err := runtime.DB.NewSelect().
		Model(m).
		Where("id = ?", m.ID).
		Where("owner_id = ?", m.OwnerID).
		Where("owner_type = ?", m.OwnerType).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			runtime.Logger.Warnf("card does not exists")
		} else {
			runtime.Logger.Errorf("get card failed: %s", err.Error())
		}
	}

	return err
}

// List: list card by given conditoins
func (m *Card) List(ctx context.Context) ([]*Card, error) {
	var cards []*Card
	err := runtime.DB.NewSelect().
		Model(&cards).
		Where("owner_id = ?", m.OwnerID).
		Where("owner_type = ?", m.OwnerType).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cards, nil
		} else {
			return nil, err
		}
	}

	return cards, nil
}

// Delete: delete(soft) card
func (m *Card) Delete(ctx context.Context) error {
	res, err := runtime.DB.NewDelete().
		Model(m).
		Where("id = ?", m.ID).
		Where("owner_id = ?", m.OwnerID).
		Where("owner_type = ?", m.OwnerType).
		Exec(ctx)
	if err == nil {
		n, _ := res.RowsAffected()
		if n == 0 {
			runtime.Logger.Warnf("card [%s] does not exists", m.ID)

			return ErrCardDoesNotExists
		} else {
			runtime.Logger.Infof("card [%s] deleted", m.ID)
		}
	} else {
		runtime.Logger.Errorf("card [%s] delete failed : %s", m.ID, err.Error())
	}

	return err
}

// Debug
func (m *Card) Debug() string {
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
