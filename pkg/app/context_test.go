package app

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestContextBehavior(t *testing.T) {
	ctx1 := context.Background()

	ctx2, cancel2 := context.WithCancel(ctx1)
	defer cancel2()

	ctx3, cancel3 := context.WithCancel(ctx2)
	defer cancel3()

	ctx4, cancel4 := context.WithCancel(ctx2)
	defer cancel4()

	ctx5, cancel5 := context.WithCancel(ctx4)
	defer cancel5()

	cancel4()

	require.NoError(t, ctx1.Err())
	require.NoError(t, ctx2.Err())
	require.NoError(t, ctx3.Err())
	require.Error(t, ctx4.Err())
	require.Error(t, ctx5.Err())
}

func TestContextSrc(t *testing.T) {
	e := newTestEngine()
	hello := &hello{}
	e.Load(hello)

	ctx := e.nodes.context(e.baseContext(), hello)
	require.NotZero(t, ctx.Src())
	require.NotNil(t, ctx.JSSrc())
}

func TestContextDeviceID(t *testing.T) {
	ctx := makeTestContext()
	id := ctx.DeviceID()
	require.NotZero(t, id)

	id2 := ctx.DeviceID()
	require.Equal(t, id, id2)
}

func TestContextAppUpdateAvailable(t *testing.T) {
	e := newTestEngine()
	ctx := e.baseContext()
	require.False(t, ctx.AppUpdateAvailable())
}

func TestContextAppInstallable(t *testing.T) {
	e := newTestEngine()
	ctx := e.baseContext()
	require.False(t, ctx.IsAppInstallable())
	ctx.ShowAppInstallPrompt()
}

func TestContextReload(t *testing.T) {
	if IsClient {
		t.Skip()
	}
	newTestEngine().
		baseContext().
		Reload()
}

func TestContextNavigate(t *testing.T) {
	e := newTestEngine()
	ctx := e.baseContext()

	t.Run("navigate succeeds", func(t *testing.T) {
		ctx.Navigate("https://murlok.io")
	})

	t.Run("navigate to a bad url logs an error", func(t *testing.T) {
		ctx.Navigate("ad;lsfjk:/:;/murlok.io")
	})
}

func TestContextResolveStaticResource(t *testing.T) {
	e := newTestEngine()
	ctx := e.baseContext()
	require.Equal(t, "/test", ctx.ResolveStaticResource("/test"))
}

func TestContextScrollTo(t *testing.T) {
	e := newTestEngine()
	ctx := e.baseContext()
	ctx.ScrollTo("test")
}

func TestContextStorage(t *testing.T) {
	e := newTestEngine()
	ctx := e.baseContext()

	t.Run("local storage is set", func(t *testing.T) {
		require.NotZero(t, ctx.LocalStorage())
	})

	t.Run("session storage is set", func(t *testing.T) {
		require.NotZero(t, ctx.SessionStorage())
	})
}

func TestContextEncryptDecryptStruct(t *testing.T) {
	ctx := makeTestContext()

	expected := struct {
		Title string
		Value int
	}{
		Title: "hello",
		Value: 42,
	}

	item := expected
	item.Title = ""
	item.Value = 0

	crypted, err := ctx.Encrypt(expected)
	require.NoError(t, err)
	require.NotEmpty(t, crypted)

	err = ctx.Decrypt(crypted, &item)
	require.NoError(t, err)
	require.Equal(t, expected, item)
}

func TestContextEncryptDecryptString(t *testing.T) {
	ctx := makeTestContext()

	expected := "hello"
	item := ""

	crypted, err := ctx.Encrypt(expected)
	require.NoError(t, err)
	require.NotEmpty(t, crypted)

	err = ctx.Decrypt(crypted, &item)
	require.NoError(t, err)
	require.Equal(t, expected, item)
}

func TestContextEncryptDecryptBytes(t *testing.T) {
	ctx := makeTestContext()

	expected := []byte("hello")
	var item []byte

	crypted, err := ctx.Encrypt(expected)
	require.NoError(t, err)
	require.NotEmpty(t, crypted)

	err = ctx.Decrypt(crypted, &item)
	require.NoError(t, err)
	require.Equal(t, expected, item)
}

