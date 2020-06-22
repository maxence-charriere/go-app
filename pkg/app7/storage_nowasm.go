// +build !wasm

package app

import "encoding/json"

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
