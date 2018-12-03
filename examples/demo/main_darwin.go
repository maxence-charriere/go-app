// +build !js

package main

import (
	"net/url"

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
			Settings: mac.Settings{
				URLScheme: "goapp-demo",
			},

			OnURLOpen: func(u *url.URL) {
				app.Log("app opened with:", u)
			},
		})
	}
}
