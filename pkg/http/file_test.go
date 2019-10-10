package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileServeHTTP(t *testing.T) {
	err := os.MkdirAll("test", 0755)
	require.NoError(t, err)
	defer os.RemoveAll("test")

	content := []byte("hello world")
	err = ioutil.WriteFile(
		filepath.Join("test", "hello.txt"),
		content,
		0666,
	)
	require.NoError(t, err)

	router := Route(Files("test"))
	req := httptest.NewRequest("GET", "http://localhost/hello.txt", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Header().Get("Content-Type"), "text/plain")
	require.Equal(t, content, rec.Body.Bytes())
}

func TestFileServeHTTPDirectory(t *testing.T) {
	err := os.MkdirAll("test/foo", 0755)
	require.NoError(t, err)
	defer os.RemoveAll("test")

	router := Route(Files("test"))
	req := httptest.NewRequest("GET", "http://localhost/foo", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
