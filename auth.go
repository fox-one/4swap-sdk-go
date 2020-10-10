package fswap

import (
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
)

// GenerateToken create a new mixin authorization token
// store.Scope must be 'FULL'
func GenerateToken(store mixin.Keystore, exp time.Duration) (string, error) {
	auth, err := mixin.AuthFromKeystore(&store)
	if err != nil {
		return "", err
	}

	sig := mixin.SignRaw("GET", "/me", nil)
	id := uuid.Must(uuid.NewV4()).String()
	token := auth.SignToken(sig, id, exp)
	return token, nil
}
