package http

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGzipServeHTTP(t *testing.T) {
	body := []byte("hello gzip world")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	gz := Gzip(http.HandlerFunc(handler))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost/hello", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	gz.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))
	require.Equal(t, "text/plain", rec.Header().Get("Content-Type"))

	rzip, err := gzip.NewReader(rec.Body)
	require.NoError(t, err)

	res, err := ioutil.ReadAll(rzip)
	require.NoError(t, err)
	require.Equal(t, body, res)
}

func TestGzipServeHTTPNoAcceptEncoding(t *testing.T) {
	body := []byte("simulated image")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	gz := Gzip(http.HandlerFunc(handler))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost/hello", nil)
	gz.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "image/png", rec.Header().Get("Content-Type"))
	require.Empty(t, rec.Header().Get("Content-Encoding"))
	require.Equal(t, fmt.Sprint(len(body)), rec.Header().Get("Content-Length"))
	require.Equal(t, body, rec.Body.Bytes())
}

func TestGzipServeHTTPNotCompressible(t *testing.T) {
	body := []byte("simulated image")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	gz := Gzip(http.HandlerFunc(handler))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost/hello", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	gz.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "image/png", rec.Header().Get("Content-Type"))
	require.Empty(t, rec.Header().Get("Content-Encoding"))
	require.Equal(t, fmt.Sprint(len(body)), rec.Header().Get("Content-Length"))
	require.Equal(t, body, rec.Body.Bytes())
}
