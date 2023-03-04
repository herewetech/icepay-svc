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

import "time"

type Card struct {
	ID         string `bun:"id,pk" json:"id"`
	OwnerID    string `bun:"owner_id" json:"owner_id"`
	OwnerType  string `bun:"owner_type" json:"owner_type"`
	Number     string `bun:"number" json:"nunber"`
	CardType   string `bun:"card_type" json:"card_type"`
	Holder     string `bun:"holder" json:"holder"`
	Expiration string `bun:"expiration" json:"expiration"`
	CVV        string `bun:"cvv" json:"cvv"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero" json:"-"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
