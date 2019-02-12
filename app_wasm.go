package app

import "syscall/js"

func run() error {
	url := js.Global().
		Get("location").
		Get("href").
		String()

	return nil
}
