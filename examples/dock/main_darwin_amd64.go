package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
)

func main() {
	app.Run(&mac.Driver{
		DockURL:  "/DockMenu",
		OnRun:    func() {},
		OnReopen: func(bool) {},
	})
}
