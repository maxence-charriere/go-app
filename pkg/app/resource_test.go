//go:build !wasm
// +build !wasm

package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalDir(t *testing.T) {
	testSkipWasm(t)

	h, _ := LocalDir("test").(localResourceResolver)
	require.Equal(t, "/", h.Resolve(""))
	require.Equal(t, "/", h.Resolve("/"))
	require.Equal(t, "test/web/app.wasm", h.Resolve("/web/app.wasm"))

	close := testCreateDir(t, "test/web")
	defer close()

	resources := []string{
		"/web/test",
		"/web/app.wasm",
	}

	for _, r := range resources {
		t.Run(r, func(t *testing.T) {
			path := strings.Replace(r, "/web", "test/web", 1)
			err := ioutil.WriteFile(path, []byte("hello"), 0666)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, r, nil)
			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)
			require.Equal(t, "hello", res.Body.String())
		})
	}
}

func TestRemoteBucket(t *testing.T) {
	utests := []struct {
		scenario string
		provider ResourceResolver
	}{
		{
			scenario: "remote bucket",
			provider: RemoteBucket("https://storage.googleapis.com/test"),
		},
		{
			scenario: "remote bucket with / suffix",
			provider: RemoteBucket("https://storage.googleapis.com/test/"),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			require.Equal(t, "/", u.provider.Resolve(""))
			require.Equal(t, "/", u.provider.Resolve("/"))
			require.Equal(t, "https://storage.googleapis.com/test/web/app.wasm", u.provider.Resolve("/web/app.wasm"))
		})
	}
}
