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

/* {{{ [Response codes && messages] */
const (
	CodeTenantDoesNotExists        = 11401001
	CodeTenantWrongPassword        = 11401002
	CodeTenantInvalidAuthorization = 11401010
	CodeTenantGetError             = 11500001
)

const (
	MsgTenantDoesNotExists        = "Tenant does not exists"
	MsgTenantWrongPassword        = "Wrong tenant password"
	MsgTenantInvalidAuthorization = "Invalid authorization information"
	MsgTenantGetError             = "Get tenant from database error"
)

/* }}} */

type TenantPostToken struct {
	AccessToken   string `json:"access_token" xml:"access_token"`
	RefreshToken  string `json:"refresh_token" xml:"refresh_token"`
	AccessExpiry  int64  `json:"access_expiry" xml:"access_exipry"`
	RefreshExpiry int64  `json:"refresh_expiry" xml:"refresh_expiry"`
	TokenType     string `json:"token_type" xml:"token_type"`
}

type TenantPostRefresh struct {
	AccessToken  string `json:"access_token" xml:"access_token"`
	AccessExpiry int64  `json:"access_expiry" xml:"access_exipry"`
	TokenType    string `json:"token_type" xml:"token_type"`
}

type TenantPutPassword struct {
	Changed bool `json:"changed" xml:"changed"`
}

type TenantGetMe struct {
	ID    string `json:"id" xml:"id"`
	Email string `json:"email" xml:"email"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
