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

// NewFactory creates a component factory.
func NewFactory() Factory {
	return make(factory)
}

type factory map[string]reflect.Type

func (f factory) Register(c Component) (name string, err error) {
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
	f[name] = rtype
	return name, nil
}

func (f factory) Registered(name string) bool {
	_, ok := f[name]
	return ok
}

func (f factory) New(name string) (Component, error) {
	rtype, ok := f[name]
	if !ok {
		return nil, errors.Errorf("component %s is not registered", name)
	}

	rval := reflect.New(rtype)

	// Here we are not checking the cast because only component cant go in the
	// factory.
	c := rval.Interface().(Component)
	return c, nil
}

// NewFactoryWithLogs returns a decorated version of the given factory that logs
// all the operations.
// It uses the default logger.
func NewFactoryWithLogs(f Factory) Factory {
	return &factoryWithLogs{
		base: f,
	}
}

type factoryWithLogs struct {
	base Factory
}

func (f *factoryWithLogs) Register(c Component) (name string, err error) {
	name, err = f.base.Register(c)
	if err != nil {
		Errorf("registering %T failed: %s", c, err)
		return "", err
	}

	Logf("%T are registered as %s", c, name)
	return name, nil
}

func (f *factoryWithLogs) Registered(name string) bool {
	ok := f.base.Registered(name)
	Log(name, "registered:", ok)
	return ok
}

func (f *factoryWithLogs) New(name string) (Component, error) {
	Log("creating component", name)

	c, err := f.base.New(name)
	if err != nil {
		Errorf("creating component %s failed: %s", name, err)
	}
	return c, err
}

// NewConcurrentFactory returns a decorated version of the given factory that
// is safe for concurrent operations.
func NewConcurrentFactory(f Factory) Factory {
	return &concurrentFactory{
		base: f,
	}
}

type concurrentFactory struct {
	mutex sync.Mutex
	base  Factory
}

func (f *concurrentFactory) Register(c Component) (name string, err error) {
	f.mutex.Lock()
	name, err = f.base.Register(c)
	f.mutex.Unlock()
	return name, err
}

func (f *concurrentFactory) Registered(name string) bool {
	f.mutex.Lock()
	ok := f.base.Registered(name)
	f.mutex.Unlock()
	return ok
}

func (f *concurrentFactory) New(name string) (Component, error) {
	f.mutex.Lock()
	c, err := f.base.New(name)
	f.mutex.Unlock()
	return c, err
}
