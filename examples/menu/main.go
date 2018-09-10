package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
)

func main() {
	app.Run(&mac.Driver{
		URL: "/Home",
		MenubarConfig: mac.MenuBarConfig{
			// Overrides the default edit menu.
			EditURL: "/EditMenu",

			// Adds the custom menu in the menubar.
			CustomURLs: []string{"/CustomMenu"},
		},
	})
}
