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

func TestGzipHandler(t *testing.T) {
	tests := []struct {
		scenario string
		function func(t *testing.T, serv *httptest.Server)
	}{
		{
			scenario: "request without Accept-Encoding serves a non gzipped file",
			function: testGzipHandlerServeWithoutAcceptEncoding,
		},
		{
			scenario: "request serves a non gzipped file",
			function: testGzipHandlerServeNonGzippedFile,
		},
		{
			scenario: "request serves a gzipped file",
			function: testGzipHandlerServeGzippedFile,
		},
	}

	handler := FileHandler("test")
	handler = GzipHandler(handler, "test")
	serv := httptest.NewServer(handler)
	defer serv.Close()

	require.NoError(t, os.Mkdir("test", 0755))
	defer os.RemoveAll("test")

	filename := filepath.Join("test", "hello.txt")
	err := ioutil.WriteFile(filename, []byte("hello world"), 0666)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t, serv)
		})
	}
}

func testGzipHandlerServeWithoutAcceptEncoding(t *testing.T, serv *httptest.Server) {
	filename := filepath.Join("test", "hello.txt")
	err := ioutil.WriteFile(filename, []byte("hello world"), 0666)
	require.NoError(t, err)
	defer os.Remove(filename)

	req, err := http.NewRequest(http.MethodGet, serv.URL+"/hello.txt", nil)
	require.NoError(t, err)

	res, err := serv.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Empty(t, res.Header.Get("Content-Encoding"))
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
}

func testGzipHandlerServeNonGzippedFile(t *testing.T, serv *httptest.Server) {
	req, err := http.NewRequest(http.MethodGet, serv.URL+"/hello.txt", nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")

	res, err := serv.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Empty(t, res.Header.Get("Content-Encoding"))
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
}

func testGzipHandlerServeGzippedFile(t *testing.T, serv *httptest.Server) {
	gzipname := filepath.Join("test", "hello.txt.gz")
	err := ioutil.WriteFile(gzipname, []byte("qsdcvfbnmj"), 0666)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, serv.URL+"/hello.txt", nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")

	res, err := serv.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
}
