// +build !wasm

package main

import (
	"net/http"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
)

func main() {
	http.Handle("/", &app.Handler{
		Name:   "app demo",
		Wasm:   "demo.wasm",
		WebDir: "web",
	})

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Error(err)
	}
}
