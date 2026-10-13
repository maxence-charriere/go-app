//go:build !wasm
// +build !wasm

package app

import (
	"context"
	_ "embed"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	Route("/", func() Composer { return &preRenderTestCompo{} })
}

type preRenderTestCompo struct {
	Compo
}

func (c *preRenderTestCompo) Render() UI {
	return Div().
		ID("pre-render-ok").
		Body(
			Img().Src("/web/resolve-static-resource-test.jpg"),
		)
}

type requestContextTestCompo struct {
	Compo
	onPreRender func(Context)
}

func (c *requestContextTestCompo) OnPreRender(ctx Context) {
	c.onPreRender(ctx)
}

func (c *requestContextTestCompo) Render() UI {
	return Div().Text("request context")
}

func TestHandlerRequestContext(t *testing.T) {
	type requestKey struct{}
	for _, test := range []struct {
		name           string
		cancelBefore   bool
		cancelDuring   bool
		cancelEncoding bool
		expired        bool
	}{
		{name: "values, deadline, and render cleanup"},
		{name: "canceled before rendering", cancelBefore: true},
		{name: "canceled during rendering", cancelDuring: true},
		{name: "canceled during encoding", cancelEncoding: true},
		{name: "deadline exceeded before rendering", expired: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			deadline := time.Now().Add(time.Hour)
			if test.expired {
				deadline = time.Now().Add(-time.Second)
			}
			parent, cancel := context.WithDeadline(
				context.WithValue(context.Background(), requestKey{}, "request value"), deadline)
			defer cancel()
			if test.cancelBefore {
				cancel()
			}

			var renderContext context.Context
			var renderError error
			const path = "/request-context-test"
			Route(path, func() Composer {
				return &requestContextTestCompo{onPreRender: func(ctx Context) {
					renderContext = ctx.Context
					if test.cancelDuring {
						cancel()
					}
					renderError = ctx.Err()
				}}
			})
			t.Cleanup(func() {
				routes.mu.Lock()
				defer routes.mu.Unlock()
				delete(routes.routes, path)
			})

			r := httptest.NewRequest(http.MethodGet, path, nil).WithContext(parent)
			w := httptest.NewRecorder()
			h := Handler{}
			if test.cancelEncoding {
				h.HTML = func() HTMLHtml {
					cancel()
					return Html()
				}
			}
			h.ServeHTTP(w, r)

			if test.cancelBefore || test.expired {
				require.Nil(t, renderContext, "a canceled request must not enqueue pre-render callbacks")
				require.Equal(t, http.StatusServiceUnavailable, w.Code)
				require.Empty(t, w.Body.String())
				return
			}
			require.NotNil(t, renderContext)
			require.Equal(t, "request value", renderContext.Value(requestKey{}))
			gotDeadline, ok := renderContext.Deadline()
			require.True(t, ok)
			require.Equal(t, deadline, gotDeadline)
			require.Equal(t, context.Canceled, renderContext.Err())
			if test.cancelDuring || test.cancelEncoding {
				if test.cancelDuring {
					require.Equal(t, context.Canceled, renderError)
				} else {
					require.NoError(t, renderError)
				}
				require.Equal(t, http.StatusServiceUnavailable, w.Code)
				require.Empty(t, w.Body.String())
			} else {
				require.NoError(t, renderError)
				require.NoError(t, parent.Err(), "render cleanup must not cancel the request context")
				require.Equal(t, http.StatusOK, w.Code)
				require.Contains(t, w.Body.String(), "<div>request context</div>")
			}
		})
	}
}

