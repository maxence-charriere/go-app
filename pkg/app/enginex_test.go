package app

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEngineXLoad(t *testing.T) {
	t.Run("load loads a new body", func(t *testing.T) {
		e := newTestEngine()
		e.load(&hello{})
		require.IsType(t, &hello{}, e.body.(HTML).body()[0])
	})

	t.Run("loading a non mountable component panics", func(t *testing.T) {
		e := newTestEngine()
		require.Panics(t, func() {
			e.load(&compoWithNilRendering{})
		})
	})

	t.Run("load updates body", func(t *testing.T) {
		e := newTestEngine()
		e.load(&hello{})
		e.load(&bar{})
		require.IsType(t, &bar{}, e.body.(HTML).body()[0])
	})

	t.Run("load body update with a non mountable component panics", func(t *testing.T) {
		e := newTestEngine()
		e.load(&hello{})
		require.Panics(t, func() {
			e.load(&compoWithNilRendering{})
		})
	})
}

func TestEngineXNavigate(t *testing.T) {
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
		require.IsType(t, &notFound{}, e.body.(HTML).body()[0])
	})
}

func newTestEngine() *engineX {
	url, _ := url.Parse("/")
	routes := makeRouter()

	return newEngineX(context.Background(),
		&routes,
		nil,
		url,
		Body,
	)
}
