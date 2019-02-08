// +build !wasm

package app

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	gzip *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) *gzipResponseWriter {
	w.Header().Set("Content-Encoding", "gzip")

	return &gzipResponseWriter{
		ResponseWriter: w,
		gzip:           gzip.NewWriter(w),
	}
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.gzip.Write(b)
}

func (w *gzipResponseWriter) Close() error {
	return w.gzip.Close()
}

type gzipHandler struct {
	http.Handler
}

func newGzipHandler(h http.Handler) http.Handler {
	return &gzipHandler{Handler: h}
}

func (h *gzipHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if accept := req.Header.Get("Accept-Encoding"); !strings.Contains(accept, "gzip") {
		h.Handler.ServeHTTP(res, req)
		return
	}

	w := newGzipWriter(res)
	defer w.Close()

	h.Handler.ServeHTTP(w, req)
}
