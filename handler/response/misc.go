/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file misc.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 02/27/2023
 */

package response

/* {{{ [Response codes && messages] */
const (
	CodeInvalidEmailOrPassword = 20400001
	CodeInvalidParameter       = 20400002
	CodeAuthFailed             = 20401001
	CodeAuthInternal           = 20401500
	CodeAuthInformationMissing = 20401404
	CodeEncodeFailed           = 20500001
	CodeDecodeFailed           = 20500002
	CodeFirebaseFailed         = 20500010
	CodeTargetNotFound         = 20404001
	CodeTimeout                = 20408001
)

const (
	MsgInvalidEmailOrPassword = "Invalid email or password"
	MsgInvalidParameter       = "Invalid parameter"
	MsgAuthFailed             = "Authorization failed"
	MsgAuthInternal           = "Authorization internal error"
	MsgAuthInformationMissing = "Authorization information missing"
	MsgEncodeFailed           = "Encode failed"
	MsgDecodeFailed           = "Decode failed"
	MsgFirebaseFailed         = "Firebase failed"
	MsgTargetNotFound         = "Target not found"
	MsgTimeout                = "Timeout"
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
