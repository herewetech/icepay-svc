/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file misc.go
 * @package utils
 * @author Dr.NP <np@herewe.tech>
 * @since 03/05/2023
 */

package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"unicode"
)

const (
	CardDinersClub      = "Diners Club"
	CardAmericanExpress = "American Express"
	CardJCB             = "JCB"
	CardVISA            = "VISA"
	CardMasterCard      = "Master Card"
	CardChinaUnionPay   = "China Union Pay"
	CardDebit           = "Debit"
	CardInvalid         = "Invalid"
)

func ValidCardNumber(number string) string {
	for _, c := range number {
		if !unicode.IsDigit(c) {
			return CardInvalid
		}
	}

	if len(number) == 14 {
		vPrefix, _ := strconv.Atoi(number[0:6])
		if (vPrefix >= 300000 && vPrefix <= 305999) ||
			(vPrefix >= 309500 && vPrefix <= 309599) ||
			(vPrefix >= 360000 && vPrefix <= 369999) ||
			(vPrefix >= 380000 && vPrefix <= 389999) {
			// DinersClub
			return CardDinersClub
		}
	}

	if len(number) == 15 {
		vPrefix, _ := strconv.Atoi(number[0:6])
		if (vPrefix >= 340000 && vPrefix <= 349999) ||
			(vPrefix >= 370000 && vPrefix <= 379999) {
			// AmericanExpress
			return CardAmericanExpress
		}
	}

	if len(number) == 16 {
		vPrefix, _ := strconv.Atoi(number[0:6])
		if vPrefix >= 352800 && vPrefix <= 358999 {
			// JCB
			return CardJCB
		}

		if strings.HasPrefix(number, "4") {
			// VISA
			return CardVISA
		}

		if vPrefix >= 510000 && vPrefix <= 559999 {
			// MasterCard
			return CardMasterCard
		}

		if vPrefix >= 622126 && vPrefix <= 622925 {
			// ChinaUnionPay
			return CardChinaUnionPay
		}
	}

	if len(number) > 13 && len(number) < 20 {
		return CardDebit
	}

	return CardInvalid
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
