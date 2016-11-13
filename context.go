package app

import (
	"fmt"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

var (
	contexts = map[uid.ID]Contexter{}
)

// Contexter represents the support where a component can be mounted.
// eg a window.
type Contexter interface {
	// The ID of the context.
	ID() uid.ID

	// Mounts the component and renders it in the context.
	Mount(c markup.Componer)

	// Renders an element.
	Render(elem *markup.Element)

	// If applicable, return the position of the context.
	Position() (x float64, y float64)

	// If applicable, moves the context.
	Move(x float64, y float64)

	// If applicable, return the size of the context.
	Size() (width float64, height float64)

	// If applicable, resizes the context.
	Resize(width float64, height float64)

	// If applicable, set the icon targeted by path.
	SetIcon(path string)

	// Close the context.
	// Should call markup.Dismount on its root component.
	// Should call UnregisterContext on itself.
	// Should perform additional cleanup if required.
	Close()
}

// Context returns the context of c.
// c must be mounted.
func Context(c markup.Componer) (ctx Contexter, err error) {
	var root *markup.Element

	if root, err = markup.ComponentRoot(c); err != nil {
		return
	}

	ctx, registered := contexts[root.ContextID]
	if !registered {
		err = fmt.Errorf("context %v has been closed", root.ID)
	}

	return
}

// RegisterContext registers c.
// Should be used only in a driver implementation.
func RegisterContext(c Contexter) {
	if len(c.ID()) == 0 {
		log.Panicf("context %T is invalid. ID must be set", c)
	}

	if _, registered := contexts[c.ID()]; registered {
		log.Panicf("context %T with id %v is already registered", c, c.ID())
	}

	contexts[c.ID()] = c
}

// UnregisterContext unregisters c.
// Should be used only in a driver implementation.
func UnregisterContext(c Contexter) {
	delete(contexts, c.ID())
}

// ZeroContext is a placeholder context.
// It's used as a replacement for non available or non implemented features.
//
// Use of methods from a ZeroContext doesn't do anything.
type ZeroContext struct {
	id          uid.ID
	placeholder string
	root        markup.Componer
}

// NewZeroContext creates a ZeroContext.
func NewZeroContext(placeholder string) (ctx *ZeroContext) {
	ctx = &ZeroContext{
		id:          uid.Context(),
		placeholder: placeholder,
	}

	RegisterContext(ctx)
	return
}

// ID returns the ID of the context.
func (c *ZeroContext) ID() uid.ID {
	return c.id
}

// Mount is a placeholder method to satisfy the Contexter interface.
// It does nothing.
func (c *ZeroContext) Mount(component markup.Componer) {
	markup.Mount(component, c.ID())
	log.Infof("%T is mounted into %v (%v)", component, c.placeholder, c.ID())
}

// Render is a placeholder method to satisfy the Contexter interface.
// It does nothing.
func (c *ZeroContext) Render(elem *markup.Element) {
	log.Infof("rendering:\n\033[32m%v\033[00m", elem.HTML())
}

// Size is a placeholder method to satisfy the Contexter interface.
func (c *ZeroContext) Size() (width float64, height float64) {
	return
}

// Resize is a placeholder method to satisfy the Contexter interface.
// It does nothing.
func (c *ZeroContext) Resize(width float64, height float64) {
	log.Infof("%v (%v) simulates a resize of %v x %v", c.placeholder, c.ID(), width, height)
}

// Position is a placeholder method to satisfy the Contexter interface.
func (c *ZeroContext) Position() (x float64, y float64) {
	return
}

// Move is a placeholder method to satisfy the Contexter interface.
// It does nothing.
func (c *ZeroContext) Move(x float64, y float64) {
	log.Infof("%v (%v) simulates a move to (%v, %v)", c.placeholder, c.ID(), x, y)
}

// SetIcon is a placeholder method to satisfy the Contexter interface.
// It does nothing.
func (c *ZeroContext) SetIcon(path string) {
	log.Infof("%v (%v) simulates set icon with %v", c.placeholder, c.ID(), path)
}

// Close is a closes the context.
func (c *ZeroContext) Close() {
	markup.Dismount(c.root)
	UnregisterContext(c)
	log.Infof("%v (%v) is closed", c.placeholder, c.ID())
}
