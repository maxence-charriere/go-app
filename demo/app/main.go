package main

import (
	"github.com/maxence-charriere/app/pkg/app"
)

const (
	icon = "https://storage.googleapis.com/murlok-github/icon-192.png"
)

func main() {
	app.Route("/", &hello{})
	app.Route("/city", &city{})
	app.Run()
}
