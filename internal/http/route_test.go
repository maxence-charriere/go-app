package http

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouteHandler(t *testing.T) {
	tests := []struct {
		scenario string
		function func(t *testing.T, serv *httptest.Server)
	}{
		{
			scenario: "request is routed to a file",
			function: testRouteHandlerServeFile,
		},
		{
			scenario: "request is routed to a page",
			function: testRouteHandlerServePage,
		},
		{
			scenario: "request to root is routed to a page",
			function: testRouteHandlerServeRootPage,
		},
	}

	serv := httptest.NewServer(&RouteHandler{
		Files:  FileHandler("test"),
		Pages:  &PageHandler{WebDir: "test"},
		WebDir: "test",
	})
	defer serv.Close()

	require.NoError(t, os.Mkdir("test", 0755))
	defer os.RemoveAll("test")

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t, serv)
		})
	}
}

func testRouteHandlerServeFile(t *testing.T, serv *httptest.Server) {
	filename := filepath.Join("test", "hello.txt")
	err := ioutil.WriteFile(filename, []byte("hello world"), 0666)
	require.NoError(t, err)

	res, err := serv.Client().Get(serv.URL + "/hello.txt")
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
}

func testRouteHandlerServePage(t *testing.T, serv *httptest.Server) {
	res, err := serv.Client().Get(serv.URL + "/hello")
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "text/html", res.Header.Get("Content-Type"))
}

func testRouteHandlerServeRootPage(t *testing.T, serv *httptest.Server) {
	res, err := serv.Client().Get(serv.URL)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "text/html", res.Header.Get("Content-Type"))
}
