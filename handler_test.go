// +build !wasm

package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("body:", string(body))
}

func testHandlerServeFile(t *testing.T, serv *httptest.Server) {
	defer os.RemoveAll("test.txt")
	ioutil.WriteFile("test.txt", []byte("hello"), 0666)

	client := serv.Client()
	url := serv.URL + "/test.txt"

	res, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if s := string(body); s != "hello" {
		t.Error("expected body:", "hello")
		t.Fatal("returned body:", s)
	}
}

func testHandlerServeGzipFile(t *testing.T, serv *httptest.Server) {
	defer os.RemoveAll("test.txt")
	ioutil.WriteFile("test.txt", []byte("hello"), 0666)

	client := serv.Client()
	url := serv.URL + "/test.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if s := string(body); s != "hello" {
		t.Error("expected body:", "hello")
		t.Fatal("returned body:", s)
	}
}
