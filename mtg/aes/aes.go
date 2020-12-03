package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// PKCS7Padding PKCS7补码, 可以参考下http://blog.studygolang.com/167.html
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// UnPKCS7Padding 去除PKCS7的补码
func UnPKCS7Padding(data []byte) []byte {
	length := len(data)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(data[length-1])
	if length <= unpadding {
		return nil
	}
	return data[:(length - unpadding)]
}

// Encrypt aes encrypt
func Encrypt(data []byte, key, iv []byte) ([]byte, error) {
	ckey, err := aes.NewCipher(key)
	if nil != err {
		return nil, err
	}

	encrypter := cipher.NewCBCEncrypter(ckey, iv)

	// PKCS7补码
	str := PKCS7Padding(data, 16)
	out := make([]byte, len(str))

	encrypter.CryptBlocks(out, str)
	return out, nil
}

// Decrypt aes decrypt
func Decrypt(data []byte, key, iv []byte) ([]byte, error) {
	ckey, err := aes.NewCipher(key)
	if nil != err {
		return nil, err
	}

	if len(data)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	decrypter := cipher.NewCBCDecrypter(ckey, iv)

	out := make([]byte, len(data))
	decrypter.CryptBlocks(out, data)

	// 去除PKCS7补码
	out = UnPKCS7Padding(out)
	return out, nil
}
