package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/win"
)

func mainDesktop() {
	app.Run(&win.Driver{
		URL: "/Hello",
	})
}