func TestHandlerCancelPendingRender(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{})
	release := make(chan struct{})
	workerFinished := make(chan struct{})
	finished := make(chan struct{})
	const path = "/cancel-pending-render-test"
	Route(path, func() Composer {
		return &requestContextTestCompo{onPreRender: func(ctx Context) {
			ctx.Async(func() {
				defer close(workerFinished)
				close(started)
				<-release // Deliberately ignores cancellation.
			})
		}}
	})
	t.Cleanup(func() {
		cancel()
		close(release)
		select {
		case <-finished:
		case <-time.After(5 * time.Second):
			t.Error("handler did not finish")
		}
		select {
		case <-started:
			select {
			case <-workerFinished:
			case <-time.After(5 * time.Second):
				t.Error("async worker did not finish after release")
			}
		default:
		}
		routes.mu.Lock()
		defer routes.mu.Unlock()
		delete(routes.routes, path)
	})
	r := httptest.NewRequest(http.MethodGet, path, nil).WithContext(ctx)
	w := httptest.NewRecorder()
	h := Handler{}
	go func() {
		h.ServeHTTP(w, r)
		close(finished)
	}()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("async pre-render work did not start")
	}
	cancel()
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("handler waited for canceled pre-render work")
	}
	select {
	case <-workerFinished:
		t.Fatal("worker should still be pending when the handler returns")
	default:
	}
	require.Empty(t, w.Body.String(), "canceled rendering must not write partial HTML")
	require.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestHandlerServePageWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: LocalDir(""),
		Title:     "Handler testing",
		Scripts: []string{
			"web/hello.js",
			"http://boo.com/bar.js",
		},
		Styles: []string{
			"web/foo.css",
			"/web/bar.css",
			"http://boo.com/bar.css",
		},
		RawHeaders: []string{
			`<meta http-equiv="refresh" content="30">`,
		},
		Image: "/web/test.png",
	}
	h.Icon.Maskable = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `<html lang="en">`)
	require.Contains(t, body, `href="/web/foo.css"`)
	require.Contains(t, body, `href="/web/bar.css"`)
	require.Contains(t, body, `href="http://boo.com/bar.css"`)
	require.Contains(t, body, `src="/web/hello.js"`)
	require.Contains(t, body, `src="http://boo.com/bar.js"`)
	require.Contains(t, body, `href="/manifest.webmanifest"`)
	require.Contains(t, body, `href="/app.css"`)
	require.Contains(t, body, `<meta http-equiv="refresh" content="30">`)
	require.Contains(t, body, `<div id="pre-render-ok">`)
	require.Contains(t, body, `content="https:///web/test.png"`)
	require.Contains(t, body, `<img src="/web/resolve-static-resource-test.jpg">`)

	t.Log(body)
}

func TestHandlerServePageWithRemoteBucket(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Title:     "Handler testing",
		Resources: RemoteBucket("https://storage.googleapis.com/go-app/"),
		Scripts: []string{
			"/web/hello.js",
			"http://boo.com/bar.js",
		},
		Styles: []string{
			"web/foo.css",
			"/web/bar.css",
			"http://boo.com/bar.css",
		},
		RawHeaders: []string{
			`<meta http-equiv="refresh" content="30">`,
		},
		Image: "/web/test.png",
	}
	h.Icon.Maskable = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `href="https://storage.googleapis.com/go-app/web/foo.css"`)
	require.Contains(t, body, `href="https://storage.googleapis.com/go-app/web/bar.css"`)
	require.Contains(t, body, `href="http://boo.com/bar.css"`)
	require.Contains(t, body, `src="https://storage.googleapis.com/go-app/web/hello.js"`)
	require.Contains(t, body, `src="http://boo.com/bar.js"`)
	require.Contains(t, body, `href="/manifest.webmanifest"`)
	require.Contains(t, body, `href="/app.css"`)
	require.Contains(t, body, `<meta http-equiv="refresh" content="30">`)
	require.Contains(t, body, `<div id="pre-render-ok">`)
	require.Contains(t, body, `content="https://storage.googleapis.com/go-app/web/test.png"`)
	require.Contains(t, body, `<img src="https://storage.googleapis.com/go-app/web/resolve-static-resource-test.jpg">`)

	t.Log(body)
}

