package app

// func TestContext(t *testing.T) {
// 	ctx := &testContext{
// 		id: uid.Context(),
// 	}
// 	RegisterContext(ctx)
// 	defer UnregisterContext(ctx)

// 	compo := &Hello{}
// 	ctx.Mount(compo)
// 	ctxBis := Context(compo)
// 	if ctx != ctxBis {
// 		t.Error("ctx and ctx bis should be equals")
// 	}
// }

// func TestContextPanic(t *testing.T) {
// 	defer func() { recover() }()

// 	ctx := &testContext{
// 		id: uid.Context(),
// 	}
// 	RegisterContext(ctx)

// 	compo := &Hello{}
// 	ctx.Mount(compo)
// 	UnregisterContext(ctx)
// 	Context(compo)
// 	t.Error("should panic")
// }

// func TestContextByID(t *testing.T) {
// 	ctx := newTestContext("TestContextByID")
// 	defer ctx.Close()

// 	ctxBis, err := ContextByID(ctx.ID())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if ctx != ctxBis {
// 		t.Error("ctx and ctxBis should be the same context")
// 	}

// 	if _, err := ContextByID("Ctx-42"); err == nil {
// 		t.Error("err should not be nil")
// 	}
// }

// func TestRegisterContext(t *testing.T) {
// 	ctx := &testContext{
// 		id: uid.Context(),
// 	}
// 	RegisterContext(ctx)
// 	defer UnregisterContext(ctx)

// 	ctxBis, registered := contexts[ctx.ID()]
// 	if !registered {
// 		t.Fatal("ctx should be registered")
// 	}
// 	if ctxBis != ctx {
// 		t.Error("ctxBis and ctx should be equal")
// 	}
// }

// func TestRegisterContextNoID(t *testing.T) {
// 	defer func() { recover() }()

// 	ctx := &testContext{}
// 	RegisterContext(ctx)
// 	t.Error("should panic")
// }

// func TestRegisterContextAlreadyRegistered(t *testing.T) {
// 	defer func() { recover() }()

// 	ctx := newTestContext("context test")
// 	RegisterContext(ctx)
// 	RegisterContext(ctx)
// 	t.Error("should panic")
// }
