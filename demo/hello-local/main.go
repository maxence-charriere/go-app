package main

import (
	"fmt"
	"net/http"

	"github.com/maxence-charriere/go-app/pkg/app"
)

func main() {
	fmt.Println("starting local server")

	h := &app.Handler{
		Title:  "Hello Demo",
		Author: "Maxence Charriere",
	}

	if err := http.ListenAndServe(":7000", h); err != nil {
		panic(err)
	}
}
