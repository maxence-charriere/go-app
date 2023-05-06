package app

import (
	"runtime"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

const (
	wasmExecJS   = ""
	wasmDriverJS = ""
	appJS        = ""
	appWorkerJS  = ""
	manifestJSON = ""
	appCSS       = ""
)

var (
	errBadInstruction = errors.New("unsupported instruction").
		WithTag("os", runtime.GOOS)
)

func GenerateStaticWebsite(dir string, h *Handler, pages ...string) error {
	panic(errBadInstruction)
}
