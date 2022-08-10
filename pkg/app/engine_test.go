package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngineInit(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	assert.NotZero(t, e.FrameRate)
	assert.NotNil(t, e.Page)
	assert.NotNil(t, e.LocalStorage)
	assert.NotNil(t, e.SessionStorage)
	assert.NotNil(t, e.StaticResourceResolver)
	assert.NotNil(t, e.Body)
	assert.NotNil(t, e.dispatches)
	assert.NotNil(t, e.componentUpdates)
	assert.NotNil(t, e.deferables)
}

func TestEngineDispatch(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	e.Dispatch(Dispatch{})

	require.Len(t, e.dispatches, 1)

	d := <-e.dispatches
	require.Equal(t, Update, d.Mode)
	require.Equal(t, e.Body, d.Source)
}

func TestEngineEmit(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	foo := &foo{Bar: "bar"}
	e.Mount(foo)
	e.Consume()
	require.Empty(t, e.dispatches)
	require.Empty(t, e.componentUpdates)

	bar := foo.getChildren()[0].(*bar)

	emitted := false
	e.Emit(bar, func() {
		emitted = true
	})
	require.False(t, emitted)
	require.Len(t, e.dispatches, 1)

	e.Consume()
	require.True(t, emitted)
	require.Empty(t, e.dispatches)
}

func TestEngineHandleDispatch(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		e := engine{}
		e.init()
		defer e.Close()

		bar := &bar{}
		e.Mount(bar)
		e.Consume()

		called := false
		e.handleDispatch(Dispatch{
			Mode:     Update,
			Source:   bar,
			Function: func(Context) { called = true },
		})
		require.True(t, called)
		require.NotEmpty(t, e.componentUpdates)
	})

	t.Run("defer", func(t *testing.T) {
		e := engine{}
		e.init()
		defer e.Close()

		bar := &bar{}
		e.Mount(bar)
		e.Consume()

		called := false
		e.handleDispatch(Dispatch{
			Mode:     Defer,
			Source:   bar,
			Function: func(Context) { called = true },
		})
		require.Empty(t, e.componentUpdates)
		require.Len(t, e.deferables, 1)
		require.False(t, called)
	})

	t.Run("next", func(t *testing.T) {
		e := engine{}
		e.init()
		defer e.Close()

		bar := &bar{}
		e.Mount(bar)
		e.Consume()

		called := false
		e.handleDispatch(Dispatch{
			Mode:     Next,
			Source:   bar,
			Function: func(Context) { called = true },
		})
		require.True(t, called)
		require.Empty(t, e.componentUpdates)
	})
}

func TestEngineAddComponentUpdate(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	h := &hello{}
	e.addComponentUpdate(h)
	require.Empty(t, e.componentUpdates)

	e.Mount(h)
	e.Consume()
	require.Empty(t, e.dispatches)
	require.Empty(t, e.componentUpdates)

	e.addComponentUpdate(h)
	require.Len(t, e.componentUpdates, 1)
	require.True(t, e.componentUpdates[h])

	e.addComponentUpdate(h)
	require.Len(t, e.componentUpdates, 1)
}

func TestPreventComponentUpdate(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	h := &hello{}
	e.Mount(h)
	e.Consume()
	require.Empty(t, e.dispatches)
	require.Empty(t, e.componentUpdates)

	e.preventComponentUpdate(h)
	require.Len(t, e.componentUpdates, 1)
	require.False(t, e.componentUpdates[h])
}

func TestEngineHandleComponentUpdates(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	foo := &foo{Bar: "bar"}
	e.Mount(foo)
	e.Consume()
	require.Empty(t, e.dispatches)
	require.Empty(t, e.componentUpdates)
	bar := foo.root.(*bar)

	e.addComponentUpdate(foo)
	e.addComponentUpdate(bar)
	require.Len(t, e.componentUpdates, 2)

	e.handleComponentUpdates()
	require.Empty(t, e.componentUpdates)
}

func TestEngineExecDeferableEvents(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	h := &hello{}
	e.Mount(h)
	e.Consume()
	require.Empty(t, e.dispatches)
	require.Empty(t, e.componentUpdates)

	called := false

	e.addDeferable(Dispatch{
		Mode:   Defer,
		Source: h,
		Function: func(Context) {
			called = true
		},
	})
	require.Len(t, e.deferables, 1)

	e.handleDeferables()
	require.True(t, called)
	require.Empty(t, e.deferables)
}

func TestEngineHandlePost(t *testing.T) {
	isAppHandleCalled := false
	isHandleACalled := false
	isHandleBCalled := false
	isHandleCCalled := false

	e := engine{
		ActionHandlers: map[string]ActionHandler{
			"/test": func(ctx Context, a Action) {
				isAppHandleCalled = true
			},
		},
	}
	e.init()
	defer e.Close()

	h := &hello{}
	e.Mount(h)
	e.Consume()

	e.Handle("/test", h, func(ctx Context, a Action) {
		isHandleACalled = true
	})

	e.Handle("/test", h, func(ctx Context, a Action) {
		isHandleBCalled = true
	})

	f := &foo{}
	e.Handle("/test", f, func(ctx Context, a Action) {
		isHandleCCalled = true
	})

	e.Post(Action{Name: "/test"})
	e.Consume()

	require.True(t, isAppHandleCalled)
	require.True(t, isHandleACalled)
	require.True(t, isHandleBCalled)
	require.False(t, isHandleCCalled)
}
