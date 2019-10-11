package http

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/maxence-charriere/app/pkg/log"
)

// Gzip returns a decorated version of the given handler that gzip responses
// bodies.
func Gzip(h http.Handler) http.Handler {
	return &zip{
		handler: h,
	}
}

type zip struct {
	handler      http.Handler
	contentTypes []string
}

func (z *zip) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	if !strings.Contains(acceptEncoding, "gzip") {
		z.handler.ServeHTTP(w, r)
		return
	}

	log.Info("gzipping").
		T("path", r.URL.Path)

	gz := gzip.NewWriter(w)
	proxy := proxyWriter{
		header:      w.Header,
		write:       gz.Write,
		writeHeader: w.WriteHeader,
		close:       gz.Close,
	}

	w.Header().Set("Content-Encoding", "gzip")
	z.handler.ServeHTTP(&proxy, r)
	proxy.Close()
}
