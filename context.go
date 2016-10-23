package app

import (
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

// Contexter represents the support where a component can be mounted.
// eg a window.
type Contexter interface {
	// The ID of the context.
	ID() uid.ID

	// Mounts the component and renders it in the context.
	Mount(c markup.Componer)

	// If applicable, resizes the context.
	Resize(width float64, height float64)

	// If applicable, moves the context.
	Move(x float64, y float64)
}
