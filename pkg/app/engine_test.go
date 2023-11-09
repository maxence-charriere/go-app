package app

import (
	"bytes"
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEngineBaseContext(t *testing.T) {
	e := newTestEngine()
	ctx := e.baseContext()
	require.NotNil(t, ctx.Context)
	require.NotNil(t, ctx.page)
	require.NotNil(t, ctx.resolveURL)
	require.NotNil(t, ctx.navigate)
	require.NotNil(t, ctx.localStorage)
	require.NotNil(t, ctx.sessionStorage)
	require.NotNil(t, ctx.dispatch)
	require.NotNil(t, ctx.defere)
	require.NotNil(t, ctx.async)
	require.NotNil(t, ctx.addComponentUpdate)
	require.NotNil(t, ctx.removeComponentUpdate)
	require.NotNil(t, ctx.handleAction)
	require.NotNil(t, ctx.postAction)
	require.NotNil(t, ctx.observeState)
	require.NotNil(t, ctx.getState)
	require.NotNil(t, ctx.setState)
	require.NotNil(t, ctx.delState)

	require.NotNil(t, ctx.notifyComponentEvent)
}

func TestEngineLoad(t *testing.T) {
	t.Run("load loads a new body", func(t *testing.T) {
		e := newTestEngine()
		e.Load(&hello{})
		require.IsType(t, &hello{}, e.body.body()[0])
	})

	t.Run("loading a non mountable component panics", func(t *testing.T) {
		e := newTestEngine()
		err := e.Load(&compoWithNilRendering{})
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("load updates body", func(t *testing.T) {
		e := newTestEngine()
		e.Load(&hello{})
		e.Load(&bar{})
		require.IsType(t, &bar{}, e.body.body()[0])
	})

	t.Run("load body update with a non mountable component panics", func(t *testing.T) {
		e := newTestEngine()
		e.Load(&hello{})
		err := e.Load(&compoWithNilRendering{})
		require.Error(t, err)
		t.Log(err)
	})
}

func TestEngineNavigate(t *testing.T) {
	t.Run("url is loaded and history is updated", func(t *testing.T) {
		e := newTestEngine()
		e.routes.route("/hello", NewZeroComponentFactory(&hello{}))

		destination, _ := url.Parse("/hello")
		e.Navigate(destination, true)
		require.Equal(t, "/hello", e.lastVisitedURL.Path)
	})

	t.Run("url is loaded and history is not updated", func(t *testing.T) {
		e := newTestEngine()
		e.routes.route("/hello", NewZeroComponentFactory(&hello{}))

		destination, _ := url.Parse("/hello")
		e.Navigate(destination, false)
		require.Equal(t, "/hello", e.lastVisitedURL.Path)
	})

	t.Run("mailto is loaded", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("mailto:contact@murlok.io")
		e.Navigate(destination, true)
	})

	t.Run("external url is opened", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("https://murlok.io")
		e.Navigate(destination, true)
	})

	t.Run("navigation on current page is skipped", func(t *testing.T) {
		e := newTestEngine()
		e.routes.route("/hello", NewZeroComponentFactory(&hello{}))

		destination, _ := url.Parse("/hello#bye")
		e.Navigate(destination, true)
		lastVisitedURL := e.lastVisitedURL

		e.Navigate(destination, true)
		require.Equal(t, lastVisitedURL, e.lastVisitedURL)
	})

	t.Run("url with fragment is loaded", func(t *testing.T) {
		e := newTestEngine()
		e.routes.route("/hello", NewZeroComponentFactory(&hello{}))

		destination, _ := url.Parse("/hello#bye")
		e.Navigate(destination, true)
		require.Equal(t, "/hello", e.lastVisitedURL.Path)
		require.Equal(t, "bye", e.lastVisitedURL.Fragment)
	})

	t.Run("fragment navigation after initial load", func(t *testing.T) {
		e := newTestEngine()
		e.routes.route("/hello", NewZeroComponentFactory(&hello{}))

		destination, _ := url.Parse("/hello")
		e.Navigate(destination, true)
		require.Equal(t, "/hello", e.lastVisitedURL.Path)
		require.Empty(t, e.lastVisitedURL.Fragment)

		destination, _ = url.Parse("/hello#bye")
		e.Navigate(destination, true)
		require.Equal(t, "bye", e.lastVisitedURL.Fragment)
	})

	t.Run("url with prefix root is loaded", func(t *testing.T) {
		e := newTestEngine()
		e.routes.route("/", NewZeroComponentFactory(&hello{}))

		os.Setenv("GOAPP_ROOT_PREFIX", "/prefix")
		destination, _ := url.Parse("/prefix")
		e.Navigate(destination, true)
		require.Equal(t, "/prefix", e.lastVisitedURL.Path)
	})

	t.Run("not found component is loaded", func(t *testing.T) {
		e := newTestEngine()

		destination, _ := url.Parse("/hello")
		e.Navigate(destination, true)
		require.IsType(t, &notFound{}, e.body.body()[0])
	})
}

func TestEngineInternalURL(t *testing.T) {
	t.Run("destination is internal URL", func(t *testing.T) {
		os.Setenv("GOAPP_INTERNAL_URLS", `["https://murlok.io"]`)
		defer os.Unsetenv("GOAPP_INTERNAL_URLS")

		e := newTestEngine()
		destination, _ := url.Parse("https://murlok.io/warrior")
		require.True(t, e.internalURL(destination))
	})

	t.Run("destination is internal URL", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("https://murlok.io/warrior")
		require.False(t, e.internalURL(destination))
	})
}

