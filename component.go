package app

import (
	"html/template"
	"net/url"
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

// ZeroCompo is the type to redefine when writing an empty component.
// Every instances of an empty struct is given the same memory address, which
// causes problem for indexing components.
// ZeroCompo have a placeholder field to avoid that.
type ZeroCompo struct {
	placeholder byte
}

// CompoBuilder is the interface that describes a component factory.
type CompoBuilder interface {
	// Register registers component of type c into the builder.
	// Components must be registered to be used.
	// During a rendering, it allows to create components of same type as c when
	// a tag named like c is found.
	Register(c Component) error

	// New creates a component named n.
	New(n string) (Component, error)
}
