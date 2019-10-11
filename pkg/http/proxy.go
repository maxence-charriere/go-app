package http

import (
	"io"
	"net/http"
	"sync"
)

// ProxyWriter is a wrapper for http.ResponseWriter that provide way to override
// Write, get status code or close a writer.
type ProxyWriter struct {
	// The writer used to write data.
	Writer io.Writer

	// The writer used to write and access headers.
	HeaderWriter http.ResponseWriter

	// A function that when set, is called once before writes operations occurs.
	BeforeWrite func(*ProxyWriter)

	once       sync.Once
	statusCode int
}

// Header returns the Header.
func (w *ProxyWriter) Header() http.Header {
	return w.HeaderWriter.Header()
}

func (w *ProxyWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// WriteHeader writes the response status code.
func (w *ProxyWriter) WriteHeader(statusCode int) {
	w.once.Do(func() {
		if w.BeforeWrite != nil {
			w.BeforeWrite(w)
		}
	})

	w.statusCode = statusCode
	w.HeaderWriter.WriteHeader(statusCode)
}

// Close closes the writer when it is a io.Closer.
func (w *ProxyWriter) Close() error {
	if closer, ok := w.Writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// StatusCode returns the response status code.
func (w *ProxyWriter) StatusCode() int {
	return w.statusCode
}
