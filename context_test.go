package app

import (
	"testing"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

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
	c.root = component
}

// Render is a placeholder method to satisfy the Contexter interface.
// It does nothing.
func (c *ZeroContext) Render(s markup.Sync) {
	log.Infof("%v rendering: %v", s.Scope, s.Node.Markup())
}

// Close is a closes the context.
func (c *ZeroContext) Close() {
	markup.Dismount(c.root)
	UnregisterContext(c)
}

func TestContext(t *testing.T) {
	ctx := &ZeroContext{
		id: uid.Context(),
	}
	RegisterContext(ctx)
	defer UnregisterContext(ctx)

	compo := &Hello{}
	ctx.Mount(compo)
	ctxBis := Context(compo)
	if ctx != ctxBis {
		t.Error("ctx and ctx bis should be equals")
	}
}

func TestContextPanic(t *testing.T) {
	defer func() { recover() }()

	ctx := &ZeroContext{
		id: uid.Context(),
	}
	RegisterContext(ctx)

	compo := &Hello{}
	ctx.Mount(compo)
	UnregisterContext(ctx)
	Context(compo)
	t.Error("should panic")
}

func TestContextByID(t *testing.T) {
	ctx := NewZeroContext("TestContextByID")
	defer ctx.Close()

	ctxBis, err := ContextByID(ctx.ID())
	if err != nil {
		t.Fatal(err)
	}
	if ctx != ctxBis {
		t.Error("ctx and ctxBis should be the same context")
	}

	if _, err := ContextByID("Ctx-42"); err == nil {
		t.Error("err should not be nil")
	}
}

func TestRegisterContext(t *testing.T) {
	ctx := &ZeroContext{
		id: uid.Context(),
	}
	RegisterContext(ctx)
	defer UnregisterContext(ctx)

	ctxBis, registered := contexts[ctx.ID()]
	if !registered {
		t.Fatal("ctx should be registered")
	}
	if ctxBis != ctx {
		t.Error("ctxBis and ctx should be equal")
	}
}

func TestRegisterContextNoID(t *testing.T) {
	defer func() { recover() }()

	ctx := &ZeroContext{}
	RegisterContext(ctx)
	t.Error("should panic")
}

func TestRegisterContextAlreadyRegistered(t *testing.T) {
	defer func() { recover() }()

	ctx := NewZeroContext("context test")
	RegisterContext(ctx)
	RegisterContext(ctx)
	t.Error("should panic")
}
