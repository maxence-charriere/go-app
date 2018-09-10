// +build darwin,amd64

package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/examples/nav"
)

func main() {
	app.Import(&nav.City{})

	app.Run(&mac.Driver{
		URL: "/nav.City",
	})
}
