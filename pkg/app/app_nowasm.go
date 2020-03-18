// +build !wasm

package app

import (
	"net/url"
	"runtime"

	"github.com/maxence-charriere/go-app/v6/pkg/log"
)

var (
	window *browserWindow
)

func run() {
	panicNoWasm()
}

func navigate(u *url.URL, updateHistory bool) error {
	panicNoWasm()
	return nil
}

func reload() {
	panicNoWasm()
}

func newContextMenu(menuItems ...MenuItemNode) {
	panicNoWasm()
}

func panicNoWasm() {
	log.Errorf("invalid go architecture").
		T("required", "wasm").
		T("current", runtime.GOARCH).
		Panic()
}

func web(path string) string {
	return path
}
