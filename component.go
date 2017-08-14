package app

import (
	"github.com/murlokswarm/markup"
	"github.com/satori/go.uuid"
)

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
