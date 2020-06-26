// +build !wasm

package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

// GenerateStaticWebsite generates the files to run a PWA built with go-app as a
// static website in the specified directory. Static websites can be used with
// hosts such as Github Pages.
//
// Note that app.wasm must still be built separately and put into the web
// directory.
func GenerateStaticWebsite(dir string, h *Handler) error {
	if err := os.MkdirAll(filepath.Join(dir, "web"), 0755); err != nil {
		return errors.New("creating directory for static website failed").
			Tag("directory", dir).
			Wrap(err)
	}

	resources := []struct {
		filename string
		path     string
	}{
		{
			filename: "index.html",
			path:     "/",
		},
		{
			filename: "wasm_exec.js",
			path:     "/wasm_exec.js",
		},
		{
			filename: "app.js",
			path:     "/app.js",
		},
		{
			filename: "app-worker.js",
			path:     "/app-worker.js",
		},
		{
			filename: "manifest.json",
			path:     "/manifest.json",
		},
		{
			filename: "app.css",
			path:     "/app.css",
		},
	}

	for _, r := range resources {
		filename := filepath.Join(dir, r.filename)

		f, err := os.Create(filename)
		if err != nil {
			return errors.New("create file failed").
				Tag("filename", filename).
				Wrap(err)
		}
		defer f.Close()

		req, err := http.NewRequest(http.MethodGet, "http://go-app.io"+r.path, nil)
		if err != nil {
			return errors.New("creating file request failed").
				Tag("filename", filename).
				Tag("path", r.path).
				Wrap(err)
		}

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if n, err := f.Write(rec.Body.Bytes()); err != nil {
			return errors.New("writing file failed").
				Tag("filename", filename).
				Tag("bytes-written", n).
				Wrap(err)
		}

	}

	return nil
}
