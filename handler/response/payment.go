/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file payment.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package response

/* {{{ [Response codes && messages] */
const (
	CodePaymentCreateFailed = 13500001
)

const (
	MsgPaymentCreateFailed = "Create payment failed"
)

/* }}} */

type PaymentPost struct {
	TransactionID string `json:"transaction_id" xml:"transaction_id"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
