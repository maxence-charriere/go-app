package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
)

func main() {
	app.Import(&Hello{})

	switch app.Kind {
	case "web":
		mainWeb()

	default:
		mainDesktop()
	}
}

func mainWeb() {
	app.Run(&web.Driver{
		URL: "/Hello",
	})
}
