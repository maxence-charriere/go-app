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

	h, _ := LocalDir("test").(localDir)
	require.Equal(t, "test", h.Static())
	require.Equal(t, "test/web/app.wasm", h.AppWASM())

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
		provider ResourceProvider
	}{
		{
			scenario: "remote bucket",
			provider: RemoteBucket("https://storage.googleapis.com/test"),
		},
		{
			scenario: "remote bucket with web suffix",
			provider: RemoteBucket("https://storage.googleapis.com/test/web/"),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			require.Equal(t, "https://storage.googleapis.com/test", u.provider.Static())
			require.Equal(t, "https://storage.googleapis.com/test/web/app.wasm", u.provider.AppWASM())
		})
	}
}
