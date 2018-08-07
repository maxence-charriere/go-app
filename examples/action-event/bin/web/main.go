package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
	actionevent "github.com/murlokswarm/app/examples/action-event"
)

func main() {
	app.Import(&actionevent.Clickbox{})
	app.Import(&actionevent.ClickListener{})

	app.Run(&web.Driver{
		URL: "/actionevent.Clickbox",
	})
}
