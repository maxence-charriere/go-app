package app

import "github.com/murlokswarm/markup"

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
	markup.Register(c)
}
