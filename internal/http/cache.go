package http

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// CacheHandler returns a decorated version of the given cache that injects
// cache related headers.
func CacheHandler(h http.Handler, webDir string) http.Handler {
	return &cacheHandler{
		Handler: h,
		webDir:  webDir,
	}
}

type cacheHandler struct {
	http.Handler

	once         sync.Once
	etag         string
	cacheControl string
	webDir       string
}

func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.init)

	if r.URL.Path == "/.etag" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if h.etag == "" {
		h.Handler.ServeHTTP(w, r)
		return
	}

	w.Header().Set("ETag", h.etag)
	w.Header().Set("Cache-Control", "no-cache")

	etag := r.Header.Get("If-None-Match")
	if etag == h.etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	h.Handler.ServeHTTP(w, r)
}

func (h *cacheHandler) init() {
	if etag := GetEtag(h.webDir); etag != "" {
		h.etag = etagHeaderValue(etag)
	}
}

// GenerateEtag generates an etag.
func GenerateEtag() string {
	t := time.Now().UTC().String()
	return fmt.Sprintf(`%x`, sha1.Sum([]byte(t)))
}

// GetEtag returns the etag for the given web directory.
func GetEtag(webDir string) string {
	filename := filepath.Join(webDir, ".etag")

	etag, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(etag)
}

func etagHeaderValue(etag string) string {
	return `"` + etag + `"`
}
