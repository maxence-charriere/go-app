package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
)

func main() {
	app.Import(&Hello{})

	app.Run(&web.Driver{
		URL: "/Hello",
	})
}
