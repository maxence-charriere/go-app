// +build !wasm

package app

import "net/url"

var (
	window *browserWindow
)

func getenv(k string) string {
	panic(errNoWasm)
}

func keepBodyClean() func() {
	panic(errNoWasm)
}

func navigate(u *url.URL, updateHistory bool) error {
	panic(errNoWasm)
}

func newContextMenu(menuItems ...MenuItemNode) {
	panic(errNoWasm)
}

func reload() {
	panic(errNoWasm)
}

func run() {
	panic(errNoWasm)
}
