package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

func PKCS7Padding(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	suffix := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, suffix...)
}

func UnPKCS7Padding(data []byte) []byte {
	l := len(data)
	if l == 0 {
		return nil
	}

	n := int(data[l-1])
	if n >= l {
		return nil
	}

	return data[:(l - n)]
}

// Encrypt aes encrypt
func Encrypt(data []byte, key, iv []byte) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if nil != err {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(b, iv)
	paddingData := PKCS7Padding(data)
	dst := make([]byte, len(paddingData))
	cbc.CryptBlocks(dst, paddingData)

	return dst, nil
}

// Decrypt aes decrypt
func Decrypt(data []byte, key, iv []byte) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if nil != err {
		return nil, err
	}

	if len(data)%b.BlockSize() != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	cbc := cipher.NewCBCDecrypter(b, iv)
	dst := make([]byte, len(data))
	cbc.CryptBlocks(dst, data)
	body := UnPKCS7Padding(dst)

	// validate padding
	n := len(dst) - len(body)
	if ok := n >= 1 && n <= aes.BlockSize; !ok {
		return nil, fmt.Errorf("padding %d out of range [1:%d]", n, aes.BlockSize)
	}

	for _, v := range dst[len(body):] {
		if int(v) != n {
			return nil, fmt.Errorf("invalid padding, expect %d but got %d", n, v)
		}
	}

	return body, nil
}
