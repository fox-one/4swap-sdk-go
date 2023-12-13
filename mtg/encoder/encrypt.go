package encoder

import (
	"bytes"
	"crypto/ed25519"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"io/ioutil"

	"github.com/fox-one/4swap-sdk-go/v2/mtg/aes"
	"golang.org/x/crypto/curve25519"
)

func Decrypt(b []byte, privateKey ed25519.PrivateKey) ([]byte, error) {
	r := bytes.NewReader(b)

	// read
	pub := make([]byte, 32)
	if _, err := io.ReadFull(r, pub); err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(r)

	key, iv, err := keyPairsToAesKeyIv(privateKey, pub)
	if err != nil {
		return nil, err
	}

	return decryptWithAseKeyIv(data, key, iv, md5.New())
}

func Encrypt(body []byte, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) ([]byte, error) {
	key, iv, err := keyPairsToAesKeyIv(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	prefix := privateKey.Public().(ed25519.PublicKey)
	return encryptWithAesKeyIv(body, prefix, key, iv, md5.New())
}

func keyPairsToAesKeyIv(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (key, iv []byte, err error) {
	var pri, pub [32]byte

	privateKeyToCurve25519(&pri, privateKey)
	if err = publicKeyToCurve25519(&pub, publicKey); err != nil {
		return
	}

	var dst []byte
	if dst, err = curve25519.X25519(pri[:], pub[:]); err != nil {
		return
	}

	if l := len(dst); l != 32 {
		err = fmt.Errorf("bad scalar multiplication length: %d, expected %d", l, 32)
		return
	}

	key, iv = dst, md5Hash(dst)
	return
}

func encryptWithAesKeyIv(body, prefix, key, iv []byte, h hash.Hash) ([]byte, error) {
	if h != nil {
		if _, err := h.Write(body); err != nil {
			return nil, err
		}

		body = append(h.Sum(nil), body...)
	}

	data, err := aes.Encrypt(body, key, iv)
	if err != nil {
		return nil, err
	}

	dst := make([]byte, len(prefix)+len(data))
	n := copy(dst, prefix)
	_ = copy(dst[n:], data)
	return dst, nil
}

func decryptWithAseKeyIv(data, key, iv []byte, h hash.Hash) ([]byte, error) {
	b, err := aes.Decrypt(data, key, iv)
	if err != nil {
		return nil, err
	}

	if h != nil {
		r := bytes.NewReader(b)
		sig := make([]byte, h.Size())
		if _, err := io.ReadFull(r, sig); err != nil {
			return nil, err
		}

		b, _ = ioutil.ReadAll(r)
		if _, err := h.Write(b); err != nil {
			return nil, err
		}

		if !bytes.Equal(h.Sum(nil), sig) {
			return nil, fmt.Errorf("invalid signature")
		}
	}

	return b, nil
}
