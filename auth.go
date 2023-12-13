package fswap

import (
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/gofrs/uuid"
)

// GenerateToken create a new mixin authorization token
// store.Scope must be 'FULL'
func GenerateToken(clientID, sessionID, sessionKey string, exp time.Duration) (string, error) {
	auth, err := mixin.AuthFromKeystore(&mixin.Keystore{
		ClientID:   clientID,
		SessionID:  sessionID,
		PrivateKey: sessionKey,
		Scope:      "FULL",
	})

	if err != nil {
		return "", err
	}

	sig := mixin.SignRaw("GET", "/me", nil)
	id := uuid.Must(uuid.NewV4()).String()
	token := auth.SignToken(sig, id, exp)
	return token, nil
}
