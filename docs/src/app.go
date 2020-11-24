// +build wasm

package main

import (
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

func main() {
	for path, new := range pages() {
		app.Route("/"+path, new())
	}

	app.Run()
}
