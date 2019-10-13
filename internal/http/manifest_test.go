package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManifestCanHandle(t *testing.T) {
	man := Manifest{}

	req := httptest.NewRequest(http.MethodGet, "http://localhost/manifest.json", nil)
	require.True(t, man.CanHandle(req))

	req = httptest.NewRequest(http.MethodGet, "http://localhost/manifest", nil)
	require.False(t, man.CanHandle(req))
}

func TestManifestServeHTTP(t *testing.T) {
	handler := Manifest{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost/manifest.json", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	require.Equal(t, lastModified, rec.Header().Get("Last-Modified"))

	t.Log(rec.Body.String())
}
