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
	"context"
	"fmt"
	"icepay-svc/runtime"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/option"
)

type Sign struct {
	Issuer    string
	Sub       string
	Name      string
	Type      string
	ExpiresIn time.Duration
}

type JWT struct {
	Token  string
	Expiry int64
}

type Auth struct {
	firebaseApp *firebase.App
}

func NewAuth() *Auth {
	s := new(Auth)

	fa, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(runtime.Config.Firebase.CredentialsFile))
	if err != nil {
		runtime.Logger.Fatal(err)
	}

	s.firebaseApp = fa

	return s
}

/* {{{ [Methods] */

// Sign JWT token
func (s *Auth) JWTSign(sign *Sign) (*JWT, error) {
	now := time.Now()
	exp := now.Add(sign.ExpiresIn).Unix()
	if sign.Issuer == "" {
		sign.Issuer = runtime.EnvPrefix + "::" + runtime.AppName
	}

	claims := jwt.MapClaims{
		"issuer": sign.Issuer,
		"sub":    sign.Sub,
		"name":   sign.Name,
		"iat":    now.Unix(),
		"exp":    exp,
		"type":   sign.Type,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, err := token.SignedString([]byte(runtime.Config.Auth.JWTAccessSecret))
	if err != nil {
		runtime.Logger.Errorf("sign JWT token failed : %s", err)

		return nil, err
	}

	return &JWT{
		Token:  ts,
		Expiry: exp,
	}, nil
}

// Valid JWT token
func (s *Auth) JWTValid(ts string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// Error
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}

		return []byte(runtime.Config.Auth.JWTRefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid claims format")
}

// Authenticate with firebase admin, idToken => claims
func (s *Auth) FirebaseAuth(ctx context.Context, idToken string) (map[string]interface{}, error) {
	fbClient, err := s.firebaseApp.Auth(ctx)
	if err != nil {
		return nil, err
	}

	gtk, err := fbClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return gtk.Claims, nil
}

// Deprecated, do nothing now...
func (s *Auth) CheckPassword(target, targetType, password string) (bool, error) {
	return true, nil
}

func (s *Auth) CheckPaymentPassword(client, password string) (bool, error) {
	return true, nil
}

/* }}} */

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
