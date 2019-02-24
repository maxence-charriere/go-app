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
		Icon:         "logo.png",
		LoadingLabel: "loading",
		Name:         "app demo",
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
