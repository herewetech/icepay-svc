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
	CodePaymentDeleteFailed = 13500002
	CodePaymentUpdateFailed = 13500003
	CodePaymentGetFailed    = 13500004
	CodePaymentListFailed   = 13500005
	CodePaymentNotifyFailed = 13500098
	CodePaymentWaitFailed   = 13500099
)

const (
	MsgPaymentCreateFailed = "Create payment failed"
	MsgPaymentDeleteFailed = "Delete payment failed"
	MsgPaymentUpdateFailed = "Update payment failed"
	MsgPaymentGetFailed    = "Get payment failed"
	MsgPaymentListFailed   = "List payment failed"
	MsgPaymentNotifyFailed = "Notify payment failed"
	MsgPaymentWaitFailed   = "Wait payment failed"
)

/* }}} */

type PaymentPost struct {
	TransactionID string `json:"transaction_id" xml:"transaction_id"`
}

type PaymentPut struct {
	TransactionID     string `json:"transaction_id" xml:"transaction_id"`
	TransactionStatus string `json:"transaction_status" xml:"transaction_status"`
}

type PaymentGet struct {
	ID       string `json:"id" xml:"id"`
	Client   string `json:"client" xml:"client"`
	Tenant   string `json:"tenant" xml:"tenant"`
	Amount   int64  `json:"amount" xml:"amount"`
	Currency string `json:"currency" xml:"currency"`
	Status   string `json:"status" xml:"status"`
	Detail   string `json:"detail" xml:"detail"`
}

type PaymentGetList struct{}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
