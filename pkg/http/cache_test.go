package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManifestServeHTTP(t *testing.T) {
	body := []byte("hello")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	cache := MemoryCache(http.HandlerFunc(handler), 42).(*memoryCache)
	rec := httptest.NewRecorder()
	cache.ServeHTTP(rec, r)

	v, cached := cache.get("/hello")
	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, cached)
	require.Equal(t, "text/plain", v.contentType)
	require.Equal(t, "5", v.contentLength)
	require.Equal(t, body, v.body)
	require.Equal(t, body, rec.Body.Bytes())
}

func TestManifestServeHTTPCachedContent(t *testing.T) {
	body := []byte("hello")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	cache := MemoryCache(http.HandlerFunc(handler), 42).(*memoryCache)
	cache.ServeHTTP(httptest.NewRecorder(), r)
	cache.ServeHTTP(httptest.NewRecorder(), r)

	v, cached := cache.get("/hello")
	require.True(t, cached)
	require.Equal(t, "text/plain", v.contentType)
	require.Equal(t, "5", v.contentLength)
	require.Equal(t, body, v.body)
}

func TestManifestServeHTTPNotCache(t *testing.T) {
	body := []byte("simulate image")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	cache := MemoryCache(http.HandlerFunc(handler), 42).(*memoryCache)
	rec := httptest.NewRecorder()
	cache.ServeHTTP(rec, r)

	_, cached := cache.get("/hello")
	require.False(t, cached)
	require.Equal(t, "image/png", rec.Header().Get("Content-Type"))
	require.Equal(t, fmt.Sprint(len(body)), rec.Header().Get("Content-Length"))
	require.Equal(t, body, rec.Body.Bytes())
}

func TestManifestServeHTTPCacheError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body := strings.TrimPrefix(r.URL.Path, "/")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}

	r1 := httptest.NewRequest("GET", "http://localhost/hello", nil)
	r2 := httptest.NewRequest("GET", "http://localhost/world", nil)

	cache := MemoryCache(http.HandlerFunc(handler), 5).(*memoryCache)
	cache.ServeHTTP(httptest.NewRecorder(), r1)
	cache.ServeHTTP(httptest.NewRecorder(), r2)

	_, cached := cache.get("/hello")
	require.True(t, cached)

	_, cached = cache.get("/world")
	require.False(t, cached)
}

func TestMemoryCacheSet(t *testing.T) {
	cache := memoryCache{capacity: 42}
	err := cache.set("/test", cacheValue{body: []byte("hello")})
	require.NoError(t, err)
	require.Equal(t, 42, cache.capacity, "capacity")
	require.Equal(t, 5, cache.size, "size")
}

func TestMemoryCacheSetInsufficientCapacity(t *testing.T) {
	cache := memoryCache{capacity: 7}

	err := cache.set("/test", cacheValue{body: []byte("hello")})
	require.NoError(t, err)

	err = cache.set("/test", cacheValue{body: []byte("world")})
	require.Error(t, err)
	require.Equal(t, 7, cache.capacity, "capacity")
	require.Equal(t, 5, cache.size, "size")
}

func TestMemoryCacheGet(t *testing.T) {
	cache := memoryCache{capacity: 42}
	err := cache.set("/test", cacheValue{
		contentEncoding: "gzip",
		contentType:     "application/json",
		body:            []byte("hello"),
	})
	require.NoError(t, err)

	v, cached := cache.get("/test")
	require.True(t, cached, "value for /test is not cached")

	require.Equal(t, cacheValue{
		contentEncoding: "gzip",
		contentType:     "application/json",
		body:            []byte("hello"),
	}, v)
}

func TestMemoryGetNotCached(t *testing.T) {
	cache := memoryCache{capacity: 42}
	_, cached := cache.get("/test")
	require.False(t, cached, "value for /test is cached")
}
