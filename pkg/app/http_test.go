// +build !wasm

package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerServePageWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Title: "Handler testing",
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
	}
	h.Icon.AppleTouch = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `href="/web/foo.css"`)
	require.Contains(t, body, `href="/web/bar.css"`)
	require.Contains(t, body, `href="http://boo.com/bar.css"`)
	require.Contains(t, body, `src="/web/hello.js"`)
	require.Contains(t, body, `src="http://boo.com/bar.js"`)
	require.Contains(t, body, `href="/manifest.json"`)
	require.Contains(t, body, `href="/app.css"`)
	require.Contains(t, body, `<meta http-equiv="refresh" content="30">`)

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
	}
	h.Icon.AppleTouch = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `href="https://storage.googleapis.com/go-app/web/foo.css"`)
	require.Contains(t, body, `href="https://storage.googleapis.com/go-app/web/bar.css"`)
	require.Contains(t, body, `href="http://boo.com/bar.css"`)
	require.Contains(t, body, `src="https://storage.googleapis.com/go-app/web/hello.js"`)
	require.Contains(t, body, `src="http://boo.com/bar.js"`)
	require.Contains(t, body, `href="/manifest.json"`)
	require.Contains(t, body, `href="/app.css"`)
	require.Contains(t, body, `<meta http-equiv="refresh" content="30">`)

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
			"web/foo.css",
			"/web/bar.css",
			"http://boo.com/bar.css",
		},
		RawHeaders: []string{
			`<meta http-equiv="refresh" content="30">`,
		},
	}
	h.Icon.AppleTouch = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `href="/go-app/web/foo.css"`)
	require.Contains(t, body, `href="/go-app/web/bar.css"`)
	require.Contains(t, body, `href="http://boo.com/bar.css"`)
	require.Contains(t, body, `src="/go-app/web/hello.js"`)
	require.Contains(t, body, `src="http://boo.com/bar.js"`)
	require.Contains(t, body, `href="/go-app/manifest.json"`)
	require.Contains(t, body, `href="/go-app/app.css"`)
	require.Contains(t, body, `<meta http-equiv="refresh" content="30">`)

	t.Log(body)
}

func TestHandlerServeWasmExecJS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/wasm_exec.js", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Equal(t, wasmExecJS, w.Body.String())
}

func TestHandlerServeAppJSWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.ServeHTTP(w, r)
	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `register("/app-worker.js"`)
	require.Contains(t, body, `fetch("/web/app.wasm"`)
	require.Contains(t, body, "GOAPP_VERSION")
	require.Contains(t, body, `"GOAPP_STATIC_RESOURCES_URL":""`)
	require.Contains(t, body, `"GOAPP_ROOT_PREFIX":""`)
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
	require.Contains(t, body, `register("/app-worker.js"`)
	require.Contains(t, body, `fetch("https://storage.googleapis.com/go-app/web/app.wasm"`)
	require.Contains(t, body, "GOAPP_VERSION")
	require.Contains(t, body, `"GOAPP_STATIC_RESOURCES_URL":"https://storage.googleapis.com/go-app"`)
	require.Contains(t, body, `"GOAPP_ROOT_PREFIX":""`)
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
	require.Contains(t, body, `register("/go-app/app-worker.js"`)
	require.Contains(t, body, `fetch("/go-app/web/app.wasm"`)
	require.Contains(t, body, "GOAPP_VERSION")
	require.Contains(t, body, `"GOAPP_STATIC_RESOURCES_URL":"/go-app"`)
	require.Contains(t, body, `"GOAPP_ROOT_PREFIX":"/go-app"`)
}

