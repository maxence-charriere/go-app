package app

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/maxence-charriere/go-app/v9/pkg/logs"
	"github.com/stretchr/testify/require"
)

func testSkipNonWasm(t *testing.T) {
	if goarch := runtime.GOARCH; goarch != "wasm" {
		t.Skip(logs.New("skipping test").
			WithTag("reason", "unsupported architecture").
			WithTag("required-architecture", "wasm").
			WithTag("current-architecture", goarch),
		)
	}
}

func testSkipWasm(t *testing.T) {
	if goarch := runtime.GOARCH; goarch == "wasm" {
		t.Skip(logs.New("skipping test").
			WithTag("reason", "unsupported architecture").
			WithTag("required-architecture", "!= than wasm").
			WithTag("current-architecture", goarch),
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
