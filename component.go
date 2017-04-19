package app

import (
	"net/url"
	"text/template"

	"github.com/murlokswarm/markup"
	"github.com/satori/go.uuid"
)

// Componer is the interface that describes a component.
type Componer interface {
	// Render should returns a markup.
	// The markup can be a template string following the text/template standard
	// package rules.
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

// Hrefer is the interface that wraps OnHref method.
// OnHref is called when the component is mounted from a click on a link that
// targets the component in the href attribute.
type Hrefer interface {
	OnHref(URL *url.URL)
}

// TemplateFuncMapper is the interface that wraps FuncMaps method.
type TemplateFuncMapper interface {
	// Allows to add custom functions to the template used to render the
	// component.
	// Note that funcs named json and time are already implemented to handle
	// structs as prop and time format. Overloads of these will be ignored.
	// See https://golang.org/pkg/text/template/#Template.Funcs for more details.
	FuncMaps() template.FuncMap
}

// RegisterComponent allows the app to create a component of type c when found
// into a markup.
// Should be called in an init func following the component implementation.
func RegisterComponent(c Componer) {
	markup.Register(c)
}

// ComponentID returns the id of c. Panic if c is not mounted.
func ComponentID(c Componer) uuid.UUID {
	return markup.ID(c)
}

// ComponentByID returns the component associated with id.
// Panic if no component with id is mounted.
func ComponentByID(id uuid.UUID) Componer {
	return markup.Component(id)
}
