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

func TestGzipServeHttp(t *testing.T) {
	body := []byte("hello gzip world")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	gz := Gzip(http.HandlerFunc(handler))
	rec := httptest.NewRecorder()
	gz.ServeHTTP(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))
	require.Equal(t, "text/plain", rec.Header().Get("Content-Type"))

	rzip, err := gzip.NewReader(rec.Body)
	require.NoError(t, err)

	res, err := ioutil.ReadAll(rzip)
	require.NoError(t, err)
	require.Equal(t, body, res)
}

func TestGzipServeHttpNoCompressible(t *testing.T) {
	body := []byte("simulated image")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	gz := Gzip(http.HandlerFunc(handler))
	rec := httptest.NewRecorder()
	gz.ServeHTTP(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "image/png", rec.Header().Get("Content-Type"))
	require.Empty(t, rec.Header().Get("Content-Encoding"))
	require.Equal(t, body, rec.Body.Bytes())
}

func TestGzipServeHttpNoWrite(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	gz := Gzip(http.HandlerFunc(handler))
	rec := httptest.NewRecorder()
	gz.ServeHTTP(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Header().Get("Content-Encoding"))
	require.Empty(t, rec.Header().Get("Content-Type"))
	require.Empty(t, rec.Body.Len())
}

func TestGzipServeHttpNoContentType(t *testing.T) {
	body := []byte("simulated image")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	gz := Gzip(http.HandlerFunc(handler))
	rec := httptest.NewRecorder()
	gz.ServeHTTP(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Header().Get("Content-Encoding"))
	require.Empty(t, rec.Header().Get("Content-Type"))
	require.Equal(t, body, rec.Body.Bytes())
}