func TestHandlerServePageWithGitHubPages(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Title:     "Handler testing",
		Resources: GitHubPages("go-app"),
		Scripts: []string{
			"/web/hello.js",
			"http://boo.com/bar.js",
		},
		Styles: []string{
			"/web/foo.css",
			"/web/bar.css",
			"http://boo.com/bar.css",
		},
		RawHeaders: []string{
			`<meta http-equiv="refresh" content="30">`,
		},
	}
	h.Icon.Maskable = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `href="/go-app/web/foo.css"`)
	require.Contains(t, body, `href="/go-app/web/bar.css"`)
	require.Contains(t, body, `href="http://boo.com/bar.css"`)
	require.Contains(t, body, `src="/go-app/web/hello.js"`)
	require.Contains(t, body, `src="http://boo.com/bar.js"`)
	require.Contains(t, body, `href="/go-app/manifest.webmanifest"`)
	require.Contains(t, body, `href="/go-app/app.css"`)
	require.Contains(t, body, `<meta http-equiv="refresh" content="30">`)
	require.Contains(t, body, `<div id="pre-render-ok">`)
	require.Contains(t, body, `<img src="/go-app/web/resolve-static-resource-test.jpg">`)
	t.Log(body)
}

func TestHandlerServeWasmExecJS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/wasm_exec.js", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Equal(t, wasmExecJS(), w.Body.String())
	t.Log(w.Body.String())
}

func TestHandlerServeAppJSWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		InternalURLs: []string{"https://redirect.me"},
	}
	h.ServeHTTP(w, r)
	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `"/app-worker.js"`)
	require.Contains(t, body, `fetchWithProgress("/web/app.wasm"`)
	require.Contains(t, body, "GOAPP_VERSION")
	require.Contains(t, body, `"GOAPP_STATIC_RESOURCES_URL":"/web"`)
	require.Contains(t, body, `"GOAPP_ROOT_PREFIX":"/"`)
	require.Contains(t, body, `"GOAPP_INTERNAL_URLS":"[\"https://redirect.me\"]"`)
}

func TestHandlerServeAppJSWithRemoteBucket(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: RemoteBucket("https://storage.googleapis.com/go-app/"),
	}
	h.ServeHTTP(w, r)
	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `"/app-worker.js"`)
	require.Contains(t, body, `fetchWithProgress("https://storage.googleapis.com/go-app/web/app.wasm"`)
	require.Contains(t, body, "GOAPP_VERSION")
	require.Contains(t, body, `"GOAPP_STATIC_RESOURCES_URL":"https://storage.googleapis.com/go-app/web"`)
	require.Contains(t, body, `"GOAPP_ROOT_PREFIX":"/"`)
}

func TestHandlerServeAppJSWithGitHubPages(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: GitHubPages("go-app"),
	}
	h.ServeHTTP(w, r)
	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `"/go-app/app-worker.js"`)
	require.Contains(t, body, `fetchWithProgress("/go-app/web/app.wasm"`)
	require.Contains(t, body, "GOAPP_VERSION")
	require.Contains(t, body, `"GOAPP_STATIC_RESOURCES_URL":"/go-app/web"`)
	require.Contains(t, body, `"GOAPP_ROOT_PREFIX":"/go-app"`)
}

func TestHandlerServeAppJSWithEnv(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Env: map[string]string{
			"FOO": "foo",
			"BAR": "bar",
		},
	}
	h.ServeHTTP(w, r)
	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, "GOAPP_VERSION")
	require.Contains(t, body, `"FOO":"foo"`)
	require.Contains(t, body, `"BAR":"bar"`)
	require.Contains(t, body, `"GOAPP_STATIC_RESOURCES_URL":"/web"`)
	require.Contains(t, body, `"GOAPP_ROOT_PREFIX":"/"`)
}

func TestHandlerServeAppWorkerJSWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app-worker.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Scripts: []string{"web/hello.js"},
		Styles:  []string{"/web/hello.css"},
		CacheableResources: []string{
			"web/hello.png",
			"http://test.io/hello.png",
		},
	}
	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `self.addEventListener("install", (event) => {`)
	require.Contains(t, body, `self.addEventListener("activate", (event) => {`)
	require.Contains(t, body, `self.addEventListener("fetch", (event) => {`)
	require.Contains(t, body, `"/web/hello.css"`)
	require.Contains(t, body, `"/web/hello.js"`)
	require.Contains(t, body, `"/web/hello.png"`)
	require.Contains(t, body, `"http://test.io/hello.png"`)
	require.Contains(t, body, `"/wasm_exec.js"`)
	require.Contains(t, body, `"/app.js"`)
	require.Contains(t, body, `"/web/app.wasm"`)
	require.Contains(t, body, `"/"`)
}

