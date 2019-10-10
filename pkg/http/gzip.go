package http

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

// DefaultGzipContentTypes returns the mime types that are gzipped by default.
func DefaultGzipContentTypes() []string {
	return []string{
		"application/javascript",
		"application/json",
		"application/wasm",
		"application/x-javascript",
		"application/x-tar",
		"image/svg+xml",
		"text/css",
		"text/html",
		"text/plain",
		"text/xml",
	}
}

// Gzip returns a decorated version of the given handler that gzip responses
// bodies with the given content types. It uses DefaultGzipContentTypes when
// there is no content types specified.
func Gzip(h http.Handler, contentTypes ...string) http.Handler {
	if len(contentTypes) == 0 {
		contentTypes = DefaultGzipContentTypes()
	}

	return &zip{
		handler:      h,
		contentTypes: contentTypes,
	}
}

type zip struct {
	handler      http.Handler
	contentTypes []string
}

func (z *zip) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var writer io.Writer = w
	var once sync.Once
	var gz *gzip.Writer

	defer func() {
		if gz != nil {
			gz.Close()
		}
	}()

	proxy := proxyWriter{
		header: w.Header,
		write: func(b []byte) (int, error) {
			once.Do(func() {
				if z.isCompressible(w.Header().Get("Content-Type")) {
					gz = gzip.NewWriter(w)
					writer = gz
					w.Header().Set("Content-Encoding", "gzip")
				}
			})
			return writer.Write(b)
		},
		writeHeader: w.WriteHeader,
	}

	z.handler.ServeHTTP(proxy, r)
}

func (z *zip) isCompressible(contentType string) bool {
	if contentType == "" {
		return false
	}

	for _, t := range z.contentTypes {
		if strings.Contains(contentType, t) {
			return true
		}
	}

	return false
}
