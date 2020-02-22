package main

import (
	"fmt"
	"net/http"

	"github.com/maxence-charriere/app/pkg/app"
)

func main() {
	fmt.Println("starting docker server")

	h := &app.Handler{
		Title:  "Hello Demo from Docker",
		Author: "Maxence Charriere",
	}

	if err := http.ListenAndServe(":7000", h); err != nil {
		panic(err)
	}
}
