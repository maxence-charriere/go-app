//go:build !wasm
// +build !wasm

package app

import (
	stderrors "errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
)

// GenerateStaticWebsite generates the files to run a PWA built with go-app as a
// static website in the specified directory. Static websites can be used with
// hosts such as Github Pages.
//
// Generation stops if a resource does not return HTTP 200 or cannot be saved.
// Files successfully generated before an error are retained.
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

			// Fetch before opening the file so a failed response cannot truncate it.
			page, err := createStaticPage(server.URL + path)
			if err != nil {
				return errors.New("creating page failed").
					WithTag("path", path).
					WithTag("filename", filename).
					Wrap(err)
			}
			f, err := createStaticFile(dir, filename)
			if err != nil {
				return errors.New("creating file failed").
					WithTag("path", path).
					WithTag("filename", filename).
					Wrap(err)
			}

			if err := writeStaticPage(f, page); err != nil {
				return errors.New("saving page failed").
					WithTag("path", path).
					WithTag("filename", filename).
					Wrap(err)
			}
		}
	}

	return nil
}

func createStaticDir(dir, path string) error {
	return os.MkdirAll(filepath.Join(dir, filepath.Dir(path)), 0755)
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

func writeStaticPage(f io.WriteCloser, page []byte) error {
	n, err := f.Write(page)
	if err == nil && n != len(page) {
		err = io.ErrShortWrite
	}
	if err != nil {
		err = errors.New("writing page failed").WithTag("bytes-written", n).Wrap(err)
	}
	if closeErr := f.Close(); closeErr != nil {
		err = stderrors.Join(err, errors.New("closing page file failed").Wrap(closeErr))
	}
	return err
}

func createStaticPage(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.New("creating http request failed").
			WithTag("path", path).
			Wrap(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("http request failed").
			WithTag("path", path).
			Wrap(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected http response status").
			WithTag("path", path).
			WithTag("status-code", res.StatusCode).
			WithTag("status", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("reading request body failed").
			WithTag("path", path).
			Wrap(err)
	}
	return body, nil
}
