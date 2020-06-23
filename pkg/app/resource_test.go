package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStaticResource(t *testing.T) {
	utests := []struct {
		scenario           string
		staticResourcesURL string
		path               string
		expected           string
	}{
		{
			scenario: "non-static resource is skipped",
			path:     "/hello",
			expected: "/hello",
		},
		{
			scenario: "non-static resource without slash is skipped",
			path:     "hello",
			expected: "hello",
		},
		{
			scenario:           "non-static resource with remote root dir is skipped",
			staticResourcesURL: "https://storage.googleapis.com/go-app",
			path:               "/hello",
			expected:           "/hello",
		},
		{
			scenario:           "non-static resource without slash and with remote root dir is skipped",
			staticResourcesURL: "https://storage.googleapis.com/go-app",
			path:               "hello",
			expected:           "hello",
		},
		{
			scenario: "static resource",
			path:     "/web/hello.css",
			expected: "/web/hello.css",
		},
		{
			scenario: "static resource without slash",
			path:     "web/hello.css",
			expected: "/web/hello.css",
		},
		{
			scenario:           "static resource with remote root dir is resolved",
			staticResourcesURL: "https://storage.googleapis.com/go-app",
			path:               "/web/hello.css",
			expected:           "https://storage.googleapis.com/go-app/web/hello.css",
		},
		{
			scenario:           "static resource without slash and with remote root dir is resolved",
			staticResourcesURL: "https://storage.googleapis.com/go-app",
			path:               "web/hello.css",
			expected:           "https://storage.googleapis.com/go-app/web/hello.css",
		},
		{
			scenario: "resolved static resource is skipped",
			path:     "https://storage.googleapis.com/go-app/web/hello.css",
			expected: "https://storage.googleapis.com/go-app/web/hello.css",
		},
		{
			scenario:           "resolved static resource with remote root dir is skipped",
			staticResourcesURL: "https://storage.googleapis.com/go-app",
			path:               "https://storage.googleapis.com/go-app/web/hello.css",
			expected:           "https://storage.googleapis.com/go-app/web/hello.css",
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			staticResourcesURL = u.staticResourcesURL
			defer func() {
				staticResourcesURL = ""
			}()

			res := StaticResource(u.path)
			require.Equal(t, u.expected, res)
		})
	}
}

func TestLocalDir(t *testing.T) {
	utests := []struct {
		scenario string
		provider ResourceProvider
	}{
		{
			scenario: "from working directory",
			provider: LocalDir("."),
		},
		{
			scenario: "from web directory",
			provider: LocalDir("web"),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			h := u.provider.(localDir)

			require.Empty(t, h.URL())
			require.Equal(t, "/web/app.wasm", h.AppWASM())
			require.Equal(t, "/web/robot.txt", h.RobotTxt())

			err := os.MkdirAll(h.path, 0755)
			require.NoError(t, err)
			defer os.RemoveAll(h.path)

			resources := []string{
				"/web/test",
				h.AppWASM(),
				h.RobotTxt(),
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
			require.Equal(t, "https://storage.googleapis.com/test", u.provider.URL())
			require.Equal(t, "https://storage.googleapis.com/test/web/app.wasm", u.provider.AppWASM())
			require.Equal(t, "https://storage.googleapis.com/test/web/robot.txt", u.provider.RobotTxt())
		})
	}
}
