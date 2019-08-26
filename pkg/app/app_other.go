// +build !wasm

package app

func navigate(url string) {
}

func reload() {
}

func render(c Compo) error {
	return ErrNoWasm
}

func run() error {
	return ErrNoWasm
}
