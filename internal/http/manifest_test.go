package http

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestHandler(t *testing.T) {
	serv := httptest.NewServer(&ManifestHandler{})
	defer serv.Close()

	res, err := serv.Client().Get(serv.URL)
	require.NoError(t, err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.NotEmpty(t, body)
	t.Logf("manifest: %s", body)
}