func TestHandlerServeAppWorkerJSWithRemoteBucket(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app-worker.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: RemoteBucket("https://storage.googleapis.com/go-app/"),
		Scripts:   []string{"web/hello.js"},
		Styles:    []string{"/web/hello.css"},
		CacheableResources: []string{
			"web/hello.png",
			"http://test.io/hello.png",
		},
	}
	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `self.addEventListener("install", (event) => {`)
	require.Contains(t, body, `self.addEventListener("activate", (event) => {`)
	require.Contains(t, body, `self.addEventListener("fetch", (event) => {`)
	require.Contains(t, body, `"https://storage.googleapis.com/go-app/web/hello.css"`)
	require.Contains(t, body, `"https://storage.googleapis.com/go-app/web/hello.js"`)
	require.Contains(t, body, `"https://storage.googleapis.com/go-app/web/hello.png"`)
	require.Contains(t, body, `"http://test.io/hello.png"`)
	require.Contains(t, body, `"/wasm_exec.js"`)
	require.Contains(t, body, `"/app.js"`)
	require.Contains(t, body, `"https://storage.googleapis.com/go-app/web/app.wasm"`)
	require.Contains(t, body, `"/"`)
}

func TestHandlerServeAppWorkerJSWithGitHubPages(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app-worker.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: GitHubPages("go-app"),
		Scripts:   []string{"web/hello.js"},
		Styles:    []string{"/web/hello.css"},
		CacheableResources: []string{
			"web/hello.png",
			"http://test.io/hello.png",
		},
	}
	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `self.addEventListener("install", (event) => {`)
	require.Contains(t, body, `self.addEventListener("activate", (event) => {`)
	require.Contains(t, body, `self.addEventListener("fetch", (event) => {`)
	require.Contains(t, body, `"/go-app/web/hello.css"`)
	require.Contains(t, body, `"/go-app/web/hello.js"`)
	require.Contains(t, body, `"/go-app/web/hello.png"`)
	require.Contains(t, body, `"http://test.io/hello.png"`)
	require.Contains(t, body, `"/go-app/wasm_exec.js"`)
	require.Contains(t, body, `"/go-app/app.js"`)
	require.Contains(t, body, `"/go-app/web/app.wasm"`)
	require.Contains(t, body, `"/go-app"`)
}

//go:embed testdata/app-worker.test.cjs
var appWorkerTests string

//go:embed testdata/script-config.test.cjs
var scriptConfigTests string

func TestJavaScriptTemplateSyntax(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node.js is required to check JavaScript template syntax")
	}
	for name, source := range map[string]string{"app.js": appJS, "app-worker.js": DefaultAppWorkerJS} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, node, "--check")
			cmd.Stdin = strings.NewReader(source)
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "%s", output)
		})
	}
}

