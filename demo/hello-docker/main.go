package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/maxence-charriere/go-app/pkg/app"
)

func main() {
	fmt.Println("starting docker server")

	h := &app.Handler{
		Title:   "Hello Demo from Docker",
		Author:  "Maxence Charriere",
		Version: os.Getenv("APP_VERSION"),
	}

	if err := http.ListenAndServe(":7000", h); err != nil {
		panic(err)
	}
}
