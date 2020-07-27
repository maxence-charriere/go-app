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

	utests := []struct {
		scenario string
		provider ResourceProvider
	}{
		{
			scenario: "from web directory",
			provider: LocalDir("web"),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			h := u.provider.(localDir)

			require.Empty(t, h.StaticResources())
			require.Equal(t, "/web/app.wasm", h.AppWASM())
			require.Equal(t, "/web/robots.txt", h.RobotsTxt())
			require.Equal(t, "/web/ads.txt", h.AdsTxt())

			close := testCreateDir(t, "web")
			defer close()

			resources := []string{
				"/web/test",
				h.AppWASM(),
				h.RobotsTxt(),
			}

			for _, r := range resources {
				t.Run(r, func(t *testing.T) {
					path := strings.Replace(r, "/web", h.path, 1)
					err := ioutil.WriteFile(path, stob("hello"), 0666)
					require.NoError(t, err)

					req := httptest.NewRequest(http.MethodGet, r, nil)
					res := httptest.NewRecorder()
					h.ServeHTTP(res, req)
					require.Equal(t, "hello", res.Body.String())
				})
			}
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
			require.Equal(t, "https://storage.googleapis.com/test", u.provider.StaticResources())
			require.Equal(t, "https://storage.googleapis.com/test/web/app.wasm", u.provider.AppWASM())
			require.Equal(t, "https://storage.googleapis.com/test/web/robots.txt", u.provider.RobotsTxt())
			require.Equal(t, "https://storage.googleapis.com/test/web/ads.txt", u.provider.AdsTxt())
		})
	}
}
