package app

import (
	"reflect"

	"github.com/murlokswarm/markup"
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

// RegisterComponent allows the app to create a component of type c when found
// into a markup.
// Should be called in an init func following the component implementation.
func RegisterComponent(c Componer) {
	v := reflect.Indirect(reflect.ValueOf(c))
	t := v.Type()

	constructor := func() markup.Componer {
		compo, _ := reflect.New(t).Interface().(Componer)
		return compo
	}

	markup.RegisterComponent(t.Name(), constructor)
}

// RegisterComponentWithConstructor allows the app to create the component
// returned by h when found into a markup.
// Should be called in an init func following the component implementation.
func RegisterComponentWithConstructor(h func() Componer) {
	v := reflect.Indirect(reflect.ValueOf(h()))
	t := v.Type()

	constructor := func() markup.Componer {
		return h()
	}

	markup.RegisterComponent(t.Name(), constructor)
}
