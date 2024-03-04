package app

import (
	"encoding/json"
	"sync"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// BrowserStorage is the interface that describes a web browser storage.
type BrowserStorage interface {
	// Set sets the value to the given key. The value must be json convertible.
	Set(k string, v any) error

	// Get gets the item associated to the given key and store it in the given
	// value.
	// It returns an error if v is not a pointer.
	Get(k string, v any) error

	// Del deletes the item associated with the given key.
	Del(k string)

	// Len returns the number of items stored.
	Len() int

	// Clear deletes all items.
	Clear()

	// ForEach iterates over each item in the storage. For each item, it calls
	// the provided function f with the key of the item as its argument.
	// The order in which the items are processed is not specified and may vary
	// across different implementations of the BrowserStorage interface.
	ForEach(f func(key string))
}

type memoryStorage struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func newMemoryStorage() *memoryStorage {
	return &memoryStorage{
		data: make(map[string][]byte),
	}
}

func (s *memoryStorage) Set(k string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.data[k] = b
	s.mu.Unlock()
	return nil
}

func (s *memoryStorage) Get(k string, v any) error {
	s.mu.RLock()
	d, ok := s.data[k]
	if !ok {
		s.mu.RUnlock()
		return nil
	}

	s.mu.RUnlock()
	return json.Unmarshal(d, v)
}

func (s *memoryStorage) Del(k string) {
	s.mu.Lock()
	delete(s.data, k)
	s.mu.Unlock()
}

func (s *memoryStorage) Clear() {
	s.mu.Lock()
	for k := range s.data {
		delete(s.data, k)
	}
	s.mu.Unlock()
}

func (s *memoryStorage) Len() int {
	s.mu.RLock()
	l := len(s.data)
	s.mu.RUnlock()
	return l
}

func (s *memoryStorage) ForEach(f func(key string)) {
	for k := range s.data {
		f(k)
	}
}

type jsStorage struct {
	name  string
	mutex sync.RWMutex
}

func newJSStorage(name string) *jsStorage {
	return &jsStorage{name: name}
}

func (s *jsStorage) Set(k string, v any) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = errors.New("setting storage value failed").
				WithTag("storage-type", s.name).
				WithTag("key", k).
				Wrap(r.(error))
		}
	}()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	Window().Get(s.name).Call("setItem", k, string(b))
	return nil
}

func (s *jsStorage) Get(k string, v any) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	item := Window().Get(s.name).Call("getItem", k)
	if !item.Truthy() {
		return nil
	}

	return json.Unmarshal([]byte(item.String()), v)
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

func (s *jsStorage) Len() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.len()
}

func (s *jsStorage) len() int {
	return Window().Get(s.name).Get("length").Int()
}

func (s *jsStorage) ForEach(f func(key string)) {
	s.mutex.Lock()
	length := s.len()
	keys := make(map[string]struct{}, length)
	for i := 0; i < length; i++ {
		key := Window().Get(s.name).Call("key", i)
		if key.Truthy() {
			keys[key.String()] = struct{}{}
		}
	}
	s.mutex.Unlock()

	for key := range keys {
		f(key)
	}
}
