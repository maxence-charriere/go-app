package http

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// Gzip returns a decorated version of the given handler that gzip responses
// bodies whit the given content-types.
func Gzip(h http.Handler, contentTypes ...string) http.Handler {
	if len(contentTypes) == 0 {
		contentTypes = DefaultContentTypes
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
	acceptEncoding := r.Header.Get("Accept-Encoding")
	if !strings.Contains(acceptEncoding, "gzip") {
		z.handler.ServeHTTP(w, r)
		return
	}

	proxy := ProxyWriter{
		Writer:       w,
		HeaderWriter: w,
		BeforeWrite: func(p *ProxyWriter) {
			contentType := p.HeaderWriter.Header().Get("Content-Type")
			if isCacheableOrCompressibleContentType(z.contentTypes, contentType) {
				p.HeaderWriter.Header().Set("Content-Encoding", "gzip")
				p.HeaderWriter.Header().Del("Content-Lenght")
				p.Writer = gzip.NewWriter(w)
			}
		},
	}

	z.handler.ServeHTTP(&proxy, r)
	proxy.Close()
}
