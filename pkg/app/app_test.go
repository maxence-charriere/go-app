package app

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	Import(&Foo{})

	defer func() { recover() }()
	Import(NoPointerCompo{})
}

func TestRun(t *testing.T) {
	err := Run()

	if runtime.GOARCH != "wasm" {
		require.Error(t, err)
		return
	}

	require.NoError(t, err)
}

func TestUI(t *testing.T) {
	UI(func() {
		t.Log("boo")
	})
}

func TestPath(t *testing.T) {
	assert.Equal(t, "/app.foo", Path(&Foo{}))
}
