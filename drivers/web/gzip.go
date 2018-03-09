// +build !js

package web

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWritter struct {
	res  http.ResponseWriter
	gzip *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) *gzipResponseWritter {
	w.Header().Set("Content-Encoding", "gzip")

	return &gzipResponseWritter{
		res:  w,
		gzip: gzip.NewWriter(w),
	}
}

func (w *gzipResponseWritter) Header() http.Header {
	return w.res.Header()
}

func (w *gzipResponseWritter) Write(b []byte) (int, error) {
	return w.gzip.Write(b)
}

func (w *gzipResponseWritter) WriteHeader(statusCode int) {
	w.res.WriteHeader(statusCode)
}

func (w *gzipResponseWritter) Close() error {
	return w.gzip.Close()
}

type gzipHandler struct {
	base http.Handler
}

func newGzipHandler(h http.Handler) http.Handler {
	return &gzipHandler{
		base: h,
	}
}

func (h *gzipHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if accept := req.Header.Get("Accept-Encoding"); !strings.Contains(accept, "gzip") {
		h.base.ServeHTTP(res, req)
		return
	}

	w := newGzipWriter(res)
	defer w.Close()

	res = w
	h.base.ServeHTTP(res, req)
}
