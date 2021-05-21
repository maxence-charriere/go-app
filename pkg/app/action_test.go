package app

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	Handle("/test", func(Context, Action) {})
	require.Len(t, actionHandlers, 1)
}

func TestActionManagerHandle(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	m := actionManager{}

	h := &hello{}
	e.Mount(h)
	e.Consume()

	isHandleACalled := false
	isHandleBCalled := false
	isHandleCCalled := false
	isHandleDCalled := false

	m.handle("/test", false, h, func(ctx Context, a Action) {
		isHandleACalled = true
	})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 1)

	m.handle("/test", false, h, func(ctx Context, a Action) {
		isHandleBCalled = true
	})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 2)

	f := &foo{}
	m.handle("/test", false, f, func(ctx Context, a Action) {
		isHandleCCalled = true
	})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 3)

	m.handle("/test", true, e.Body, func(ctx Context, a Action) {
		isHandleDCalled = true
	})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 4)

	m.post(Action{Name: "/test"})
	e.Consume()
	require.True(t, isHandleACalled)
	require.True(t, isHandleBCalled)
	require.False(t, isHandleCCalled)
	require.True(t, isHandleDCalled)
	require.Len(t, m.handlers["/test"], 3)
}

func TestActionManagerCloseUnusedHandlers(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	m := actionManager{}

	h := &hello{}
	e.Mount(h)
	e.Consume()

	m.handle("/test", false, h, func(ctx Context, a Action) {})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 1)

	f := &foo{}
	m.handle("/test", false, f, func(ctx Context, a Action) {})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 2)

	m.closeUnusedHandlers()
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 1)

	e.Mount(Div())
	e.Consume()
	m.closeUnusedHandlers()
	require.Empty(t, m.handlers)
}

func TestActionBuilder(t *testing.T) {
	actionName := "/text/actionBuilder"
	createdAction := Action{}

	wg := sync.WaitGroup{}
	wg.Add(1)

	Handle(actionName, func(ctx Context, a Action) {
		createdAction = a
		wg.Done()
	})

	disp := NewClientTester(Div())
	defer disp.Close()

	newActionBuilder(disp, actionName).
		Value(42).
		Tag("foo", "bar").
		Post()

	wg.Wait()
	require.Equal(t, actionName, createdAction.Name)
	require.Equal(t, 42, createdAction.Value)
	require.Equal(t, "bar", createdAction.Tags.Get("foo"))
}
