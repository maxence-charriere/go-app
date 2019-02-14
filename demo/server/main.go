// +build !wasm

package main

import (
	"net/http"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
)

func main() {
	http.Handle("/", &app.Handler{
		Author:      "Maxence Charriere",
		Description: "A demo that shows what the app package can build.",
		Keywords: []string{
			"go",
			"golang",
			"app",
			"ui",
			"gui",
			"wasm",
			"web assembly",
		},
		Name: "app demo",
	})

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Error(err)
	}
}
