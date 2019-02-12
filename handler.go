// +build !wasm

//go:generate go run page_gen.go
//go:generate go fmt

package app

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

// Handler is a http handler that serves UI components created with this
// package.
type Handler struct {
	// The app author.
	Author string

	// The app description.
	Description string

	// The path of the icon relative to the web directory. Default is
	// "/logo.png".
	Icon string

	// The app keywords.
	Keywords []string

	// The app name.
	Name string

	// The path of the go web assembly file to serve.
	Wasm string

	// The function that returns the path of the web directory. Default is the
	// working directory.
	WebDir func() string

	once        sync.Once
	fileHandler http.Handler
	page        []byte
	webDir      string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.init)

	path := filepath.Join(h.webDir, r.URL.Path)

	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		h.fileHandler.ServeHTTP(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(h.page)
}

func (h *Handler) init() {
	if h.Wasm == "" {
		wasm, _ := os.Getwd()
		wasm = filepath.Base(wasm)
		if wasm == "" {
			wasm = "app"
		}

		h.Wasm = "/" + wasm + ".wasm"
	}

	if h.Icon == "" {
		h.Icon = "/logo.png"
	}

	h.webDir = "."
	if h.WebDir != nil {
		h.webDir = h.WebDir()
	}
	h.webDir, _ = filepath.Abs(h.webDir)

	h.fileHandler = h.newFileHandler()
	h.page = h.newPage()
}

func (h *Handler) newFileHandler() http.Handler {
	handler := http.FileServer(http.Dir(h.webDir))
	handler = newGzipHandler(handler)
	return handler
}

func (h *Handler) newPage() []byte {
	b := bytes.Buffer{}
	tmpl := template.Must(template.New("page").Parse(pageHTML))

	if err := tmpl.Execute(&b, struct {
		Author      string
		CSS         []string
		DefaultCSS  string
		Description string
		Icon        string
		Keywords    string
		Name        string
		Scripts     []string
		Wasm        string
	}{
		Author:      h.Author,
		CSS:         h.filepathsFromDir(h.webDir, ".css"),
		DefaultCSS:  pageCSS,
		Description: h.Description,
		Icon:        h.Icon,
		Keywords:    strings.Join(h.Keywords, ", "),
		Name:        h.Name,
		Scripts:     h.filepathsFromDir(h.webDir, ".js"),
		Wasm:        h.Wasm,
	}); err != nil {
		panic(err)
	}

	return b.Bytes()
}

func (h *Handler) filepathsFromDir(dirPath string, extensions ...string) []string {
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

		path = path[len(h.webDir):]
		filepaths = append(filepaths, path)
		return nil
	}

	filepath.Walk(dirPath, walker)
	return filepaths
}
