package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

func encrypt(key string, v []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, v, nil), nil
}

func decrypt(key string, v []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(v) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := v[:nonceSize], v[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