func TestContextNotificationService(t *testing.T) {
	ctx := makeTestContext()
	ctx.Notifications()
}

func TestContextDispatch(t *testing.T) {
	t.Run("function is executed when source element is mounted", func(t *testing.T) {
		e := newTestEngine()

		hello := &hello{}
		e.Load(hello)
		called := false

		ctx := e.nodes.context(e.baseContext(), hello)
		ctx.Dispatch(func(ctx Context) {
			called = true
		})

		e.ConsumeAll()
		require.True(t, called)
	})

	t.Run("function is skipped when source element is not mounted", func(t *testing.T) {
		e := newTestEngine()

		hello := &hello{}
		called := false

		ctx := e.nodes.context(e.baseContext(), hello)
		ctx.Dispatch(func(ctx Context) {
			called = true
		})

		e.ConsumeAll()
		require.False(t, called)
	})
}

func TestContextDefer(t *testing.T) {
	t.Run("function is executed when source element is mounted", func(t *testing.T) {
		e := newTestEngine()

		hello := &hello{}
		e.Load(hello)
		called := false

		ctx := e.nodes.context(e.baseContext(), hello)
		ctx.Defer(func(ctx Context) {
			called = true
		})

		e.ConsumeAll()
		require.True(t, called)
	})

	t.Run("function is skipped dispatched when source element is not mounted", func(t *testing.T) {
		e := newTestEngine()

		hello := &hello{}
		called := false

		ctx := e.nodes.context(e.baseContext(), hello)
		ctx.Defer(func(ctx Context) {
			called = true
		})

		e.ConsumeAll()
		require.False(t, called)
	})
}

func TestContextAfter(t *testing.T) {
	waitForWorkers := func(t *testing.T, e *engineX) {
		t.Helper()
		done := make(chan struct{})
		go func() {
			e.waitAsync(context.Background())
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("delayed worker did not finish")
		}
	}

	for _, duration := range []time.Duration{time.Millisecond, 0, -time.Millisecond} {
		t.Run(duration.String(), func(t *testing.T) {
			e := newTestEngine()
			hello := &hello{}
			require.NoError(t, e.Load(hello))
			ctx := e.nodes.context(e.baseContext(), hello)
			calls := 0
			ctx.After(duration, func(Context) { calls++ })

			waitForWorkers(t, e)
			e.ConsumeAll()
			require.Equal(t, 1, calls)
		})
	}

	for _, test := range []struct {
		name         string
		duration     time.Duration
		cancelBefore bool
	}{
		{"canceled before scheduling", time.Hour, true},
		{"canceled with ready timer", 0, true},
		{"canceled after scheduling", time.Hour, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			e := newTestEngine()
			var cancel context.CancelFunc
			e.ctx, cancel = context.WithCancel(e.ctx)
			defer cancel()
			hello := &hello{}
			require.NoError(t, e.Load(hello))
			ctx := e.nodes.context(e.baseContext(), hello)
			called := false
			if test.cancelBefore {
				cancel()
			}
			ctx.After(test.duration, func(Context) { called = true })
			cancel()

			waitForWorkers(t, e)
			e.ConsumeAll()
			require.False(t, called)
		})
	}

	t.Run("canceled after callback is queued", func(t *testing.T) {
		e := newTestEngine()
		var cancel context.CancelFunc
		e.ctx, cancel = context.WithCancel(e.ctx)
		defer cancel()
		hello := &hello{}
		require.NoError(t, e.Load(hello))
		ctx := e.nodes.context(e.baseContext(), hello)
		called := false
		ctx.After(0, func(Context) { called = true })

		waitForWorkers(t, e)
		require.NotEmpty(t, e.dispatches)
		cancel()
		e.ConsumeAll()
		require.False(t, called)
	})

	t.Run("nil callback", func(t *testing.T) {
		e := newTestEngine()
		hello := &hello{}
		require.NoError(t, e.Load(hello))
		ctx := e.nodes.context(e.baseContext(), hello)
		ctx.After(0, nil)

		waitForWorkers(t, e)
		require.NotPanics(t, e.ConsumeAll)
	})
}

