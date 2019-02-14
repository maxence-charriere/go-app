package main

import (
	"log"

	"github.com/murlokswarm/app"
)

func main() {
	app.Import(&Hello{})

	app.DefaultPath = "/hello"

	if err := app.Run(); err != nil {
		log.Print(err)
	}
}
