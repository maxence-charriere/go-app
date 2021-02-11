// +build wasm

package app

import (
	"runtime"

	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

const (
	wasmExecJS   = ""
	appJS        = ""
	appWorkerJS  = ""
	manifestJSON = ""
	appCSS       = ""
)

var (
	errBadInstruction = errors.New("unsupported instruction").
		Tag("architecture", runtime.GOARCH)
)

func GenerateStaticWebsite(dir string, h *Handler, pages ...string) error {
	panic(errBadInstruction)
}
