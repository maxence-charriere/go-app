// +build !wasm

package app

func render(c Compo) error {
	return ErrNoWasm
}

func run() error {
	return ErrNoWasm
}
