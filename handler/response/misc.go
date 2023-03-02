/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file client.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 02/27/2023
 */

package response

/* {{{ [Response codes && messages] */
const (
	CodeInvalidEmailOrPassword = 20400001
	CodeAuthFailed             = 20401001
	CodeAuthInternal           = 20401500
)

const (
	MsgInvalidEmailOrPassword = "Invalid email or password"
	MsgAuthInternal           = "Authorization failed"
)

/* }}} */

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */