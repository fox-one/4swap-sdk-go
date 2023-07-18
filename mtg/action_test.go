package mtg

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checksum(t *testing.T) {
	var b [100]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		t.Error(err)
	}

	sum := checksum(b[:])
	assert.Len(t, sum, 4)

	h := sha256.Sum256(b[:])
	h = sha256.Sum256(h[:])
	assert.True(t, bytes.HasPrefix(h[:], sum))
}
