// +build !js

package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
	"github.com/murlokswarm/app/drivers/win"
)

func main() {
	app.Import(&Hello{})

	switch app.Kind {
	case "web":
		app.Run(&web.Driver{
			URL: "/Hello",
		})

	default:
		app.Run(&win.Driver{
			URL: "/Hello",
		})
	}
}
