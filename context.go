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

	// If applicable, resizes the context.
	Resize(width float64, height float64)

	// If applicable, moves the context.
	Move(x float64, y float64)

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

type ZeroContext struct {
	id          uid.ID
	placeholder string
}

func NewZeroContext(placeholder string) (ctx *ZeroContext) {
	ctx = &ZeroContext{
		id:          uid.Context(),
		placeholder: placeholder,
	}

	registerContext(ctx)
	return
}

func (c *ZeroContext) ID() uid.ID {
	return c.id
}

func (c *ZeroContext) Mount(component markup.Componer) {
	log.Infof("%T is mounted into %v (%v)", component, c.placeholder, c.ID())
}

func (c *ZeroContext) Resize(width float64, height float64) {
	log.Infof("%v (%v) simulates a resize of %v x %v", c.placeholder, c.ID(), width, height)
}

func (c *ZeroContext) Move(x float64, y float64) {
	log.Infof("%v (%v) simulates a move to (%v, %v)", c.placeholder, c.ID(), x, y)
}

func (c *ZeroContext) SetIcon(path string) {
	log.Infof("%v (%v) simulates set icon with %v", c.placeholder, c.ID(), path)
}
