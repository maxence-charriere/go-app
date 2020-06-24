// +build !wasm

package app

import (
	"net/url"
	"os"
)

var (
	window *browserWindow
)

func getenv(k string) string {
	return os.Getenv(k)
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
