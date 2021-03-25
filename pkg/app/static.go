// +build !wasm

package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

// GenerateStaticWebsite generates the files to run a PWA built with go-app as a
// static website in the specified directory. Static websites can be used with
// hosts such as Github Pages.
//
// Note that app.wasm must still be built separately and put into the web
// directory.
func GenerateStaticWebsite(dir string, h *Handler, pages ...string) error {
	if dir == "" {
		dir = "."
	}

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
			filename: "manifest.webmanifest",
			path:     "/manifest.webmanifest",
		},
		{
			filename: "app.css",
			path:     "/app.css",
		},
	}

	for _, p := range pages {
		if p == "" {
			continue
		}

		resources = append(resources, struct {
			filename string
			path     string
		}{
			filename: p + ".html",
			path:     p,
		})
	}

	server := httptest.NewServer(h)
	defer server.Close()

	for _, r := range resources {
		filename := filepath.Join(dir, r.filename)

		f, err := os.Create(filename)
		if err != nil {
			return errors.New("create file failed").
				Tag("filename", filename).
				Wrap(err)
		}
		defer f.Close()

		req, err := http.NewRequest(http.MethodGet, server.URL+r.path, nil)
		if err != nil {
			return errors.New("creating file request failed").
				Tag("filename", filename).
				Tag("path", r.path).
				Wrap(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return errors.New("http request failed").
				Tag("filename", filename).
				Tag("path", r.path).
				Wrap(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return errors.New(res.Status)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New("reading request body failed").
				Tag("filename", filename).
				Tag("path", r.path).
				Wrap(err)
		}

		if n, err := f.Write(body); err != nil {
			return errors.New("writing file failed").
				Tag("filename", filename).
				Tag("bytes-written", n).
				Wrap(err)
		}
	}

	return nil
}
