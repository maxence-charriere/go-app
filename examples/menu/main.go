package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
)

func main() {
	app.Run(&mac.Driver{
		MenubarConfig: mac.MenuBarConfig{
			// Overrides the default edit menu.
			EditURL: "/editmenu",

			// Adds the custom menu in the menubar.
			CustomURLs: []string{"/custommenu"},
		},

		OnRun: func() {
			newWindow()
		},

		OnReopen: func(hasVisibleWindow bool) {
			if !hasVisibleWindow {
				newWindow()
			}
		},
	})
}
func newWindow() {
	app.NewWindow(app.WindowConfig{
		Title:           "menu",
		TitlebarHidden:  true,
		Width:           1280,
		Height:          768,
		BackgroundColor: "#21252b",
		DefaultURL:      "/home",
	})
}
