//go:build wasm
// +build wasm

package app

import (
	"runtime"

	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

func GenerateStaticWebsite(dir string, h *Handler, pages ...string) error {
	panic(errors.New("unsupported instruction").
		WithTag("architecture", runtime.GOARCH))
}

func wasmExecJS() string {
	return ""
}
