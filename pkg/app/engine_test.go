package app

import (
	"bytes"
	"context"
	"net/url"
	"os"
	"runtime"
	"testing"
	"time"

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

	t.Run("tel is loaded", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("tel:+18005551234")
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

func TestEngineTel(t *testing.T) {
	t.Run("destination is tel", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("tel:+18005551234")
		require.True(t, e.tel(destination))
	})

	t.Run("destination is not tel", func(t *testing.T) {
		e := newTestEngine()
		destination, _ := url.Parse("/hello")
		require.False(t, e.tel(destination))
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

func TestEngineQueueCancellation(t *testing.T) {
	for _, queueTest := range []struct {
		name    string
		queue   func(*engineX) chan func()
		enqueue func(*engineX, func())
	}{
		{"dispatch", func(e *engineX) chan func() { return e.dispatches }, (*engineX).dispatch},
		{"defer", func(e *engineX) chan func() { return e.defers }, (*engineX).defere},
	} {
		t.Run(queueTest.name, func(t *testing.T) {
			for _, test := range []struct {
				name         string
				full         bool
				cancelBefore bool
				cancelDuring bool
			}{
				{name: "enqueue normally"},
				{name: "wait for queue space", full: true},
				{name: "already canceled", cancelBefore: true},
				{name: "already canceled with full queue", full: true, cancelBefore: true},
				{name: "cancel sender with full queue", full: true, cancelDuring: true},
			} {
				t.Run(test.name, func(t *testing.T) {
					ctx, cancel := context.WithCancel(context.Background())
					e := &engineX{
						ctx:        ctx,
						dispatches: make(chan func(), 1),
						defers:     make(chan func(), 1),
					}
					queue := queueTest.queue(e)
					t.Cleanup(func() {
						cancel()
						// Release the sender even if a broken implementation ignores cancellation.
						select {
						case <-queue:
						default:
						}
					})
					var calls []int
					if test.full {
						queue <- func() { calls = append(calls, 1) }
					}
					if test.cancelBefore {
						cancel()
					}

					started := make(chan struct{})
					e.async(func() {
						close(started)
						queueTest.enqueue(e, func() { calls = append(calls, 2) })
					})
					<-started
					if test.cancelDuring {
						cancel()
					} else if test.full && !test.cancelBefore {
						(<-queue)()
					}

					finished := make(chan struct{})
					go func() {
						e.waitAsync(context.Background())
						close(finished)
					}()
					select {
					case <-finished:
					case <-time.After(5 * time.Second):
						t.Fatal("queue sender did not finish")
					}

					if test.cancelBefore || test.cancelDuring {
						if test.full {
							require.Len(t, queue, 1)
							(<-queue)()
							require.Equal(t, []int{1}, calls)
						}
						require.Empty(t, queue)
						return
					}
					require.Len(t, queue, 1)
					(<-queue)()
					if test.full {
						require.Equal(t, []int{1, 2}, calls)
					} else {
						require.Equal(t, []int{2}, calls)
					}
				})
			}
		})
	}
}

func TestEngineAsync(t *testing.T) {
	e := newTestEngine()

	called := false
	e.async(func() {
		called = true
	})

	e.waitAsync(context.Background())
	require.True(t, called)
}

func TestEngineAsyncCompletion(t *testing.T) {
	e := newTestEngine()
	for range 3 {
		called := false
		e.async(func() {
			e.async(func() { called = true })
			runtime.Goexit()
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		completed := e.waitAsync(ctx)
		cancel()
		require.True(t, completed, "nested work must complete even when its parent exits with Goexit")
		require.True(t, called)
	}
}

func TestEngineConsumeCancellation(t *testing.T) {
	for _, test := range []struct {
		name    string
		consume func(*engineX)
	}{
		{"next", (*engineX).ConsumeNext},
		{"all", (*engineX).ConsumeAll},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("pending async work", func(t *testing.T) {
				e := newTestEngine()
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
				e.ctx = ctx
				release := make(chan struct{})
				e.async(func() { <-release })
				finished := make(chan struct{})
				t.Cleanup(func() {
					cancel()
					close(release)
					cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
					defer stop()
					require.True(t, e.waitAsync(cleanup))
				})
				go func() {
					test.consume(e)
					close(finished)
				}()
				select {
				case <-finished:
				case <-time.After(5 * time.Second):
					t.Fatal("consumption did not stop while async work was still pending")
				}
			})

			for _, when := range []string{"before consumption", "in dispatch", "in defer"} {
				t.Run(when, func(t *testing.T) {
					e := newTestEngine()
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					e.ctx = ctx
					var calls []string
					e.dispatch(func() {
						calls = append(calls, "dispatch")
						if when == "in dispatch" {
							cancel()
						}
					})
					if when != "in defer" {
						e.dispatch(func() { calls = append(calls, "unexpected dispatch") })
					}
					e.defere(func() {
						calls = append(calls, "defer")
						cancel()
					})
					e.defere(func() { calls = append(calls, "unexpected defer") })
					if when == "before consumption" {
						cancel()
					}
					test.consume(e)
					switch when {
					case "before consumption":
						require.Empty(t, calls)
					case "in dispatch":
						require.Equal(t, []string{"dispatch"}, calls)
					case "in defer":
						require.Equal(t, []string{"dispatch", "defer"}, calls)
					}
				})
			}
		})
	}
}

func TestEngineConsumeNextEmptyQueueCancellation(t *testing.T) {
	e := newTestEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	e.ctx = ctx
	finished := make(chan struct{})
	go func() {
		e.ConsumeNext()
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("consumption did not stop with an empty dispatch queue")
	}
}

func TestEngineConsumeAsyncDispatches(t *testing.T) {
	e := newTestEngine()
	var calls []int
	e.async(func() {
		e.async(func() {
			e.dispatch(func() {
				calls = append(calls, 1)
				e.async(func() {
					e.dispatch(func() { calls = append(calls, 2) })
				})
			})
		})
	})
	e.ConsumeNext()
	require.Equal(t, []int{1}, calls)
	e.ConsumeAll()
	require.Equal(t, []int{1, 2}, calls)
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
		require.Equal(t, "<!doctype html>\n<html>\n  <body>\n    <span></span>\n    bye\n  </body>\n</html>", b.String())
	})
}

func newTestEngine() *engineX {
	return NewTestEngine().(*engineX)
}
