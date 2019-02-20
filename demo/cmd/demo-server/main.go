// +build !wasm

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/maxence-charriere/app"
)

func main() {
	// Setup the http handler to serve the web assembly app.
	http.Handle("/", &app.Handler{
		Icon:    "logo.png",
		Loading: "loading",
		Name:    "app demo",

		// The path of the directory that contains the wasm app file and the
		// other resources like css files.
		WebDir: "web",

		// The name of the wasm file that contains the app.
		Wasm: "demo.wasm",
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("Defaulting to port %s", port)
	}

	// Launches the server :S.
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
