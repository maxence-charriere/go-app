// +build !wasm

//go:generate go run page_gen.go
//go:generate go fmt

package app

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

// Handler is a http handler that serves UI components created with this
// package.
type Handler struct {
	// The app author.
	Author string

	// The app description.
	Description string

	// The path of the icon relative to the web directory.
	Icon string

	// The app keywords.
	Keywords []string

	// The text displayed while loading the a page.
	Loading string

	// The app name.
	Name string

	// The path of the go web assembly file to serve relative to the web
	// directory.
	Wasm string

	// The he path of the web directory. Default is the working directory.
	WebDir string

	// WebDirFunc is a func that returns the path of the web directory. The
	// returned string overrides WebDir when defined.
	WebDirFunc func() string

	once         sync.Once
	fileHandler  http.Handler
	lastModified string
	page         []byte
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.init)

	path := filepath.Join(h.WebDir, r.URL.Path)

	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		h.fileHandler.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Last-Modified", h.lastModified)
	w.Header().Set("Context-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(h.page)
}

func (h *Handler) init() {
	h.Wasm = h.getWasm()
	h.WebDir = h.getWebDir()
	h.fileHandler = h.newFileHandler(h.WebDir)
	h.lastModified = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	h.page = h.newPage()
}

func (h *Handler) getWasm() string {
	wasm := h.Wasm
	if !strings.HasSuffix(wasm, ".wasm") {
		wasm += ".wasm"
	}
	return "/" + wasm
}

func (h *Handler) getWebDir() string {
	webdir := h.WebDir
	if h.WebDirFunc != nil {
		webdir = h.WebDirFunc()
	}
	if webdir == "" {
		webdir = "."
	}
	webdir, _ = filepath.Abs(webdir)
	return webdir
}

func (h *Handler) newFileHandler(webDir string) http.Handler {
	handler := http.FileServer(http.Dir(webDir))
	handler = newGzipHandler(handler)
	return handler
}

func (h *Handler) newPage() []byte {
	b := bytes.Buffer{}
	gz := gzip.NewWriter(&b)
	defer gz.Close()

	tmpl := template.Must(template.New("page").Parse(pageHTML))

	if err := tmpl.Execute(gz, struct {
		AppJS       string
		Author      string
		CSS         []string
		DefaultCSS  string
		Description string
		Icon        string
		Keywords    string
		Loading     string
		Name        string
		Scripts     []string
		Wasm        string
	}{
		AppJS:       pageJS,
		Author:      h.Author,
		CSS:         h.filepathsFromDir(h.WebDir, ".css"),
		DefaultCSS:  pageCSS,
		Description: h.Description,
		Icon:        h.Icon,
		Keywords:    strings.Join(h.Keywords, ", "),
		Loading:     h.Loading,
		Name:        h.Name,
		Scripts:     h.filepathsFromDir(h.WebDir, ".js"),
		Wasm:        h.Wasm,
	}); err != nil {
		panic(err)
	}

	gz.Close()
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

		path = path[len(dirPath):]
		filepaths = append(filepaths, path)
		return nil
	}

	filepath.Walk(dirPath, walker)
	return filepaths
}
