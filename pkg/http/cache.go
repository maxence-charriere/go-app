package http

import (
	"bytes"
	"errors"
	"net/http"
	"sync"

	"github.com/maxence-charriere/app/pkg/log"
)

// MemoryCache returns a decorated version of the given http.Handler that caches
// and serves responses bodies with the given contents type. It uses
// DefaultContentTypes when there is no content types specified.
func MemoryCache(h http.Handler, capacity int, contentTypes ...string) http.Handler {
	if len(contentTypes) == 0 {
		contentTypes = DefaultContentTypes
	}

	return &memoryCache{
		handler:      h,
		capacity:     capacity,
		contentTypes: contentTypes,
	}
}

type memoryCache struct {
	handler      http.Handler
	capacity     int
	contentTypes []string

	once   sync.Once
	mu     sync.RWMutex
	size   int
	values map[string]cacheValue
}

func (c *memoryCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	c.mu.RLock()
	if v, cached := c.get(path); cached {
		w.Header().Set("Content-Encoding", v.contentEncoding)
		w.Header().Set("Content-Type", v.contentType)
		w.Header().Set("Content-Length", v.contentLength)

		if n, err := w.Write(v.body); err != nil {
			log.Error("writing cached data failed").
				T("error", err).
				T("path", path).
				T("content-encoding", v.contentEncoding).
				T("content-type", v.contentType).
				T("content-length", v.contentLength).
				T("bytes written", n)
		}

		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()

	var once sync.Once
	var isCacheable bool
	var buffer bytes.Buffer

	proxy := proxyWriter{
		header: w.Header,
		write: func(b []byte) (int, error) {
			once.Do(func() {
				contentType := w.Header().Get("Content-Type")
				isCacheable = isCacheableOrCompressibleContentType(c.contentTypes, contentType)
			})

			if isCacheable {
				if n, err := buffer.Write(b); err != nil {
					return n, err
				}
			}
			return w.Write(b)

		},
		writeHeader: w.WriteHeader,
	}
	c.handler.ServeHTTP(proxy, r)

	if !isCacheable {
		return
	}

	c.mu.Lock()
	if _, cached := c.get(path); cached {
		return
	}

	v := cacheValue{
		contentEncoding: w.Header().Get("Content-Encoding"),
		contentType:     w.Header().Get("Content-Type"),
		contentLength:   w.Header().Get("Content-Length"),
		body:            buffer.Bytes(),
	}

	if err := c.set(path, v); err != nil {
		log.Error("caching response body failed").
			T("error", err).
			T("path", path).
			T("cache capacity", c.capacity).
			T("cache size", c.size).
			T("body length", len(v.body))
	}
	c.mu.Unlock()
}

func (c *memoryCache) set(path string, v cacheValue) error {
	if c.size+len(v.body) > c.capacity {
		return errors.New("insufficient capacity")
	}

	if c.values == nil {
		c.values = make(map[string]cacheValue)
	}
	c.values[path] = v
	c.size += len(v.body)

	return nil
}

func (c *memoryCache) get(path string) (cacheValue, bool) {
	v, cached := c.values[path]
	return v, cached
}

type cacheValue struct {
	contentEncoding string
	contentType     string
	contentLength   string
	body            []byte
}
