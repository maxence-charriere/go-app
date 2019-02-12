// +build !wasm

package app

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	testHandler(t, &Handler{})
}

func TestHandlerWithWebDir(t *testing.T) {
	testHandler(t, &Handler{
		WebDir: func() string { return "." },
	})
}

func testHandler(t *testing.T, h *Handler) {
	tests := []struct {
		scenario string
		function func(t *testing.T, serv *httptest.Server)
	}{
		{
			scenario: "serve a page success",
			function: testHandlerServePage,
		},
		{
			scenario: "serve a file success",
			function: testHandlerServeFile,
		},
	}

	serv := httptest.NewServer(h)

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t, serv)
		})
	}

	serv.Close()
}

func testHandlerServePage(t *testing.T, serv *httptest.Server) {
	res, err := serv.Client().Get(serv.URL)
	require.NoError(t, err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	t.Log("body:", string(body))
}

func testHandlerServeFile(t *testing.T, serv *httptest.Server) {
	defer os.RemoveAll("test.txt")
	ioutil.WriteFile("test.txt", []byte("hello world"), 0666)

	client := serv.Client()
	url := serv.URL + "/test.txt"

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	gzipReader, err := gzip.NewReader(res.Body)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(gzipReader)
	require.NoError(t, err)
	require.Equal(t, "hello world", string(body))
}
