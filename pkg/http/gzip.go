package http

import (
	"compress/gzip"
	"io"
	"net/http"
	"sync"
)

// Gzip returns a decorated version of the given handler that gzip responses
// bodies with the given content types. It uses DefaultContentTypes when there
// is no content types specified.
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
				contentType := w.Header().Get("Content-Type")
				if isCacheableOrCompressibleContentType(z.contentTypes, contentType) {
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
