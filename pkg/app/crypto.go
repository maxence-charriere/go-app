package app

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
)

func encrypt(key string, v []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithRandomNonce(c)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nil, nil, v, nil), nil
}

func decrypt(key string, v []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithRandomNonce(c)
	if err != nil {
		return nil, err
	}

	// Preserve the legacy error for an incomplete 12-byte nonce prefix.
	// NewGCMWithRandomNonce reports a NonceSize of zero.
	if len(v) < 12 {
		return nil, errors.New("ciphertext too short")
	}

	return gcm.Open(nil, nil, v, nil)
}
