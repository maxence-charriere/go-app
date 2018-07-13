package test

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app/internal/core"
)

// A Element implementation for tests.
type Element struct {
	core.Elem
	id uuid.UUID
}

// NewElement creates a new element.
func NewElement(d *Driver) *Element {
	elem := &Element{
		id: uuid.New(),
	}

	d.elems.Put(elem)
	return elem
}

// ID satisfies the app.Element interface.
func (e *Element) ID() uuid.UUID {
	return e.id
}
