package main

import (
	"fmt"
	"net/http"

	"github.com/maxence-charriere/app/pkg/app"
)

func main() {
	fmt.Println("starting local server")

	if err := http.ListenAndServe(":7777", &app.Handler{}); err != nil {
		panic(err)
	}
}
