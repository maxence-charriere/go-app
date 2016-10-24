package app

import (
	"testing"

	"github.com/murlokswarm/uid"
)

func TestRegisterContext(t *testing.T) {
	ctx := &ZeroContext{
		id: uid.Context(),
	}

	registerContext(ctx)
	unregisterContext(ctx)
}

func TestRegisterContextNoID(t *testing.T) {
	defer func() { recover() }()

	ctx := &ZeroContext{}
	registerContext(ctx)
	t.Error("should panic")
}

func TestRegisterContextAlreadyRegistered(t *testing.T) {
	defer func() { recover() }()

	ctx := NewZeroContext("context test")
	registerContext(ctx)
	registerContext(ctx)
	t.Error("should panic")
}

func TestZeroContext(t *testing.T) {
	ctx := NewZeroContext("context test")
	defer unregisterContext(ctx)

	t.Log(ctx.ID())
	ctx.Resize(42, 42)
	ctx.Move(42, 42)
	ctx.SetIcon("test.png")

	hello := &Hello{}
	ctx.Mount(hello)
}
