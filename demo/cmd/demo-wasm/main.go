// +build wasm

package main

import (
	"log"

	"github.com/maxence-charriere/app/pkg/app"
)

// The app entry point.
func main() {
	// Imports the hello component declared above in order to make it loadable
	// in a page or usable in other components.
	//
	// Imported component can be used as URL or html tags by referencing them by
	// their lowercased names.
	// E.g:
	//  Hello   => hello
	//  foo.Bar => foo.bar
	app.Import(&Hello{})
	app.Import(&app.ContextMenu{})

	// Defines the component to load when an URL without path is loaded.
	app.DefaultPath = "/hello"

	// Runs the app in the browser.
	if err := app.Run(); err != nil {
		log.Print(err)
	}
}
