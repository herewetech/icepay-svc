/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file card.go
 * @package request
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package request

type CardPost struct {
	Number string `json:"number" xml:"number"`
}

type CardDelete struct{}

type CardGet struct{}

type CardGetList struct{}

type CardUpdate struct {
	Holder     string `json:"holder" xml:"holder"`
	Expiration string `json:"expiration" xml:"expiration"`
	CVV        string `json:"cvv" xml:"cvv"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
