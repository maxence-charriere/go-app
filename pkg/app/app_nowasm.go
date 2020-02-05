// +build !wasm

package app

import (
	"net/url"
	"runtime"

	"github.com/maxence-charriere/app/pkg/log"
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
	log.Errorf("bad go architecture").
		T("required", "wasm").
		T("current", runtime.GOARCH).
		Panic()
}
