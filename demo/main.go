// +build !wasm

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/murlokswarm/app"
)

func main() {
	http.Handle("/", &app.Handler{
		Name:   "app demo",
		Wasm:   "demo.wasm",
		WebDir: "web",
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("Defaulting to port %s", port)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
