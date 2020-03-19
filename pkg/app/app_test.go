package app

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveStaticResourcePath(t *testing.T) {
	tests := []struct {
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

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			remoteRootDir = test.remoteRootDir
			defer func() {
				remoteRootDir = ""
			}()

			res := ResolveStaticResourcePath(test.path)
			require.Equal(t, test.expected, res)
		})
	}
}

func BenchmarkHTML(b *testing.B) {
	var buffer bytes.Buffer

	root := Div().
		Class("Menu").
		Body(
			Text("☰"),
			Main().
				Class("Content").
				Body(
					H1().
						Body(
							Text("Hello,"),
							Text("world"),
						),
					Input().
						Placeholder("What is your name?").
						AutoFocus(true),
				),
		)

	for i := 0; i < b.N; i++ {
		buffer.Reset()
		root.html(&buffer)
	}
}

func BenchmarkHTMLWithRoot(b *testing.B) {
	var buffer bytes.Buffer

	for i := 0; i < b.N; i++ {
		buffer.Reset()

		root := Div().
			Class("Menu").
			Body(
				Text("☰"),
				Main().
					Class("Content").
					Body(
						H1().
							Body(
								Text("Hello,"),
								Text("world"),
							),
						Input().
							Placeholder("What is your name?").
							AutoFocus(true),
					),
			)

		root.html(&buffer)
	}
}
