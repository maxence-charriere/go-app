package app

import (
	"html/template"
	"net/url"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// Component is the interface that describes a component.
// Should be implemented on a non empty struct pointer.
type Component interface {
	// Render should return a string describing the component with HTML5
	// standard.
	// It supports standard Go html/template API.
	// Pipeline is based on the component struct.
	// See https://golang.org/pkg/text/template and
	// https://golang.org/pkg/html/template for template usage.
	Render() string
}

// Mounter is the interface that wraps OnMount method.
// OnMount si called when a component is mounted.
type Mounter interface {
	OnMount()
}

// Dismounter is the interface that wraps OnDismount method.
// OnDismount si called when a component is dismounted.
type Dismounter interface {
	OnDismount()
}

// Navigator is the interface that wraps OnNavigate method.
// OnNavigate is called when a component is navigated to.
type Navigator interface {
	OnNavigate(u url.URL)
}

// Mapper is the interface that wraps FuncMaps method.
type Mapper interface {
	// Allows to add custom functions to the template used to render the
	// component.
	//
	// Funcs named raw, json and time are reserved. They handle raw html code,
	// json conversions and time format.
	// They can't be overloaded.
	// See https://golang.org/pkg/text/template/#Template.Funcs for more details.
	FuncMaps() template.FuncMap
}

// ZeroCompo is the type to use as base for an empty compone.
// Every instances of an empty struct is given the same memory address, which
// causes problem for indexing components.
// ZeroCompo have a placeholder field to avoid that.
type ZeroCompo struct {
	placeholder byte
}

// Factory is the interface that describes a component factory.
type Factory interface {
	// RegisterComponent registers a component under its type name lowercased.
	RegisterComponent(c Component) (name string, err error)

	// NewComponent creates the named component.
	// It returns an error if there is no component registered under name.
	NewComponent(name string) (Component, error)
}

// A factory that implements the Factory interface.
type factory map[string]reflect.Type

func (f factory) RegisterComponent(c Component) (name string, err error) {
	if err = ensureValidComponent(c); err != nil {
		return
	}

	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	name = normalizeComponentName(t.String())
	f[name] = t
	return
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

func ensureValidComponent(c Component) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return errors.Errorf("%T must be implemented on a struct pointer", c)
	}

	if v = v.Elem(); v.Kind() != reflect.Struct {
		return errors.Errorf("%T must be implemented on a struct pointer", c)
	}

	if v.NumField() == 0 {
		return errors.Errorf("%T can't be implemented on an empty struct pointer", c)
	}
	return nil
}

func normalizeComponentName(name string) string {
	name = strings.ToLower(name)
	if pkgsep := strings.IndexByte(name, '.'); pkgsep != -1 {
		pkgname := name[:pkgsep]
		if pkgname == "main" {
			name = name[pkgsep+1:]
		}
	}
	return name
}
