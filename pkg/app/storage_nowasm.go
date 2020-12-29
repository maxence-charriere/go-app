// +build !wasm

package app

import (
	"encoding/json"
	"fmt"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

func init() {
	LocalStorage = make(memoryStorage)
	SessionStorage = make(memoryStorage)
}

type memoryStorage map[string][]byte

func (s memoryStorage) Set(k string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	s[k] = b
	return nil
}

func (s memoryStorage) Get(k string, v interface{}) error {
	if _, ok := s[k]; !ok {
		return nil
	}
	return json.Unmarshal(s[k], v)
}

func (s memoryStorage) Del(k string) {
	delete(s, k)
}

func (s memoryStorage) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s memoryStorage) Len() int {
	return len(s)
}

func (s memoryStorage) Key(i int) (string, error) {
	j := 0
	for k := range s {
		fmt.Println(k)

		if i == j {
			return k, nil
		}
		j++
	}

	return "", errors.New("index out of range").
		Tag("index", i).
		Tag("len", s.Len())
}
