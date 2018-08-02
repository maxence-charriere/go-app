package app

import (
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// NewFactory creates a component factory.
func NewFactory() *Factory {
	return &Factory{
		types: make(map[string]reflect.Type),
	}
}

// Factory represents a factory that creates components.
// It is safe for concurrent operations.
type Factory struct {
	mutex sync.Mutex
	types map[string]reflect.Type
}

// RegisterCompo registers the given component.
func (f *Factory) RegisterCompo(c Compo) (name string, err error) {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return "", errors.New("component is not a pointer")
	}

	if v = v.Elem(); v.Kind() != reflect.Struct {
		return "", errors.New("component is not based on a struct")
	}

	if v.NumField() == 0 {
		return "", errors.New("component does not have fields")
	}

	t := v.Type()
	name = strings.ToLower(t.String())
	name = strings.TrimPrefix(name, "main.")
	f.types[name] = t
	return name, nil

}

// IsCompoRegistered reports whether the named component is registered.
func (f *Factory) IsCompoRegistered(name string) bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	_, ok := f.types[name]
	return ok
}

// NewCompo creates the named component.
func (f *Factory) NewCompo(name string) (Compo, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	t, ok := f.types[name]
	if !ok {
		return nil, errors.Errorf("component %s is not registered", name)
	}

	v := reflect.New(t)
	c := v.Interface().(Compo)

	return c, nil
}
