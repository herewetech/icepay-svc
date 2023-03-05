/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file card.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package response

/* {{{ [Response codes && messages] */
const (
	CodeCardCreateFailed = 12500001
	CodeCardDeleteFailed = 12500002
	CodeCardGetFailed    = 12500003
)

const (
	MsgCardCreateFailed = "Create card failed"
	MsgCardDeleteFailed = "Delete card failed"
	MsgCardGetFailed    = "Get card failed"
)

/* }}} */

type CardPost struct {
	ID       string `json:"id" xml:"id"`
	Number   string `json:"number" xml:"number"`
	CardType string `json:"card_type" xml:"card_type"`
}

type CardDelete struct{}

type CardGet struct {
	ID       string `json:"id" xml:"id"`
	Number   string `json:"number" xml:"number"`
	CardType string `json:"card_type" xml:"card_type"`
}

type CardGetList struct {
	List  []*CardGet `json:"list" xml:"list"`
	Total int        `json:"total" xml:"total"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
