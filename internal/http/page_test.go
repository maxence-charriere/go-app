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

func TestPageHandler(t *testing.T) {
	serv := httptest.NewServer(&PageHandler{
		WebDir: "test",
	})
	defer serv.Close()

	require.NoError(t, os.Mkdir("test", 0755))
	defer os.RemoveAll("test")

	cssname := filepath.Join("test", "test.css")
	err := ioutil.WriteFile(cssname, []byte(".test{}"), 0666)
	require.NoError(t, err)

	jsname := filepath.Join("test", "test.js")
	err = ioutil.WriteFile(jsname, []byte("alert('hello')"), 0666)
	require.NoError(t, err)

	res, err := serv.Client().Get(serv.URL)
	require.NoError(t, err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	bodyString := string(body)
	require.NoError(t, err)
	assert.Contains(t, bodyString, `<link type="text/css" rel="stylesheet" href="/test.css">`)
	assert.Contains(t, bodyString, `<script src="/test.js"></script>`)
	t.Logf("page: %s", body)
}