func TestHandlerJavaScriptConfiguration(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node.js is required to exercise generated JavaScript")
	}
	for _, test := range []struct {
		name  string
		value string
	}{
		{"plain", "loading {progress}%"},
		{"empty", ""},
		{"quotes", "double\" single' backtick`"},
		{"backslashes", `C:\folder\file`},
		{"whitespace", "first\nsecond\r\n\tend\b\f"},
		{"unicode", "台灣 café 😀 \u2028\u2029"},
		{"markup", "</script><script>globalThis.injected = true</script>&"},
		{"injection", `"; globalThis.injected = true; //`},
		{"placeholders", "{{.Env}} {{.LoadingLabel}} {{.Wasm}} {{.WorkerJS}} {{.WasmContentLength}} {{.WasmContentLengthHeader}} {{.Version}} {{.ResourcesToCache}}"},
	} {
		t.Run(test.name, func(t *testing.T) {
			// makeAppJS exports these values to the server environment too.
			for _, key := range []string{"GOAPP_SCRIPT_TEST", "GOAPP_INTERNAL_URLS", "GOAPP_VERSION", "GOAPP_STATIC_RESOURCES_URL", "GOAPP_ROOT_PREFIX"} {
				t.Setenv(key, "")
			}
			h := Handler{
				Version:                 test.value,
				LoadingLabel:            test.value,
				WasmContentLength:       test.value,
				WasmContentLengthHeader: test.value,
				Env:                     map[string]string{"GOAPP_SCRIPT_TEST": test.value},
				InternalURLs:            []string{test.value},
				Resources:               PrefixedLocation("/" + test.value),
				ServiceWorkerTemplate:   DefaultAppWorkerJS,
			}
			appSource := string(h.makeAppJS())
			var resources []string
			for _, path := range []string{"/app.css", "/app.js", "/manifest.webmanifest", "/wasm_exec.js", "/web/app.wasm", "/"} {
				resources = append(resources, h.Resources.Resolve(path))
			}
			input := map[string]any{
				"app":       appSource,
				"worker":    string(h.makeAppWorkerJS()),
				"value":     test.value,
				"env":       h.Env,
				"wasmURL":   h.Resources.Resolve("/web/app.wasm"),
				"workerURL": h.Resources.Resolve("/app-worker.js"),
				"resources": resources,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, node, "-e", scriptConfigTests)
			cmd.Stdin = strings.NewReader(jsonString(input))
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "%s", output)
			require.Contains(t, string(output), "JavaScript configuration checks passed")
		})
	}
}

func TestAppWorkerLifecycleAndCaches(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node.js is required to exercise the service-worker script")
	}
	h := Handler{
		Version:               "v2",
		Resources:             LocalDir(""),
		ServiceWorkerTemplate: DefaultAppWorkerJS,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, "-", string(h.makeAppWorkerJS()))
	cmd.Stdin = strings.NewReader(appWorkerTests)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", output)
}

func TestHandlerCustomAppWorkerTemplate(t *testing.T) {
	h := Handler{
		Version:               "custom-version",
		Resources:             GitHubPages("custom-app"),
		ServiceWorkerTemplate: "version={{.Version}}; resources={{.ResourcesToCache}}",
	}
	script := string(h.makeAppWorkerJS())
	require.Contains(t, script, "version=custom-version; resources=[")
	require.Contains(t, script, `"/custom-app/web/app.wasm"`)
	require.NotContains(t, script, "{{.")
}

func TestHandlerCustomAppWorkerLiteralPlaceholders(t *testing.T) {
	h := Handler{
		Version:               "quotes\"\n{{.ResourcesToCache}} {{.Version}}",
		Resources:             LocalDir(""),
		ServiceWorkerTemplate: `raw={{.Version}}; literal="{{.Version}}"; resources={{.ResourcesToCache}}`,
	}
	script := string(h.makeAppWorkerJS())
	require.True(t, strings.HasPrefix(script, "raw="+h.Version+"; literal="+jsonString(h.Version)+"; resources=["))
	require.Contains(t, script, `"/web/app.wasm"`)
}

func TestHandlerServeManifestJSONWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/manifest.webmanifest", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Name:            "foobar",
		ShortName:       "foo",
		BackgroundColor: "#0000f0",
		ThemeColor:      "#0000ff",
	}

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/manifest+json", w.Header().Get("Content-Type"))
	require.Contains(t, body, `"short_name":"foo"`)
	require.Contains(t, body, `"name":"foobar"`)
	require.Contains(t, body, `"src":"https://raw.githubusercontent.com/maxence-charriere/go-app/master/docs/web/icon.png"`)
	require.Contains(t, body, `"background_color":"#0000f0"`)
	require.Contains(t, body, `"theme_color":"#0000ff"`)
	require.Contains(t, body, `"scope":"/"`)
	require.Contains(t, body, `"start_url":"/"`)
}

