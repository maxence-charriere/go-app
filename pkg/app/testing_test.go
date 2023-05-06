package app

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/maxence-charriere/go-app/v9/pkg/logs"
	"github.com/stretchr/testify/require"
)

func testSkipNonJS(t *testing.T) {
	if goos := runtime.GOOS; goos != "js" {
		t.Skip(logs.New("skipping test").
			WithTag("reason", "unsupported OS").
			WithTag("required-os", "js").
			WithTag("current-os", goos),
		)
	}
}

func testSkipJS(t *testing.T) {
	if goos := runtime.GOOS; goos == "js" {
		t.Skip(logs.New("skipping test").
			WithTag("reason", "unsupported OS").
			WithTag("required-os", "!= than js").
			WithTag("current-os", goos),
		)
	}
}

func testCreateDir(t *testing.T, path string) func() {
	err := os.MkdirAll(path, 0755)
	require.NoError(t, err)

	return func() {
		os.RemoveAll(path)
	}
}

func testCreateFile(t *testing.T, path, content string) {
	err := ioutil.WriteFile(path, []byte(content), 0666)
	require.NoError(t, err)
}