func TestEngineMailTo(t *testing.T) {
	t.Run("destination is mailto", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("mailto:maxence@goapp.dev")
		require.True(t, e.mailTo(destination))
	})

	t.Run("destination is not mailto", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("/hello")
		require.False(t, e.mailTo(destination))
	})
}

func TestEngineExternalNavigation(t *testing.T) {
	t.Run("destination is external navigation", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("https://murlok.io")
		require.True(t, e.externalNavigation(destination))
	})

	t.Run("destination is not external navigation", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("/hello")
		require.False(t, e.externalNavigation(destination))
	})
}

func TestEngineAsync(t *testing.T) {
	e := newTestEngine()

	called := false
	e.async(func() {
		called = true
	})

	e.goroutines.Wait()
	require.True(t, called)
}

func TestEngineStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	routes := makeRouter()
	routes.route("/", func() Composer {
		return &navigatorComponent{
			onNav: func(ctx Context) {
				ctx.Dispatch(func(ctx Context) {
					ctx.Defer(func(ctx Context) {
						cancel()
					})
				})
			},
		}
	})

	e := newTestEngine()
	e.ctx = ctx
	e.routes = &routes

	destination, _ := url.Parse("/")
	e.Navigate(destination, false)
	e.Start(0)
}

func TestEngineEncode(t *testing.T) {
	t.Run("encoding when engine did not load a component returns an error", func(t *testing.T) {
		e := newTestEngine()

		var b bytes.Buffer
		err := e.Encode(&b, Html())
		require.Error(t, err)
		require.Empty(t, b.Bytes())
	})

	t.Run("encoding a document without body returns an error", func(t *testing.T) {
		e := newTestEngine()
		compo := &hello{}
		e.Load(compo)

		var b bytes.Buffer
		err := e.Encode(&b, Html())
		require.Error(t, err)
		require.Empty(t, b.Bytes())
	})

	t.Run("encoding document succed", func(t *testing.T) {
		e := newTestEngine()
		compo := &compoWithCustomRoot{Root: Span()}
		e.Load(compo)

		var b bytes.Buffer
		err := e.Encode(&b, Html().privateBody(
			Body().privateBody(
				Text("bye"),
			),
		))
		require.NoError(t, err)
		require.Equal(t, "<!DOCTYPE html>\n<html>\n  <body>\n    <span></span>\n    bye\n  </body>\n</html>", b.String())
	})
}

func newTestEngine() *engineX {
	return NewTestEngine().(*engineX)
}
