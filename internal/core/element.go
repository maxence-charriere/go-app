package core

import "github.com/murlokswarm/app"

// ElementWithComponent is a base struct to embed in ElementWithComponent
// implementations.
type ElementWithComponent struct{}

// WhenWindow satisfies the app.ElementWithComponent interface.
func (e *ElementWithComponent) WhenWindow(func(w app.Window)) {}

// WhenPage satisfies the app.ElementWithComponent interface.
func (e *ElementWithComponent) WhenPage(func(p app.Page)) {}

// WhenDockTile satisfies the app.ElementWithComponent interface.
func (e *ElementWithComponent) WhenDockTile(func(d app.DockTile)) {}

// WhenStatusMenu satisfies the app.ElementWithComponent interface.
func (e *ElementWithComponent) WhenStatusMenu(func(s app.StatusMenu)) {}