func TestHandlerServeManifestJSONWithRemoteBucket(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/manifest.webmanifest", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources:       RemoteBucket("https://storage.googleapis.com/go-app/"),
		Name:            "foobar",
		ShortName:       "foo",
		BackgroundColor: "#0000f0",
		ThemeColor:      "#0000ff",
	}

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/manifest+json", w.Header().Get("Content-Type"))
	require.Contains(t, body, `"short_name":"foo"`)
	require.Contains(t, body, `"name":"foobar"`)
	require.Contains(t, body, `"src":"https://raw.githubusercontent.com/maxence-charriere/go-app/master/docs/web/icon.png"`)
	require.Contains(t, body, `"background_color":"#0000f0"`)
	require.Contains(t, body, `"theme_color":"#0000ff"`)
	require.Contains(t, body, `"scope":"/"`)
	require.Contains(t, body, `"start_url":"/"`)
}

func TestHandlerServeManifestJSONWithGitHubPages(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/manifest.webmanifest", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources:       GitHubPages("go-app"),
		Name:            "foobar",
		ShortName:       "foo",
		BackgroundColor: "#0000f0",
		ThemeColor:      "#0000ff",
	}

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/manifest+json", w.Header().Get("Content-Type"))
	require.Contains(t, body, `"short_name":"foo"`)
	require.Contains(t, body, `"name":"foobar"`)
	require.Contains(t, body, `"src":"https://raw.githubusercontent.com/maxence-charriere/go-app/master/docs/web/icon.png"`)
	require.Contains(t, body, `"background_color":"#0000f0"`)
	require.Contains(t, body, `"theme_color":"#0000ff"`)
	require.Contains(t, body, `"scope":"/"`)
	require.Contains(t, body, `"start_url":"/go-app"`)
}

func TestHandlerServeAppCSS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.css", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "text/css", w.Header().Get("Content-Type"))
	require.Equal(t, appCSS, w.Body.String())
}

func TestHandlerServeAppWasm(t *testing.T) {
	close := testCreateDir(t, "web")
	defer close()
	testCreateFile(t, filepath.Join("web", "app.wasm"), "wasm!")

	h := Handler{}
	h.init()

	utests := []struct {
		scenario string
		path     string
	}{
		{
			scenario: "from resource provider path",
			path:     h.Resources.Resolve("/web/app.wasm"),
		},
		{
			scenario: "from legacy v6 path",
			path:     "/app.wasm",
		},
		{
			scenario: "from legacy v6 path",
			path:     "/goapp.wasm",
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, u.path, nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, r)
			require.Equal(t, "application/wasm", w.Header().Get("Content-Type"))
			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, "wasm!", w.Body.String())
		})
	}
}

func TestHandlerServeFile(t *testing.T) {
	close := testCreateDir(t, "web")
	defer close()
	testCreateFile(t, filepath.Join("web", "hello.txt"), "hello!")

	r := httptest.NewRequest(http.MethodGet, "/web/hello.txt", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "hello!", w.Body.String())
}

func TestHandlerProxyResources(t *testing.T) {
	close := testCreateDir(t, "web")
	defer close()

	s := httptest.NewServer(&Handler{
		ProxyResources: []ProxyResource{
			{
				Path:         "/hello.txt",
				ResourcePath: "/web/hello.txt",
			},
			{
				Path:         "/plop.txt",
				ResourcePath: "/web/plop.txt",
			},
			{
				Path:         "/app.js",
				ResourcePath: "/web/app.js",
			},
		},
	})
	defer s.Close()

	utests := []struct {
		scenario string
		file     string
		body     string
		code     int
		notProxy bool
	}{
		{
			scenario: "robots.txt is fetched",
			file:     "robots.txt",
			code:     http.StatusOK,
			body:     "robots!",
		},
		{
			scenario: "sitemap.xml is fetched",
			file:     "sitemap.xml",
			code:     http.StatusOK,
			body:     "sitemap!",
		},
		{
			scenario: "ads.txt is fetched",
			file:     "ads.txt",
			code:     http.StatusOK,
			body:     "ads!",
		},
		{
			scenario: "proxy resource is fetched",
			file:     "hello.txt",
			code:     http.StatusOK,
			body:     "hello!",
		},
		{
			scenario: "proxy resource is not found",
			file:     "plop.txt",
			code:     http.StatusNotFound,
		},
		{
			scenario: "no proxy resource is not fetched",
			file:     "bye.txt",
			code:     http.StatusNotFound,
			body:     "bye!",
			notProxy: true,
		},
		{
			scenario: "app.js is not a proxy resource",
			file:     "app.js",
			code:     http.StatusOK,
			body:     "wasm!",
			notProxy: true,
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			for i := 0; i < 2; i++ {
				if u.body != "" {
					testCreateFile(t, filepath.Join("web", u.file), u.body)
				}

				url := s.URL + "/" + u.file

				res, err := http.Get(url)
				require.NoError(t, err)
				defer res.Body.Close()

				require.Equal(t, u.code, res.StatusCode)
				if u.code != http.StatusOK {
					return
				}

				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				if u.notProxy {
					require.NotEqual(t, u.body, string(body))
					return
				}

				require.Equal(t, u.body, string(body))
			}
		})
	}
}

