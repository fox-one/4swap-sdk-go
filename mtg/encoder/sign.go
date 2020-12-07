package encoder

import (
	"bytes"
	"crypto/ed25519"
	"io"
	"io/ioutil"
)

func Sign(body []byte, privateKey ed25519.PrivateKey) []byte {
	return ed25519.Sign(privateKey, body)
}

func Verify(body, sig []byte, publicKey ed25519.PublicKey) bool {
	return ed25519.Verify(publicKey, body, sig)
}

func Pack(body, sig []byte) []byte {
	b := make([]byte, len(body)+len(sig))
	n := copy(b, sig)
	copy(b[n:], body)
	return b
}

func Unpack(b []byte) (body, sig []byte, err error) {
	r := bytes.NewReader(b)
	sig = make([]byte, ed25519.SignatureSize)
	_, err = io.ReadFull(r, sig)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(r)
	return
}
