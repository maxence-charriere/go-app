//go:generate go run gen.go
//go:generate go fmt

package http

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/maxence-charriere/app/pkg/log"
)

var lastModified = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

// Page is a handler that serves the page that works with the wasm app.
type Page struct {
	Author       string
	Description  string
	Headers      []string
	Icon         string
	Keywords     []string
	LoadingLabel string
	Name         string
	ThemeColor   string
	WebDir       string
}

// CanHandle returns whether it can handle the given request.
func (p Page) CanHandle(r *http.Request) bool {
	return true
}

func (p Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer

	tmpl := template.Must(template.New("page").Parse(pageHTML))
	if err := tmpl.Execute(&b, struct {
		AppJS        string
		Author       string
		CSS          []string
		DefaultCSS   string
		Description  string
		Headers      []string
		Icon         string
		Keywords     string
		LoadingLabel string
		Name         string
		ThemeColor   string
		Scripts      []string
	}{
		AppJS:        pageJS,
		Author:       p.Author,
		CSS:          filepathsFromDir(p.WebDir, ".css"),
		DefaultCSS:   pageCSS,
		Description:  p.Description,
		Headers:      p.Headers,
		Icon:         p.Icon,
		Keywords:     strings.Join(p.Keywords, ", "),
		LoadingLabel: p.LoadingLabel,
		Name:         p.Name,
		ThemeColor:   p.ThemeColor,
		Scripts:      filepathsFromDir(p.WebDir, ".js"),
	}); err != nil {
		log.Error("generating page failed").
			T("error", err).
			T("path", r.URL.Path)
	}

	w.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Last-Modified", lastModified)
	w.Header().Set("Cache-Path-Override", "/")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Error("writing response failed").
			T("error", err).
			T("path", r.URL.Path)
	}
}

func filepathsFromDir(dirPath string, extensions ...string) []string {
	var filepaths []string

	extensionMap := make(map[string]struct{}, len(extensions))
	for _, ext := range extensions {
		extensionMap[ext] = struct{}{}
	}

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if _, ok := extensionMap[filepath.Ext(path)]; !ok {
			return nil
		}

		path = path[len(dirPath):]
		filepaths = append(filepaths, path)
		return nil
	}

	filepath.Walk(dirPath, walker)
	return filepaths
}
