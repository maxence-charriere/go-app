// +build !js

package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/drivers/web"
)

func main() {
	app.Import(&Hello{})

	switch app.Kind {
	case "web":
		app.Run(&web.Driver{
			URL: "/Hello",
		})

	default:
		app.Run(&mac.Driver{
			URL: "/Hello",
		})
	}
}
