// +build !wasm

package app

func navigate(url string) {
	Log("navigating to", url)
}

func reload() {
	Log("reloading")
}

func render(c Compo) error {
	return ErrNoWasm
}

func run() error {
	return ErrNoWasm
}
