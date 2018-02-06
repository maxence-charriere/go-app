package app

import (
	"reflect"

	"github.com/pkg/errors"
)

// Factory is the interface that describes a component factory.
type Factory interface {
	// RegisterComponent registers a component under its type name lowercased.
	RegisterComponent(c Component) (name string, err error)

	// IsRegisteredComponent reports wheter the named component is registered.
	IsRegisteredComponent(name string) bool

	// NewComponent creates the named component.
	// It returns an error if there is no component registered under name.
	NewComponent(name string) (Component, error)
}

// NewFactory creates a component factory.
func NewFactory() Factory {
	return make(factory)
}

// A factory that implements the Factory interface.
type factory map[string]reflect.Type

func (f factory) RegisterComponent(c Component) (name string, err error) {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return name, errors.Errorf("%T must be implemented on a struct pointer", c)
	}

	if v = v.Elem(); v.Kind() != reflect.Struct {
		return name, errors.Errorf("%T must be implemented on a struct pointer", c)
	}

	if v.NumField() == 0 {
		return name, errors.Errorf("%T can't be implemented on an empty struct pointer", c)
	}

	v = reflect.ValueOf(c).Elem()
	t := v.Type()

	name = normalizeComponentName(t.String())
	f[name] = t
	return
}

func (f factory) IsRegisteredComponent(name string) bool {
	_, registered := f[name]
	return registered
}

func (f factory) NewComponent(name string) (c Component, err error) {
	t, ok := f[name]
	if !ok {
		err = errors.Errorf("component %s is not registered", name)
		return
	}

	v := reflect.New(t)
	c = v.Interface().(Component)
	return
}
