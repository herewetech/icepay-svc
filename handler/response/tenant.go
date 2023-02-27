/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file tenant.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package response

import "time"

/* {{{ [Response codes && messages] */
const (
	CodeTenantDoesNotExists = 11401001
	CodeTenantWrongPassword = 11401002
	CodeTenantGetError      = 11500001
)

const (
	MsgTenantDoesNotExists = "Tenant does not exists"
	MsgTenantWrongPassword = "Wrong tenant password"
	MsgTenantGetError      = "Get tenant from database error"
)

/* }}} */

type TenantPostToken struct {
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
