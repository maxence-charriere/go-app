package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
)

func mainDesktop() {
	app.Run(&mac.Driver{
		URL: "/Hello",
	})
}
