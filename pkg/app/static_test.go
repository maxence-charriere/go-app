//go:build !wasm
// +build !wasm

package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateStaticWebsite(t *testing.T) {
	testSkipWasm(t)

	dir := "static-test"
	defer os.RemoveAll(dir)

	err := GenerateStaticWebsite(dir,
		&Handler{
			Name:      "Static Go-app",
			Title:     "Static test",
			Resources: GitHubPages("go-app"),
		},
		"/hello",
		"world",
		"/nested/foo",
	)
	require.NoError(t, err)

	files := []string{
		filepath.Join(dir),
		filepath.Join(dir, "web"),
		filepath.Join(dir, "index.html"),
		filepath.Join(dir, "wasm_exec.js"),
		filepath.Join(dir, "app.js"),
		filepath.Join(dir, "app-worker.js"),
		filepath.Join(dir, "manifest.webmanifest"),
		filepath.Join(dir, "app.css"),
		filepath.Join(dir, "hello.html"),
		filepath.Join(dir, "world.html"),
		filepath.Join(dir, "nested", "foo.html"),
	}

	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			_, err := os.Stat(f)
			require.NoError(t, err)
		})
	}
}
