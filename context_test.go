package app

import (
	"testing"

	"github.com/murlokswarm/uid"
)

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

func TestContextByID(t *testing.T) {
	ctx := NewZeroContext("TestContextByID")

	ctxBis, err := ContextByID(ctx.ID())
	if err != nil {
		t.Error(err)
	}

	if ctx != ctxBis {
		t.Error("ctx and ctxBis should be the same context")
	}

	if ctxBis, err = ContextByID("Ctx-42"); err == nil {
		t.Error("should error")
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

func TestZeroContext(t *testing.T) {
	ctx := NewZeroContext("context test")
	defer UnregisterContext(ctx)

	t.Log(ctx.ID())

	hello := &Hello{}
	ctx.Mount(hello)
}
