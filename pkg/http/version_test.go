package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersionServeHTTP(t *testing.T) {
	body := []byte("hello")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	version := Version(http.HandlerFunc(handler), "qwerty")
	rec := httptest.NewRecorder()
	version.ServeHTTP(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
	require.Equal(t, fmt.Sprint(len(body)), rec.Header().Get("Content-Length"))
	require.Equal(t, etagHeaderValue("qwerty"), rec.Header().Get("ETag"))
	require.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
	require.Equal(t, body, rec.Body.Bytes())
}

func TestVersionServeHTTPNoETag(t *testing.T) {
	body := []byte("hello")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	version := Version(http.HandlerFunc(handler), "")
	rec := httptest.NewRecorder()
	version.ServeHTTP(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
	require.Equal(t, fmt.Sprint(len(body)), rec.Header().Get("Content-Length"))
	require.Empty(t, rec.Header().Get("ETag"))
	require.Empty(t, rec.Header().Get("Cache-Control"))
	require.Equal(t, body, rec.Body.Bytes())
}

func TestVersionServeHTTPNotModified(t *testing.T) {
	body := []byte("hello")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	r.Header.Set("If-None-Match", etagHeaderValue("qwerty"))

	version := Version(http.HandlerFunc(handler), "qwerty")
	rec := httptest.NewRecorder()
	version.ServeHTTP(rec, r)

	require.Equal(t, http.StatusNotModified, rec.Code)
	require.Empty(t, rec.Header().Get("Content-Type"))
	require.Empty(t, rec.Header().Get("Content-Length"))
	require.Equal(t, etagHeaderValue("qwerty"), rec.Header().Get("ETag"))
	require.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
	require.Empty(t, rec.Body.Bytes())
}
