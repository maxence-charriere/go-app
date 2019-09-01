package app

import (
	"errors"
	"net/url"
	"reflect"
	"strings"
)

// Compo is the interface that describes a component.
// Must be implemented on a non empty struct pointer.
type Compo interface {
	// Render must return a HTML5 string. It supports standard Go html/template
	// API. The pipeline is based on the component struct. See
	// https://golang.org/pkg/text/template and
	// https://golang.org/pkg/html/template for template usage.
	Render() string
}

// Mounter is the interface that wraps OnMount method.
type Mounter interface {
	// OnMount is called when a component is mounted.
	// App.Render should not be called inside.
	OnMount()
}

// Dismounter is the interface that wraps OnDismount method.
type Dismounter interface {
	// OnDismount is called when a component is dismounted.
	// App.Render should not be called inside.
	OnDismount()
}

// CompoWithExtendedRender is the interface that wraps Funcs method.
type CompoWithExtendedRender interface {
	// Funcs returns a map of funcs to use when rendering a component.
	// Funcs named raw, json and time are reserved.
	// They handle raw html code, json conversions and time format.
	// They can't be overloaded.
	// See https://golang.org/pkg/text/template/#Template.Funcs for more details.
	Funcs() map[string]interface{}
}

// ZeroCompo is the type to use as base for empty components.
// Every instances of an empty struct is given the same memory address, which
// causes problem for indexing components.
// ZeroCompo have a placeholder field to avoid that.
type ZeroCompo struct {
	placeholder byte
}

type compoBuilder map[string]reflect.Type

func (b compoBuilder) imports(c Compo) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return errors.New("component is not a pointer")
	}
	if v = v.Elem(); v.Kind() != reflect.Struct {
		return errors.New("component is not implemented on a struct")
	}
	if v.NumField() == 0 {
		return errors.New("component is based on a struct without field. use ZeroCompo instead of struct{}")
	}

	b[compoName(c)] = v.Type()
	return nil
}

func (b compoBuilder) isImported(name string) bool {
	_, ok := b[name]
	return ok
}

func (b compoBuilder) new(name string) (Compo, error) {
	t, ok := b[name]
	if !ok {
		return nil, errors.New("component " + name + " is not imported")
	}
	return reflect.New(t).Interface().(Compo), nil
}

func compoName(c Compo) string {
	v := reflect.ValueOf(c)
	v = reflect.Indirect(v)
	name := strings.ToLower(v.Type().String())
	return strings.TrimPrefix(name, "main.")
}

func compoNameFromURLString(rawurl string) string {
	u, _ := url.Parse(rawurl)
	return compoNameFromURL(u)
}

func compoNameFromURL(u *url.URL) string {
	p := u.Path
	p = strings.TrimPrefix(p, "/")

	path := strings.SplitN(p, "/", 2)
	if len(path[0]) == 0 {
		return ""
	}

	names := strings.SplitN(path[0], "?", 2)
	name := names[0]
	name = strings.ToLower(name)
	return strings.TrimPrefix(name, "main.")
}
