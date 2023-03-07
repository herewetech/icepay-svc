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

/* {{{ [Response codes && messages] */
const (
	CodeClientDoesNotExists        = 10401001
	CodeClientWrongPassword        = 10401002
	CodeClientInvalidAuthorization = 10401010
	CodeClientGetError             = 10500001
)

const (
	MsgClientDoesNotExists        = "Client does not exists"
	MsgClientWrongPassword        = "Wrong client password"
	MsgClientInvalidAuthorization = "Invalid authorization information"
	MsgClientGetError             = "Get client from database error"
)

/* }}} */

type ClientPostToken struct {
	AccessToken   string `json:"access_token" xml:"access_token"`
	RefreshToken  string `json:"refresh_token" xml:"refresh_token"`
	AccessExpiry  int64  `json:"access_expiry" xml:"access_exipry"`
	RefreshExpiry int64  `json:"refresh_expiry" xml:"refresh_expiry"`
	TokenType     string `json:"token_type" xml:"token_type"`
}

type ClientPostRefresh struct {
	AccessToken  string `json:"access_token" xml:"access_token"`
	AccessExpiry int64  `json:"access_expiry" xml:"access_exipry"`
	TokenType    string `json:"token_type" xml:"token_type"`
}

type ClientPutPassword struct {
	Changed bool `json:"changed" xml:"changed"`
}

type ClientPutPaymentPassword struct {
	Changed bool `json:"changed" xml:"changed"`
}

type ClientGetMe struct {
	ID    string `json:"id" xml:"id"`
	Email string `json:"email" xml:"email"`
}

type ClientGetCredential struct {
	Credential string `json:"credential" xml:"credential"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
