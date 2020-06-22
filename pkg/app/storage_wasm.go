package app

import (
	"encoding/json"
	"sync"
	"syscall/js"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

func init() {
	LocalStorage = newJSStorage("localStorage")
	SessionStorage = newJSStorage("sessionStorage")
}

type jsStorage struct {
	name  string
	mutex sync.RWMutex
}

func newJSStorage(name string) *jsStorage {
	return &jsStorage{name: name}
}

func (s *jsStorage) Set(k string, v interface{}) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = errors.New("setting storage value failed").
				Tag("storage-type", s.name).
				Tag("key", k).
				Wrap(r.(js.Error))
		}
	}()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	Window().Get(s.name).Call("setItem", k, btos(b))
	return nil
}

func (s *jsStorage) Get(k string, v interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	item := Window().Get(s.name).Call("getItem", k)
	if !item.Truthy() {
		return nil
	}

	return json.Unmarshal(stob(item.String()), v)
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
