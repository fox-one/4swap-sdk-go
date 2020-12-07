package encoder

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha512"
	"errors"
	"io"
	"io/ioutil"

	"github.com/fox-one/4swap-sdk-go/mtg/aes"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/mixin-sdk-go/edwards25519"
)

func Decrypt(b []byte, privateKey ed25519.PrivateKey) ([]byte, error) {
	r := bytes.NewReader(b)

	// read
	pub := make([]byte, 32)
	if _, err := io.ReadFull(r, pub); err != nil {
		return nil, err
	}

	key, iv, err := keyPairsToAesKeyIv(privateKey, pub)
	if err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(r)
	body, err := aes.Decrypt(data, key, iv)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Encrypt(body []byte, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) ([]byte, error) {
	key, iv, err := keyPairsToAesKeyIv(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	data, err := aes.Encrypt(body, key, iv)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(data)+ed25519.PublicKeySize)
	n := copy(out, privateKey[ed25519.PublicKeySize:])
	_ = copy(out[n:], data)
	return out, nil
}

func keyPairsToAesKeyIv(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (key, iv []byte, err error) {
	var pri, pub mixin.Key
	copy(pub[:], publicKey)
	privateKeyToCurve25519(pri, privateKey)

	if !pub.CheckKey() {
		err = errors.New("public key is invalid")
		return
	}

	if !pri.CheckScalar() {
		err = errors.New("private key is invalid")
		return
	}

	var point edwards25519.ExtendedGroupElement
	var point2 edwards25519.ProjectiveGroupElement

	tmp := [32]byte(pub)
	point.FromBytes(&tmp)
	tmp = pri
	edwards25519.GeScalarMult(&point2, &tmp, &point)

	point2.ToBytes(&tmp)
	return tmp[:16], tmp[16:], nil
}

func privateKeyToCurve25519(curve25519Private [32]byte, privateKey ed25519.PrivateKey) {
	h := sha512.New()
	h.Write(privateKey.Seed())
	digest := h.Sum(nil)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(curve25519Private[:], digest)
}