func BenchmarkHandlerColdRun(b *testing.B) {
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	for b.Loop() {
		h := Handler{}
		h.ServeHTTP(w, r)
		h.ServeHTTP(w, r)
	}
}

func BenchmarkHandlerHotRun(b *testing.B) {
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()
	h := Handler{}
	h.ServeHTTP(w, r)

	for b.Loop() {
		h.ServeHTTP(w, r)
	}
}

func TestIsRemoteLocation(t *testing.T) {
	tests := []struct {
		scenario string
		path     string
		expected bool
	}{
		{
			scenario: "path with http scheme is a remote location",
			path:     "http://localhost/hello",
			expected: true,
		},
		{
			scenario: "path with https scheme is a remote location",
			path:     "https://localhost/hello",
			expected: true,
		},
		{
			scenario: "empty path is not a remote location",
			path:     "",
			expected: false,
		},
		{
			scenario: "working dir path is not a remote location",
			path:     ".",
			expected: false,
		},
		{
			scenario: "absolute path is not a remote location",
			path:     "/User/hello",
			expected: false,
		},
		{
			scenario: "relative path is not a remote location",
			path:     "./hello",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			res := isRemoteLocation(test.path)
			require.Equal(t, test.expected, res)
		})
	}
}

func TestIsStaticResourcePath(t *testing.T) {
	tests := []struct {
		scenario string
		path     string
		expected bool
	}{
		{
			scenario: "static resource path",
			path:     "/web/hello",
			expected: true,
		},
		{
			scenario: "static resource path with prefix slash",
			path:     "web/hello",
			expected: true,
		},
		{
			scenario: "static resource directory",
			path:     "/web",
			expected: false,
		},
		{
			scenario: "static resource directory without prefix slash",
			path:     "web",
			expected: false,
		},
		{
			scenario: "non static resource",
			path:     "/app.js",
			expected: false,
		},
		{
			scenario: "remote resource",
			path:     "https://localhost/hello",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			res := isStaticResourcePath(test.path)
			require.Equal(t, test.expected, res)
		})
	}
}

func TestParseHTTPResource(t *testing.T) {
	utests := []struct {
		scenario string
		in       string
		resource httpResource
	}{
		{
			scenario: "empty",
			in:       "",
		},
		{
			scenario: "url is parsed",
			in:       "https://hello.world",
			resource: httpResource{URL: "https://hello.world"},
		},
		{
			scenario: "url and simple crossorigin is parsed",
			in:       "https://hello.world crossorigin",
			resource: httpResource{
				URL:         "https://hello.world",
				CrossOrigin: "true",
			},
		},
		{
			scenario: "url and defined crossorigin is parsed",
			in:       "https://hello.world crossorigin=use-credentials",
			resource: httpResource{
				URL:         "https://hello.world",
				CrossOrigin: "use-credentials",
			},
		},
		{
			scenario: "simple crossorigin is parsed",
			in:       "crossorigin",
			resource: httpResource{CrossOrigin: "true"},
		},
		{
			scenario: "defined crossorigin is parsed",
			in:       "crossorigin=anonymous",
			resource: httpResource{CrossOrigin: "anonymous"},
		},
		{
			scenario: "out of order url and simple crossorigin is parsed",
			in:       "    crossorigin    https://hello.world ",
			resource: httpResource{
				URL:         "https://hello.world",
				CrossOrigin: "true",
			},
		},
		{
			scenario: "url and async loading is parsed",
			in:       "https://hello.world async",
			resource: httpResource{
				URL:         "https://hello.world",
				LoadingMode: "async",
			},
		},
		{
			scenario: "url and defer loading is parsed",
			in:       "https://hello.world defer",
			resource: httpResource{
				URL:         "https://hello.world",
				LoadingMode: "defer",
			},
		},
		{
			scenario: "url with crossorigin and loading is parsed",
			in:       "https://hello.world defer crossorigin",
			resource: httpResource{
				URL:         "https://hello.world",
				CrossOrigin: "true",
				LoadingMode: "defer",
			},
		},
		{
			scenario: "url is case sensitive",
			in:       "https://Hello.World defer crossorigin",
			resource: httpResource{
				URL:         "https://Hello.World",
				CrossOrigin: "true",
				LoadingMode: "defer",
			},
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			require.Equal(t, u.resource, parseHTTPResource(u.in))
		})
	}
}

