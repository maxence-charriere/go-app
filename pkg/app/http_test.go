// +build !wasm

package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerServePage(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Title: "Handler testing",
		Scripts: []string{
			"hello.js",
			"http://boo.com/bar.js",
		},
		Styles: []string{
			"foo.css",
			"/bar.css",
			"http://boo.com/bar.css",
		},
		Wasm: "test.wasm",
		Web:  "web",
	}
	h.Icon.AppleTouch = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `href="/web/foo.css"`)
	require.Contains(t, body, `href="/web/bar.css"`)
	require.Contains(t, body, `href="http://boo.com/bar.css"`)
	require.Contains(t, body, `<script src="/web/hello.js">`)
	require.Contains(t, body, `<script src="http://boo.com/bar.js">`)
	require.Contains(t, body, `href="/manifest.json"`)
	require.Contains(t, body, `href="/app.css"`)

	t.Log(body)
}

func TestHandlerServePageWithLocalWebDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Title: "Handler testing",
		Scripts: []string{
			"hello.js",
		},
		Styles: []string{
			"foo.css",
			"/bar.css",
		},
		Wasm: "test.wasm",
		Web:  "web",
	}
	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `href="/web/foo.css"`)
	require.Contains(t, body, `href="/web/bar.css"`)
	require.Contains(t, body, `<script src="/web/hello.js">`)
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

func TestHandlerServeAppJS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Equal(t, strings.ReplaceAll(appJS, "{{.Wasm}}", "/web/app.wasm"), w.Body.String())
}

func TestHandlerServeAppWorkerJS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app-worker.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Web:     "/web",
		Scripts: []string{"hello.js"},
		Styles:  []string{"hello.css"},
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
	require.Contains(t, body, `"/wasm_exec.js",`)
	require.Contains(t, body, `"/app.js",`)
	require.Contains(t, body, `"/web/app.wasm",`)
	require.Contains(t, body, `"/",`)
}

func TestHandlerServeManifestJSON(t *testing.T) {
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

func TestHandlerServeFile(t *testing.T) {
	err := os.MkdirAll(filepath.Join("test", "web"), 0755)
	require.NoError(t, err)
	defer os.RemoveAll("test")

	err = ioutil.WriteFile(filepath.Join("test", "web", "hello.txt"), []byte("hello!"), 0666)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/test/web/hello.txt", nil)
	w := httptest.NewRecorder()

	h := Handler{Web: "./test/web"}
	h.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "hello!", w.Body.String())
}