func TestContextPreventUpdate(t *testing.T) {
	e := newTestEngine()

	hello := &hello{}
	e.Load(hello)
	e.ConsumeAll()

	ctx := e.nodes.context(e.baseContext(), hello)
	ctx.Dispatch(func(ctx Context) {
		require.Contains(t, e.updates.pending[1], hello)
		require.Equal(t, 1, e.updates.pending[1][hello])
		ctx.PreventUpdate()
		require.Zero(t, e.updates.pending[1][hello])
	})

	e.ConsumeAll()
}

func TestContextUpdate(t *testing.T) {
	e := newTestEngine()

	hello := &hello{}
	e.Load(hello)
	e.ConsumeAll()

	ctx := e.nodes.context(e.baseContext(), hello)
	ctx.Update()
	ctx.Dispatch(func(ctx Context) {
		require.Contains(t, e.updates.pending[1], hello)
		require.Equal(t, 2, e.updates.pending[1][hello])
	})

	e.ConsumeAll()
}

func TestContextHandle(t *testing.T) {
	e := newTestEngine()

	actionName := "/test/context/handle"
	action := Action{}

	hello := &hello{}
	e.Load(hello)
	ctx := e.nodes.context(e.baseContext(), hello)

	ctx.Handle(actionName, func(ctx Context, a Action) {
		action = a
	})

	ctx.NewActionWithValue(actionName, 21, T("hello", "world"), Tags{"foo": "bar"})
	e.ConsumeAll()
	require.Equal(t, actionName, action.Name)
	require.Equal(t, 21, action.Value)
	require.Equal(t, "world", action.Tags.Get("hello"))
	require.Equal(t, "bar", action.Tags.Get("foo"))

	ctx.NewAction(actionName)
	e.ConsumeAll()
	require.Equal(t, actionName, action.Name)
	require.Nil(t, action.Value)
	require.Nil(t, action.Tags)
}

func TestContextStates(t *testing.T) {
	e := newTestEngine()

	hello := &hello{}
	e.Load(hello)
	ctx := e.nodes.context(e.baseContext(), hello)

	state := "/test/context/states"
	var v string

	ctx.SetState(state, "hello")
	ctx.GetState(state, &v)
	require.Equal(t, "hello", v)

	ctx.ObserveState(state, &v)
	ctx.SetState(state, "bye")

	e.ConsumeAll()
	require.Equal(t, "bye", v)

	ctx.DelState(state)
	require.Empty(t, e.states.states)
}

func TestContextResizeContent(t *testing.T) {
	e := newTestEngine()
	hello := &hello{}
	e.Load(hello)
	ctx := e.nodes.context(e.baseContext(), hello)
	ctx.ResizeContent()
	e.ConsumeAll()
}

func makeTestContext() Context {
	resolveURL := func(v string) string {
		return v
	}

	var page Page
	url, _ := url.Parse("https://goapp.dev")
	if IsServer {
		requestPage := makeRequestPage(url, resolveURL)
		page = &requestPage
	} else {
		page = makeBrowserPage(resolveURL)
	}

	var localStorage BrowserStorage
	var sessionStorage BrowserStorage
	if IsServer {
		localStorage = newMemoryStorage()
		sessionStorage = newMemoryStorage()
	} else {
		localStorage = newJSStorage("localStorage")
		sessionStorage = newJSStorage("sessionStorage")
	}

	return Context{
		Context:               context.Background(),
		page:                  func() Page { return page },
		resolveURL:            resolveURL,
		localStorage:          localStorage,
		sessionStorage:        sessionStorage,
		dispatch:              func(f func()) { f() },
		defere:                func(f func()) { f() },
		async:                 func(f func()) { f() },
		addComponentUpdate:    func(Composer, int) {},
		removeComponentUpdate: func(Composer) {},
		handleAction:          func(string, UI, bool, ActionHandler) {},
		postAction:            func(Context, Action) {},
	}
}
