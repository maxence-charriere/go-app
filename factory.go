package app

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

type compoBuilder struct {
	mutex sync.Mutex
	types map[string]reflect.Type
}

func newCompoBuilder() *compoBuilder {
	return &compoBuilder{
		types: make(map[string]reflect.Type),
	}
}

func (f *compoBuilder) register(c Compo) (name string, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

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

	name = compoName(c)
	f.types[name] = v.Type()
	return name, nil

}

func (f *compoBuilder) isRegistered(name string) bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	_, ok := f.types[name]
	return ok
}

func (f *compoBuilder) new(name string) (Compo, error) {
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
