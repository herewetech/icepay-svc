/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file client.go
 * @package request
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package request

type ClientPostToken struct {
	Email    string `json:"email" xml:"email"`
	Password string `json:"password" xml:"password"`
}

type ClientPostRefresh struct{}

type ClientPutPassword struct {
	OldPassword string `json:"old_password" xml:"old_password"`
	NewPassword string `json:"new_password" xml:"new_password"`
}

type ClientPutPaymentPassword struct {
	OldPassword string `json:"old_password" xml:"old_password"`
	NewPassword string `json:"new_password" xml:"new_password"`
}

type ClientPut struct {
	Name               string `json:"name" xml:"name"`
	Phone              string `json:"phone" xml:"phone"`
	OldPassword        string `json:"old_password" xml:"old_password"`
	NewPassword        string `json:"new_password" xml:"new_password"`
	OldPaymentPassword string `json:"old_payment_password" xml:"old_payment_password"`
	NewPaymentPassword string `json:"new_payment_password" xml:"new_payment_password"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
