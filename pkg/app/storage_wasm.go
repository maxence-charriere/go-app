package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"sync"
)

type jsStorage struct {
	name  string
	key   []byte
	mutex sync.RWMutex
}

func newJSStorage(name string) *jsStorage {
	u := Window().URL()

	key := []byte(u.Scheme + "(*_*)" + u.Host)
	for len(key) < 32 {
		key = append(key, 'o')
	}
	key = key[:32]

	return &jsStorage{
		name: name,
		key:  key,
	}
}

func (s *jsStorage) Set(k string, v interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	b, err = encrypt(b, s.key)
	if err != nil {
		return err
	}

	item := base64.StdEncoding.EncodeToString(b)
	Window().Get(s.name).Call("setItem", k, item)
	return nil
}

func (s *jsStorage) Get(k string, v interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	item := Window().Get(s.name).Call("getItem", k)
	if !item.Truthy() {
		return nil
	}

	b, err := base64.StdEncoding.DecodeString(item.String())
	if err != nil {
		return err
	}

	if b, err = decrypt(b, s.key); err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (s *jsStorage) Del(k string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	Window().Get(s.name).Call("removeItem", k)
}

func (s *jsStorage) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	Window().Get(s.name).Call("clear")
}

func encrypt(v []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
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

func decrypt(v []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
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
