/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file payment.go
 * @package service
 * @author Dr.NP <np@herewe.tech>
 * @since 03/05/2023
 */

package service

import (
	"context"
	"errors"
	"icepay-svc/model"
	"strings"
)

type Payment struct{}

func NewPayment() *Payment {
	s := new(Payment)

	return s
}

/* {{{ [Methods] */

// Create
func (s *Payment) Create(ctx context.Context, tenentID, credential string, amount int64, currency, detail string) (*model.Transaction, error) {
	if !strings.HasPrefix(credential, "icepay://") {
		return nil, errors.New("Invalid credential")
	}

	return nil, nil
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
