//go:build wasm
// +build wasm

package app

import (
	"net/http"
	"runtime"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
)

func GenerateStaticWebsite(dir string, h *Handler, pages ...string) error {
	panic(errors.New("unsupported instruction").
		WithTag("architecture", runtime.GOARCH))
}

func GenerateStaticWebsiteFromMux(dir string, h http.Handler, pages ...string) error {
	panic(errors.New("unsupported instruction").
		WithTag("architecture", runtime.GOARCH))
}

func wasmExecJS() string {
	return ""
}
