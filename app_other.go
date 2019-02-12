// +build !wasm

package app

import "github.com/pkg/errors"

func run() error {
	return errors.New("go architecture is not wasm")
}
