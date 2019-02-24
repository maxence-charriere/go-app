//go:generate go run page_gen.go
//go:generate go fmt

package http

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

// PageHandler is a handler that serves pages that works with wasm app.
type PageHandler struct {
	Author       string
	Description  string
	Icon         string
	Keywords     []string
	LoadingLabel string
	Name         string
	WebDir       string

	once         sync.Once
	lastModified string
	page         []byte
}

func (h *PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.initPage)

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Last-Modified", h.lastModified)
	w.Write(h.page)
}

func (h *PageHandler) initPage() {
	buffer := bytes.Buffer{}
	writer := gzip.NewWriter(&buffer)

	tmpl := template.Must(template.New("page").Parse(pageHTML))
	if err := tmpl.Execute(writer, struct {
		AppJS        string
		Author       string
		CSS          []string
		DefaultCSS   string
		Description  string
		Icon         string
		Keywords     string
		LoadingLabel string
		Name         string
		Scripts      []string
	}{
		AppJS:        pageJS,
		Author:       h.Author,
		CSS:          filepathsFromDir(h.WebDir, ".css"),
		DefaultCSS:   pageCSS,
		Description:  h.Description,
		Icon:         h.Icon,
		Keywords:     strings.Join(h.Keywords, ", "),
		LoadingLabel: h.LoadingLabel,
		Name:         h.Name,
		Scripts:      filepathsFromDir(h.WebDir, ".js"),
	}); err != nil {
		panic(err)
	}

	writer.Close()
	h.page = buffer.Bytes()
	h.lastModified = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
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
