package app

import (
	"github.com/murlokswarm/errors"
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

	// Renders a node following desription described in s.
	Render(s markup.Sync)
}

// Context returns the context of c.
// Panic if c is not mounted.
func Context(c Componer) Contexter {
	root := markup.Root(c)
	return ContextByID(root.ContextID)
}

// ContextByID returns the context registered under id.
// Panic if context is not registered.
func ContextByID(id uid.ID) (ctx Contexter) {
	var registered bool
	if ctx, registered = contexts[id]; !registered {
		err := errors.Newf("context %v is not registered or has been closed", id)
		log.Panic(err)
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
