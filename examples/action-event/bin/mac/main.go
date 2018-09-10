// +build darwin,amd64

package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	actionevent "github.com/murlokswarm/app/examples/action-event"
)

func main() {
	app.Import(&actionevent.Clickbox{})
	app.Import(&actionevent.ClickListener{})

	app.Run(&mac.Driver{
		URL: "/actionevent.Clickbox",
	})
}
