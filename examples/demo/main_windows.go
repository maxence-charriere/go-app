// +build !js

package main

import (
	"fmt"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
	"github.com/murlokswarm/app/drivers/win"
)

func main() {
	app.Import(&Hello{})

	fmt.Println("kind:", app.Kind)

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
