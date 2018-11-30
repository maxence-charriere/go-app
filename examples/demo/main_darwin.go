// +build !js

package main

import (
	"fmt"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/drivers/web"
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
		app.Run(&mac.Driver{
			URL: "/Hello",
		})
	}
}
