package app

import (
	"net/url"
	"reflect"
	"strings"
)

// Compo is the interface that describes a component.
// Must be implemented on a non empty struct pointer.
type Compo interface {
	// Render must return HTML 5.
	// It supports standard Go html/template API.
	// The pipeline is based on the component struct.
	// See https://golang.org/pkg/text/template and
	// https://golang.org/pkg/html/template for template usage.
	Render() string
}

// Mounter is the interface that wraps OnMount method.
type Mounter interface {
	Compo

	// OnMount is called when a component is mounted.
	// App.Render should not be called inside.
	OnMount()
}

// Dismounter is the interface that wraps OnDismount method.
type Dismounter interface {
	Compo

	// OnDismount is called when a component is dismounted.
	// App.Render should not be called inside.
	OnDismount()
}

// Navigable is the interface that wraps OnNavigate method.
type Navigable interface {
	Compo

	// OnNavigate is called when a component is navigated to.
	OnNavigate(u *url.URL)
}

// CompoWithExtendedRender is the interface that wraps Funcs method.
type CompoWithExtendedRender interface {
	Compo

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
