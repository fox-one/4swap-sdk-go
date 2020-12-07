package encoder

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

func DecodePrivateKey(s string) (ed25519.PrivateKey, error) {
	b := decodeBase64(s)
	if len(b) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key")
	}

	return b, nil
}

func DecodePublicKey(s string) (ed25519.PublicKey, error) {
	b := decodeBase64(s)
	if len(b) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key")
	}

	return b, nil
}

func decodeBase64(memo string) []byte {
	if b, err := base64.StdEncoding.DecodeString(memo); err == nil {
		return b
	}

	if b, err := base64.URLEncoding.DecodeString(memo); err == nil {
		return b
	}

	return []byte(memo)
}
