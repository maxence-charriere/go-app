package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
)

func main() {
	app.Import(
		&NavPane{},
		&Hello{},
		&Window{},
	)

	app.Run(&web.Driver{
		URL: "window",
	})
}
