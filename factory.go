package app

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// Factory is the interface that describes a component factory.
type Factory interface {
	// Register registers the given component under its type name lowercased.
	Register(c Component) (name string, err error)

	// Registered reports wheter the named component is registered.
	Registered(name string) bool

	// New creates the named component.
	New(name string) (Component, error)
}

// NewFactory creates a component factory that is safe for concurrent use.
func NewFactory() Factory {
	return &factory{
		types: make(map[string]reflect.Type),
	}
}

type factory struct {
	mutex sync.RWMutex
	types map[string]reflect.Type
}

func (f *factory) Register(c Component) (name string, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	rval := reflect.ValueOf(c)
	if rval.Kind() != reflect.Ptr {
		return "", errors.New("component is not a pointer")
	}

	if rval = rval.Elem(); rval.Kind() != reflect.Struct {
		return "", errors.New("component does not point to a struct")
	}

	if rval.NumField() == 0 {
		return "", errors.New("component does not have field")
	}

	rtype := rval.Type()
	name = normalizeComponentName(rtype.String())
	f.types[name] = rtype
	return name, nil
}

func (f *factory) Registered(name string) bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	_, ok := f.types[name]
	return ok
}

func (f *factory) New(name string) (Component, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	rtype, ok := f.types[name]
	if !ok {
		return nil, errors.Errorf("component %s is not registered", name)
	}

	rval := reflect.New(rtype)

	// Here we are not checking the cast because only component can go in the
	// factory.
	c := rval.Interface().(Component)
	return c, nil
}
