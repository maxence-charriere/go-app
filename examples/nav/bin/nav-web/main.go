package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
	"github.com/murlokswarm/app/examples/nav"
)

func main() {
	app.Import(&nav.City{})

	app.Run(&web.Driver{
		URL: "/nav.City",
	})
}
