/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file client.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package response

import "time"

/* {{{ [Response codes && messages] */
const (
	CodeInvalidEmailOrPassword = 10400001
	CodeClientDoesNotExists    = 10401001
	CodeClientWrongPassword    = 10401002
	CodeClientGetError         = 10500001
)

const (
	MsgInvalidEmailOrPassword = "Invalid email or password"
	MsgClientDoesNotExists    = "Client does not exists"
	MsgClientWrongPassword    = "Wrong client password"
	MsgClientGetError         = "Get client from database error"
)

/* }}} */

type ClientPostToken struct {
	Token  string    `json:"token" xml:"token"`
	Expiry time.Time `json:"expiry" xml:"exipry"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
