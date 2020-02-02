// +build !wasm

package app

import (
	"fmt"
	"net/url"
	"runtime"
)

var (
	errNoWasm = fmt.Errorf("go architecture is not %q: %q", "wasm", runtime.GOARCH)
	window    *browserWindow
)

func run() {
	panic(errNoWasm)
}

func navigate(u *url.URL, updateHistory bool) error {
	panic(errNoWasm)
}

func reload() {
	panic(errNoWasm)
}

func newContextMenu(menuItems ...MenuItemNode) {
	panic(errNoWasm)
}
