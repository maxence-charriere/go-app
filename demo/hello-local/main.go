package main

import (
	"fmt"
	"net/http"

	"github.com/maxence-charriere/app/pkg/app"
)

func main() {
	fmt.Println("starting local server")

	h := &app.Handler{
		Title:  "Hello Demo",
		Author: "Maxence Charriere",
		Wasm:   "app.wasm",
	}

	if err := http.ListenAndServe(":7777", h); err != nil {
		panic(err)
	}
}
