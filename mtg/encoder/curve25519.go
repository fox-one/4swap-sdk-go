package encoder

import (
	"crypto/ed25519"
	"crypto/sha512"

	"filippo.io/edwards25519"
)

func privateKeyToCurve25519(dst *[32]byte, privateKey ed25519.PrivateKey) {
	h := sha512.New()
	h.Write(privateKey.Seed())
	digest := h.Sum(nil)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(dst[:], digest)
}

func publicKeyToCurve25519(dst *[32]byte, publicKey ed25519.PublicKey) error {
	p, err := (&edwards25519.Point{}).SetBytes(publicKey[:])
	if err != nil {
		return err
	}

	copy(dst[:], p.BytesMontgomery())
	return nil
}
