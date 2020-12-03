package aes

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	var (
		body = make([]byte, 255)
		key  = make([]byte, 16)
		iv   = make([]byte, 16)
	)

	_, _ = io.ReadFull(rand.Reader, body)
	_, _ = io.ReadFull(rand.Reader, key)
	_, _ = io.ReadFull(rand.Reader, iv)

	encrypted, err := Encrypt(body, key, iv)
	require.Nil(t, err)

	decrypted, err := Decrypt(encrypted, key, iv)
	require.Nil(t, err)

	t.Log(string(decrypted), err)

	require.Equal(t, body, decrypted)
}
