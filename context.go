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
	Mount(c Componer)

	// Renders an element.
	Render(elem *markup.Element)
}

// Context returns the context of c.
// c must be mounted.
func Context(c Componer) (ctx Contexter, err error) {
	root, err := markup.ComponentRoot(c)
	if err != nil {
		return
	}

	ctx, err = ContextByID(root.ContextID)
	return
}

// ContextByID returns the context registered under id.
func ContextByID(id uid.ID) (ctx Contexter, err error) {
	var registered bool

	if ctx, registered = contexts[id]; !registered {
		err = fmt.Errorf("context %v is not registered or has been closed", id)
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
	root        Componer
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
func (c *ZeroContext) Mount(component Componer) {
	markup.Mount(component, c.ID())
}

// Render is a placeholder method to satisfy the Contexter interface.
// It does nothing.
func (c *ZeroContext) Render(elem *markup.Element) {
	log.Info(elem.HTML())
}

// Close is a closes the context.
func (c *ZeroContext) Close() {
	markup.Dismount(c.root)
	UnregisterContext(c)
}