func TestHTTPResourceToLink(t *testing.T) {
	t.Run("resource is converted to link", func(t *testing.T) {
		link := parseHTTPResource("http://hello.world").toLink()

		require.IsType(t, Link(), link)
		require.Equal(t, "http://hello.world", link.(*htmlLink).attributes["href"])
		require.NotContains(t, link.(*htmlLink).attributes, "crossorigin")
	})

	t.Run("resource is converted to link with naked cross origin", func(t *testing.T) {
		link := parseHTTPResource("http://hello.world crossorigin").toLink()

		require.IsType(t, Link(), link)
		require.Equal(t, "http://hello.world", link.(*htmlLink).attributes["href"])
		require.Equal(t, "", link.(*htmlLink).attributes["crossorigin"])
	})

	t.Run("resource is converted to link with cross origin", func(t *testing.T) {
		link := parseHTTPResource("http://hello.world crossorigin=anonymous").toLink()

		require.IsType(t, Link(), link)
		require.Equal(t, "http://hello.world", link.(*htmlLink).attributes["href"])
		require.Equal(t, "anonymous", link.(*htmlLink).attributes["crossorigin"])
	})
}

func TestHTTPResourceToScript(t *testing.T) {
	t.Run("resource is converted to script", func(t *testing.T) {
		script := parseHTTPResource("http://hello.world").toScript()

		require.IsType(t, Script(), script)
		require.Equal(t, "http://hello.world", script.(*htmlScript).attributes["src"])
		require.NotContains(t, script.(*htmlScript).attributes, "crossorigin")
		require.NotContains(t, script.(*htmlScript).attributes, "async")
		require.NotContains(t, script.(*htmlScript).attributes, "defer")
	})

	t.Run("resource is converted to script with naked cross origin", func(t *testing.T) {
		script := parseHTTPResource("http://hello.world crossorigin").toScript()

		require.IsType(t, Script(), script)
		require.Equal(t, "http://hello.world", script.(*htmlScript).attributes["src"])
		require.Equal(t, "", script.(*htmlScript).attributes["crossorigin"])
	})

	t.Run("resource is converted to script with cross origin", func(t *testing.T) {
		script := parseHTTPResource("http://hello.world crossorigin=anonymous").toScript()

		require.IsType(t, Script(), script)
		require.Equal(t, "http://hello.world", script.(*htmlScript).attributes["src"])
		require.Equal(t, "anonymous", script.(*htmlScript).attributes["crossorigin"])
	})

	t.Run("resource is converted to script with defer loading", func(t *testing.T) {
		script := parseHTTPResource("http://hello.world defer").toScript()

		require.IsType(t, Script(), script)
		require.Equal(t, "http://hello.world", script.(*htmlScript).attributes["src"])
		require.Equal(t, "true", script.(*htmlScript).attributes["defer"])
	})

	t.Run("resource is converted to script with async loading", func(t *testing.T) {
		script := parseHTTPResource("http://hello.world async").toScript()

		require.IsType(t, Script(), script)
		require.Equal(t, "http://hello.world", script.(*htmlScript).attributes["src"])
		require.Equal(t, "true", script.(*htmlScript).attributes["async"])
	})
}
