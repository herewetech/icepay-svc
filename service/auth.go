/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file auth.go
 * @package service
 * @author Dr.NP <np@herewe.tech>
 * @since 02/25/2023
 */

package service

import (
	"icepay-svc/runtime"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct{}

func NewAuth() *Auth {
	s := new(Auth)

	return s
}

func (s *Auth) JWTSign(id, name string) (string, time.Time, error) {
	exp := time.Now().Add(time.Duration(runtime.Config.Auth.JWTExpiry) * time.Minute)
	claims := jwt.MapClaims{
		"issuer": runtime.AppName,
		"sub":    id,
		"name":   name,
		"exp":    exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, err := token.SignedString([]byte(runtime.Config.Auth.JWTSecret))
	if err != nil {
		runtime.Logger.Errorf("sign JWT token failed : %s", err.Error())

		return "", time.Time{}, nil
	}

	return ts, exp, nil
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
