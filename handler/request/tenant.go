/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file tenant.go
 * @package request
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package request

type TenantPostToken struct {
	Email    string `json:"email" xml:"email"`
	Password string `json:"password" xml:"password"`
}

type TenantPostRefresh struct{}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
