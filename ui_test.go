package app

import (
	"testing"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/satori/go.uuid"
)

type testContext struct {
	id          uuid.UUID
	placeholder string
	root        Componer
}

func newTestContext(placeholder string) (ctx *testContext) {
	ctx = &testContext{
		id:          uuid.NewV1(),
		placeholder: placeholder,
	}
	Elements().Add(ctx)
	return ctx
}

func (ctx *testContext) ID() uuid.UUID {
	return ctx.id
}

func (ctx *testContext) Close() error {
	markup.Dismount(ctx.root)
	Elements().Remove(ctx)
	return nil
}

func (ctx *testContext) Mount(c Componer) {
	markup.Mount(c, ctx.ID())
	ctx.root = c
}

func (ctx *testContext) Component() Componer {
	return ctx.root
}

func (ctx *testContext) Render(s markup.Sync) {
	log.Infof("rendering: %v", s.Node.Markup())
}

func TestElementStore(t *testing.T) {
	elems := newElementStore()
	ctx := &testContext{
		id: uuid.NewV1(),
	}

	elems.Add(ctx)
	if l := elems.Len(); l != 1 {
		t.Fatal("elems should have 1 element:", l)
	}

	if ctxbis, ok := elems.Get(ctx.ID()); !ok || ctx != ctxbis {
		t.Fatal("ctx and ctxbis should be equals")
	}

	elems.Remove(ctx)
	if l := elems.Len(); l != 0 {
		t.Fatal("elems should be empty:", l)
	}
}

func TestElementStoreAddPanic(t *testing.T) {
	defer func() { recover() }()

	elems := newElementStore()
	ctx := &testContext{
		id: uuid.NewV1(),
	}
	elems.Add(ctx)
	elems.Add(ctx)

	t.Error("should panic")
}

func TestContext(t *testing.T) {
	ctx := newTestContext("test ctx")
	defer ctx.Close()

	c := &Hello{}
	ctx.Mount(c)

	if ctxbis := Context(c); ctxbis != ctx {
		t.Error("ctx and ctxbis should be equal")
	}
}

func TestContextPanic(t *testing.T) {
	defer func() { recover() }()

	ctx := &testContext{}
	c := &Hello{}
	ctx.Mount(c)
	Context(c)

	t.Error("should panic")
}

func TestUIChan(t *testing.T) {
	doneChan := make(chan interface{})
	defer close(doneChan)

	UIChan <- func() {
		doneChan <- true
	}

	<-doneChan
}
