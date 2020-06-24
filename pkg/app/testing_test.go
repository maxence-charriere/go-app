package app

import (
	"runtime"
	"testing"

	"github.com/maxence-charriere/go-app/v7/pkg/logs"
)

func testSkipNonWasm(t *testing.T) {
	if goarch := runtime.GOARCH; goarch != "wasm" {
		t.Skip(logs.New("skipping test").
			Tag("reason", "unsupported architecture").
			Tag("required-architecture", "wasm").
			Tag("current-architecture", goarch),
		)
	}
}

func testSkipWasm(t *testing.T) {
	if goarch := runtime.GOARCH; goarch == "wasm" {
		t.Skip(logs.New("skipping test").
			Tag("reason", "unsupported architecture").
			Tag("required-architecture", "!= than wasm").
			Tag("current-architecture", goarch),
		)
	}
}
