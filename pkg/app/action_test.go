package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	Handle("/test", func(Context, Action) {})
	require.Len(t, actionHandlers, 1)
}

func TestActionManagerHandle(t *testing.T) {
	var m actionManager

	source := Div()
	handlerCalled := false
	handler := func(ctx Context, a Action) {
		handlerCalled = true
	}
	m.Handle("test", source, true, handler)
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["test"], 1)

	actionHandler := m.handlers["test"][actionHandlerKey(source, handler)]
	require.NotZero(t, actionHandler)
	require.Equal(t, source, actionHandler.Source)
	require.NotNil(t, actionHandler.Function)

	actionHandler.Function(Context{}, Action{})
	require.True(t, handlerCalled)
}

func TestActionManagerPost(t *testing.T) {
	t.Run("action handler is called asynchronously", func(t *testing.T) {
		var nm nodeManager
		var am actionManager

		ctx := makeTestContext()
		source, err := nm.Mount(ctx, 1, Div())
		ctx = nm.context(ctx, source)
		require.NoError(t, err)

		handlerCalled := false
		am.Handle("test", source, true, func(ctx Context, a Action) {
			handlerCalled = true
		})

		am.Post(ctx, Action{
			Name: "test",
		})
		require.True(t, handlerCalled)
	})

	t.Run("action handler is called synchronously", func(t *testing.T) {
		var nm nodeManager
		var am actionManager

		ctx := makeTestContext()
		source, err := nm.Mount(ctx, 1, Div())
		ctx = nm.context(ctx, source)
		require.NoError(t, err)

		handlerCalled := false
		am.Handle("test", source, false, func(ctx Context, a Action) {
			handlerCalled = true
		})

		am.Post(ctx, Action{
			Name: "test",
		})
		require.True(t, handlerCalled)
	})

	t.Run("action handler is removed when source is dismounted", func(t *testing.T) {
		var m actionManager

		source := Div()
		handlerCalled := false
		m.Handle("test", source, true, func(ctx Context, a Action) {
			handlerCalled = true
		})

		m.Post(makeTestContext(), Action{
			Name: "test",
		})
		require.False(t, handlerCalled)
		require.Len(t, m.handlers, 1)
		require.Empty(t, m.handlers["test"])
	})
}

func TestActionManagerCleanup(t *testing.T) {
	var m actionManager

	m.Handle("test", Div(), true, func(ctx Context, a Action) {})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["test"], 1)

	m.Cleanup()
	require.Empty(t, m.handlers)
}
