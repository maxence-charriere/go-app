package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveStaticResourcePath(t *testing.T) {
	utests := []struct {
		scenario      string
		remoteRootDir string
		path          string
		expected      string
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
			scenario:      "non-static resource with remote root dir is skipped",
			remoteRootDir: "https://storage.googleapis.com/go-app",
			path:          "/hello",
			expected:      "/hello",
		},
		{
			scenario:      "non-static resource without slash and with remote root dir is skipped",
			remoteRootDir: "https://storage.googleapis.com/go-app",
			path:          "hello",
			expected:      "hello",
		},
		{
			scenario: "static resource is skipped",
			path:     "/web/hello.css",
			expected: "/web/hello.css",
		},
		{
			scenario: "static resource without slash is skipped",
			path:     "web/hello.css",
			expected: "web/hello.css",
		},
		{
			scenario:      "static resource with remote root dir is resolved",
			remoteRootDir: "https://storage.googleapis.com/go-app",
			path:          "/web/hello.css",
			expected:      "https://storage.googleapis.com/go-app/web/hello.css",
		},
		{
			scenario:      "static resource without slash and with remote root dir is resolved",
			remoteRootDir: "https://storage.googleapis.com/go-app",
			path:          "web/hello.css",
			expected:      "https://storage.googleapis.com/go-app/web/hello.css",
		},
		{
			scenario: "resolved static resource is skipped",
			path:     "https://storage.googleapis.com/go-app/web/hello.css",
			expected: "https://storage.googleapis.com/go-app/web/hello.css",
		},
		{
			scenario:      "resolved static resource with remote root dir is skipped",
			remoteRootDir: "https://storage.googleapis.com/go-app",
			path:          "https://storage.googleapis.com/go-app/web/hello.css",
			expected:      "https://storage.googleapis.com/go-app/web/hello.css",
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			remoteRootDir = u.remoteRootDir
			defer func() {
				remoteRootDir = ""
			}()

			res := StaticResource(u.path)
			require.Equal(t, u.expected, res)
		})
	}
}
