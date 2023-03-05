/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file payment.go
 * @package request
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package request

type PaymentPost struct {
	Credential string `json:"credential" xml:"credential"`
	Amount     int64  `json:"amount" xml:"amount"`
	Currency   string `json:"currency" xml:"currency"`
	Detail     string `json:"detail" xml:"detail"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
