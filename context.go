package app

import (
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

	// If applicable, moves the context.
	Move(x float64, y float64)

	// If applicable, resizes the context.
	Resize(width float64, height float64)

	// If applicable, set the icon targeted by path.
	SetIcon(path string)
}

func registerContext(c Contexter) {
	if len(c.ID()) == 0 {
		log.Panicf("context %T is invalid. ID must be set", c)
	}

	if _, registered := contexts[c.ID()]; registered {
		log.Panicf("context %T with id %v is already registered", c, c.ID())
	}

	contexts[c.ID()] = c
}

func unregisterContext(c Contexter) {
	delete(contexts, c.ID())
}

// ZeroContext is a placeholder context. It's used to give a support on non
// implemented or not available app components.
// eg There is an app menu on MacOS, not on Windows.
//
// Use of methods from a ZeroContext doesn't do anything.
type ZeroContext struct {
	id          uid.ID
	placeholder string
}

// NewZeroContext creates a ZeroContext.
func NewZeroContext(placeholder string) (ctx *ZeroContext) {
	ctx = &ZeroContext{
		id:          uid.Context(),
		placeholder: placeholder,
	}

	registerContext(ctx)
	return
}

// ID returns the ID of the context.
func (c *ZeroContext) ID() uid.ID {
	return c.id
}

// Mount is a placeholder method to satisfy the Contexter interface.
// It does nothing.
func (c *ZeroContext) Mount(component markup.Componer) {
	log.Infof("%T is mounted into %v (%v)", component, c.placeholder, c.ID())
}

// Resize is a placeholder method to satisfy the Contexter interface.
// It does nothing.
func (c *ZeroContext) Resize(width float64, height float64) {
	log.Infof("%v (%v) simulates a resize of %v x %v", c.placeholder, c.ID(), width, height)
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
