// +build !js

package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/win"
)

func main() {
	app.Import(&Hello{})

	app.Run(&win.Driver{
		URL: "/Hello",
	})
}
