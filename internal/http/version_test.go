package http

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEtag(t *testing.T) {
	require.Empty(t, GetEtag("."))

	etag := GenerateEtag()
	err := ioutil.WriteFile("./.etag", []byte(etag), 0666)
	require.NoError(t, err)
	defer os.Remove(".etag")

	ret := GetEtag(".")
	require.Equal(t, etag, ret)
}