func TestHandlerServeAppJSWithEnv(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Env: Environment{
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
	require.Contains(t, body, `"GOAPP_STATIC_RESOURCES_URL":""`)
	require.Contains(t, body, `"GOAPP_ROOT_PREFIX":""`)
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
	require.Contains(t, body, `self.addEventListener("install", event => {`)
	require.Contains(t, body, `self.addEventListener("activate", event => {`)
	require.Contains(t, body, `self.addEventListener("fetch", event => {`)
	require.Contains(t, body, `"/web/hello.css",`)
	require.Contains(t, body, `"/web/hello.js",`)
	require.Contains(t, body, `"/web/hello.png",`)
	require.Contains(t, body, `"http://test.io/hello.png",`)
	require.Contains(t, body, `"/wasm_exec.js",`)
	require.Contains(t, body, `"/app.js",`)
	require.Contains(t, body, `"/web/app.wasm",`)
	require.Contains(t, body, `"/",`)
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
	require.Contains(t, body, `self.addEventListener("install", event => {`)
	require.Contains(t, body, `self.addEventListener("activate", event => {`)
	require.Contains(t, body, `self.addEventListener("fetch", event => {`)
	require.Contains(t, body, `"https://storage.googleapis.com/go-app/web/hello.css",`)
	require.Contains(t, body, `"https://storage.googleapis.com/go-app/web/hello.js",`)
	require.Contains(t, body, `"https://storage.googleapis.com/go-app/web/hello.png",`)
	require.Contains(t, body, `"http://test.io/hello.png",`)
	require.Contains(t, body, `"/wasm_exec.js",`)
	require.Contains(t, body, `"/app.js",`)
	require.Contains(t, body, `"https://storage.googleapis.com/go-app/web/app.wasm",`)
	require.Contains(t, body, `"/",`)
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
	require.Contains(t, body, `self.addEventListener("install", event => {`)
	require.Contains(t, body, `self.addEventListener("activate", event => {`)
	require.Contains(t, body, `self.addEventListener("fetch", event => {`)
	require.Contains(t, body, `"/go-app/web/hello.css",`)
	require.Contains(t, body, `"/go-app/web/hello.js",`)
	require.Contains(t, body, `"/go-app/web/hello.png",`)
	require.Contains(t, body, `"http://test.io/hello.png",`)
	require.Contains(t, body, `"/go-app/wasm_exec.js",`)
	require.Contains(t, body, `"/go-app/app.js",`)
	require.Contains(t, body, `"/go-app/web/app.wasm",`)
	require.Contains(t, body, `"/go-app",`)
}

func TestHandlerServeManifestJSONWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/manifest.json", nil)
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
	require.Contains(t, body, `"short_name": "foo"`)
	require.Contains(t, body, `"name": "foobar"`)
	require.Contains(t, body, `"src": "https://storage.googleapis.com/murlok-github/icon-192.png"`)
	require.Contains(t, body, `"src": "https://storage.googleapis.com/murlok-github/icon-512.png"`)
	require.Contains(t, body, `"background_color": "#0000f0"`)
	require.Contains(t, body, `"theme_color": "#0000ff"`)
	require.Contains(t, body, `"scope": "/"`)
	require.Contains(t, body, `"start_url": "/"`)
}

func TestHandlerServeManifestJSONWithRemoteBucket(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/manifest.json", nil)
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
	require.Contains(t, body, `"short_name": "foo"`)
	require.Contains(t, body, `"name": "foobar"`)
	require.Contains(t, body, `"src": "https://storage.googleapis.com/murlok-github/icon-192.png"`)
	require.Contains(t, body, `"src": "https://storage.googleapis.com/murlok-github/icon-512.png"`)
	require.Contains(t, body, `"background_color": "#0000f0"`)
	require.Contains(t, body, `"theme_color": "#0000ff"`)
	require.Contains(t, body, `"scope": "/"`)
	require.Contains(t, body, `"start_url": "/"`)
}

func TestHandlerServeManifestJSONWithGitHubPages(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/manifest.json", nil)
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
	require.Contains(t, body, `"short_name": "foo"`)
	require.Contains(t, body, `"name": "foobar"`)
	require.Contains(t, body, `"src": "https://storage.googleapis.com/murlok-github/icon-192.png"`)
	require.Contains(t, body, `"src": "https://storage.googleapis.com/murlok-github/icon-512.png"`)
	require.Contains(t, body, `"background_color": "#0000f0"`)
	require.Contains(t, body, `"theme_color": "#0000ff"`)
	require.Contains(t, body, `"scope": "/go-app/"`)
	require.Contains(t, body, `"start_url": "/go-app/"`)
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
			path:     h.Resources.AppWASM(),
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

func TestHandlerServeRobotsTxt(t *testing.T) {
	close := testCreateDir(t, "web")
	defer close()
	testCreateFile(t, filepath.Join("web", "robots.txt"), "robot")

	s := httptest.NewServer(&Handler{})
	defer s.Close()

	test := func(t *testing.T) {
		res, err := http.Get(s.URL + "/robots.txt")
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)

		content, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "robot", btos(content))
	}

	t.Run("robots.txt", test)
	t.Run("cached robots.txt", test)
}

func TestHandlerServeRobotsTxtNotFound(t *testing.T) {
	s := httptest.NewServer(&Handler{})
	defer s.Close()

	res, err := http.Get(s.URL + "/robots.txt")
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func BenchmarkHandlerColdRun(b *testing.B) {
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
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

	for i := 0; i < b.N; i++ {
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
