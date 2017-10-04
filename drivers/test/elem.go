package test

import "github.com/google/uuid"

// A Element implementation for tests.
type Element struct {
	id uuid.UUID
}

// NewElement creates a new element.
func NewElement(d *Driver) *Element {
	elem := &Element{
		id: uuid.New(),
	}
	d.elements.Add(elem)
	return elem
}

// ID satisfies the app.Element interface.
func (e *Element) ID() uuid.UUID {
	return e.id
}
