/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file encode.go
 * @package utils
 * @author Dr.NP <np@herewe.tech>
 * @since 03/02/2023
 */

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"fmt"
)

func AESCrypt(input, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := make([]byte, len(input)+(aes.BlockSize-len(input)%aes.BlockSize))
	copy(padded, input)
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCBCEncrypter(block, iv)
	output := make([]byte, len(padded))
	stream.CryptBlocks(output, padded)

	return output, nil
}

func AESDecrypt(input, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCBCDecrypter(block, iv)
	output := make([]byte, len(input))
	stream.CryptBlocks(output, input)

	return output, nil
}

func EncryptPassword(plainText, salt, ident string) string {
	hash := sha512.New()
	hash.Write([]byte(plainText))
	hash.Write([]byte(salt))
	hash.Write([]byte(ident))

	return fmt.Sprintf("%02x", hash.Sum(nil))
}

func EncryptPaymentPassword(plainText, salt, ident string) string {
	hash := sha512.New()
	hash.Write([]byte(ident))
	hash.Write([]byte(salt))
	hash.Write([]byte(plainText))

	return fmt.Sprintf("%02x", hash.Sum(nil))
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
