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

func TestCacheHandler(t *testing.T) {
	tests := []struct {
		scenario string
		function func(t *testing.T)
	}{
		{
			scenario: "request to etag file returns a 400",
			function: testCacheHandlerRequestEtagFile,
		},
		{
			scenario: "request with no etag set returns a response without caching",
			function: testCacheHandlerRequestNoEtag,
		},
		{
			scenario: "request with etag returns 304",
			function: testCacheHandlerRequestWithEtagMatch,
		},
		{
			scenario: "request with etag returns a 200",
			function: testCacheHandlerRequestWithEtagNoMatch,
		},
	}

	// handler := FileHandler("test")
	// handler = CacheHandler(handler, "test")
	// serv := httptest.NewServer(handler)
	// defer serv.Close()

	// require.NoError(t, os.Mkdir("test", 0755))
	// defer os.RemoveAll("test")

	// etagname := filepath.Join("test", ".etag")
	// err := ioutil.WriteFile(etagname, []byte(GenerateEtag()), 0666)
	// require.NoError(t, err)

	// filename := filepath.Join("test", "hello.txt")
	// err = ioutil.WriteFile(filename, []byte("hello world"), 0666)
	// require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.scenario, test.function)
	}
}

func TestGenerateEtag(t *testing.T) {
	t.Log(GenerateEtag())
}

func testCacheHandlerRequestEtagFile(t *testing.T) {
	handler := FileHandler("test")
	handler = CacheHandler(handler, "test")
	serv := httptest.NewServer(handler)
	defer serv.Close()

	res, err := serv.Client().Get(serv.URL + "/.etag")
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 400, res.StatusCode)
}

func testCacheHandlerRequestNoEtag(t *testing.T) {
	require.NoError(t, os.Mkdir("test", 0755))
	defer os.RemoveAll("test")

	filename := filepath.Join("test", "hello.txt")
	err := ioutil.WriteFile(filename, []byte("hello world"), 0666)
	require.NoError(t, err)

	handler := FileHandler("test")
	handler = CacheHandler(handler, "test")
	serv := httptest.NewServer(handler)
	defer serv.Close()

	req, err := http.NewRequest(http.MethodGet, serv.URL+"/hello.txt", nil)
	require.NoError(t, err)
	req.Header.Set("If-None-Match", GenerateEtag())

	res, err := serv.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Empty(t, res.Header.Get("ETag"))
	assert.Empty(t, res.Header.Get("Cache-Control"))
}

func testCacheHandlerRequestWithEtagMatch(t *testing.T) {
	require.NoError(t, os.Mkdir("test", 0755))
	defer os.RemoveAll("test")

	etag := GenerateEtag()
	etagname := filepath.Join("test", ".etag")
	err := ioutil.WriteFile(etagname, []byte(etag), 0666)
	require.NoError(t, err)

	handler := FileHandler("test")
	handler = CacheHandler(handler, "test")
	serv := httptest.NewServer(handler)
	defer serv.Close()

	req, err := http.NewRequest(http.MethodGet, serv.URL+"/hello.txt", nil)
	require.NoError(t, err)
	req.Header.Set("If-None-Match", etag)

	res, err := serv.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotModified, res.StatusCode)
	assert.Equal(t, etag, res.Header.Get("ETag"))
}

func testCacheHandlerRequestWithEtagNoMatch(t *testing.T) {
	require.NoError(t, os.Mkdir("test", 0755))
	defer os.RemoveAll("test")

	etag := GenerateEtag()
	etagname := filepath.Join("test", ".etag")
	err := ioutil.WriteFile(etagname, []byte(etag), 0666)
	require.NoError(t, err)

	filename := filepath.Join("test", "hello.txt")
	err = ioutil.WriteFile(filename, []byte("hello world"), 0666)
	require.NoError(t, err)

	handler := FileHandler("test")
	handler = CacheHandler(handler, "test")
	serv := httptest.NewServer(handler)
	defer serv.Close()

	req, err := http.NewRequest(http.MethodGet, serv.URL+"/hello.txt", nil)
	require.NoError(t, err)
	req.Header.Set("If-None-Match", GenerateEtag())

	res, err := serv.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, etag, res.Header.Get("ETag"))
}
