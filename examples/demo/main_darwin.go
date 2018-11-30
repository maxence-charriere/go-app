package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
)

func main() {
	app.Import(&Hello{})

	app.Run(&mac.Driver{
		URL: "/Hello",
	})
}
