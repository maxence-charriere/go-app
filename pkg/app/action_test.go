package app

import (
	"sync"
	"testing"
	"time"

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

func TestActionManagerPostContexts(t *testing.T) {
	for _, test := range []struct {
		name  string
		async []bool
	}{
		{
			name:  "async handlers",
			async: []bool{true, true, true},
		},
		{
			name:  "mixed async and UI handlers",
			async: []bool{true, false, false},
		},
	} {
		for _, scheduling := range []string{
			"after posting",
			"during posting",
		} {
			t.Run(test.name+"/"+scheduling, func(t *testing.T) {
				var nm nodeManager
				var am actionManager
				ctx := makeTestContext()
				const posts = 20
				calls := len(test.async) * posts
				dispatches := make(chan func(), calls*2)
				ctx.dispatch = func(f func()) { dispatches <- f }

				// Delay execution before invoking the closure, so its captured
				// context is read only after Post has visited every handler.
				ready := make(chan struct{})
				if scheduling == "during posting" {
					close(ready)
				}
				var workers sync.WaitGroup
				ctx.async = func(f func()) {
					workers.Add(1)
					go func() {
						defer workers.Done()
						<-ready
						f()
					}()
				}

				type result struct {
					want UI
					got  UI
				}
				results := make(chan result, calls*2)
				for _, async := range test.async {
					source, err := nm.Mount(ctx, 1, Div())
					require.NoError(t, err)
					defer nm.Dismount(source)
					am.Handle("test", source, async, func(ctx Context, _ Action) {
						results <- result{want: source, got: ctx.Src()}
						ctx.Dispatch(func(ctx Context) {
							results <- result{want: source, got: ctx.Src()}
						})
					})
				}

				for range posts {
					am.Post(ctx, Action{Name: "test"})
				}
				if scheduling == "after posting" {
					close(ready)
				}
				done := make(chan struct{})
				go func() {
					workers.Wait()
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(5 * time.Second):
					t.Fatal("async action handlers did not finish")
				}

				// Run UI handlers and follow-up dispatches on the test goroutine.
				for len(dispatches) != 0 {
					(<-dispatches)()
				}
				require.Len(t, results, calls*2)
				for len(results) != 0 {
					result := <-results
					require.Same(t, result.want, result.got)
				}
			})
		}
	}
}
