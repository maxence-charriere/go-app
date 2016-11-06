package app

import (
	"testing"

	"github.com/murlokswarm/uid"
)

func TestRegisterContext(t *testing.T) {
	ctx := &ZeroContext{
		id: uid.Context(),
	}

	RegisterContext(ctx)
	defer UnregisterContext(ctx)

	ctxBis, registered := Context(ctx.ID())
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
	ctx.Resize(42, 42)
	ctx.Move(42, 42)
	ctx.SetIcon("test.png")

	hello := &Hello{}
	ctx.Mount(hello)
}
