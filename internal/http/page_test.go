package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageCanHandle(t *testing.T) {
	p := Page{WebDir: "test"}
	require.True(t, p.CanHandle(nil))
}

func TestPageServeHTTP(t *testing.T) {
	require.NoError(t, os.Mkdir("test", 0755))
	defer os.RemoveAll("test")

	cssname := filepath.Join("test", "test.css")
	err := ioutil.WriteFile(cssname, []byte(".test{}"), 0666)
	require.NoError(t, err)

	jsname := filepath.Join("test", "test.js")
	err = ioutil.WriteFile(jsname, []byte("alert('hello')"), 0666)
	require.NoError(t, err)

	handler := Page{WebDir: "test"}
	req := httptest.NewRequest(http.MethodGet, "http://localhost/index.html", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "text/html", rec.Header().Get("Content-Type"))
	require.Equal(t, lastModified, rec.Header().Get("Last-Modified"))
	assert.Contains(t, body, `<link type="text/css" rel="stylesheet" href="/test.css">`)
	assert.Contains(t, body, `<script src="/test.js"></script>`)
	t.Log(body)
}
