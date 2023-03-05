/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file credential.go
 * @package service
 * @author Dr.NP <np@herewe.tech>
 * @since 03/05/2023
 */

package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"icepay-svc/runtime"
	"icepay-svc/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CredentialSource struct {
	ID     string
	Expiry time.Time
}
type Credential struct{}

func NewCredential() *Credential {
	s := new(Credential)

	return s
}

/* {{{ [Methods] */
func (s *Credential) Encode(id string) (string, error) {
	sourceStr := fmt.Sprintf("%s@@%d", id, time.Now().Add(time.Duration(runtime.Config.Security.CredentialLifetime)*time.Minute).Unix())
	// AES encrypt
	cipher, err := utils.AESCrypt([]byte(sourceStr), []byte(runtime.Config.Security.AESKey))
	if err != nil {
		return "", err
	}

	cipherText := base64.StdEncoding.EncodeToString(cipher)
	str := fmt.Sprintf("icepay://%s", cipherText)

	return str, nil
}

func (s *Credential) Decode(credential string) (string, error) {
	if !strings.HasPrefix(credential, "icepay://") {
		return "", errors.New("Wrong credential format")
	}

	stream, err := base64.StdEncoding.DecodeString(credential[9:])
	if err != nil {
		return "", err
	}

	// AES decrypt
	source, err := utils.AESDecrypt(stream, []byte(runtime.Config.Security.AESKey))
	reSource := regexp.MustCompile(`(.+)@@(\d+)`)
	matches := reSource.FindSubmatch(source)
	if len(matches) == 3 {
		expiryUnixStamp, _ := strconv.ParseInt(string(matches[2]), 10, 63)
		if time.Now().Unix() < expiryUnixStamp {
			return string(matches[1]), nil
		} else {
			return "", errors.New("Credential expires")
		}
	}

	return "", errors.New("Invalid credential")
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
