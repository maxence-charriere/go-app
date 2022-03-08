//go:build !wasm
// +build !wasm

package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
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

	resources := map[string]struct{}{
		"/":                     {},
		"/wasm_exec.js":         {},
		"/app.js":               {},
		"/app-worker.js":        {},
		"/manifest.webmanifest": {},
		"/app.css":              {},
		"/web":                  {},
	}

	for path := range routes.routes {
		resources[path] = struct{}{}
	}

	for _, p := range pages {
		if p == "" {
			continue
		}
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		resources[p] = struct{}{}
	}

	server := httptest.NewServer(h)
	defer server.Close()

	for path := range resources {
		switch path {
		case "/web":
			if err := createStaticDir(filepath.Join(dir, path), ""); err != nil {
				return errors.New("creating web directory failed").Wrap(err)
			}

		default:
			filename := path
			if filename == "/" {
				filename = "/index.html"
			}

			f, err := createStaticFile(dir, filename)
			if err != nil {
				return errors.New("creating file failed").
					Tag("path", path).
					Tag("filename", filename).
					Wrap(err)
			}
			defer f.Close()

			page, err := createStaticPage(server.URL + path)
			if err != nil {
				return errors.New("creating page failed").
					Tag("path", path).
					Tag("filename", filename).
					Wrap(err)
			}

			if n, err := f.Write(page); err != nil {
				return errors.New("writing page failed").
					Tag("path", path).
					Tag("filename", filename).
					Tag("bytes-written", n).
					Wrap(err)
			}
		}
	}

	return nil
}

func createStaticDir(dir, path string) error {
	dir = filepath.Join(dir, filepath.Dir(path))
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}
	return os.MkdirAll(filepath.Join(dir), 0755)
}

func createStaticFile(dir, path string) (*os.File, error) {
	if err := createStaticDir(dir, path); err != nil {
		return nil, errors.New("creating file directory failed").Wrap(err)
	}

	filename := filepath.Join(dir, path)
	if filepath.Ext(filename) == "" {
		filename += ".html"
	}

	return os.Create(filename)
}

func createStaticPage(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.New("creating http request failed").
			Tag("path", path).
			Wrap(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("http request failed").
			Tag("path", path).
			Wrap(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("reading request body failed").
			Tag("path", path).
			Wrap(err)
	}
	return body, nil
}
