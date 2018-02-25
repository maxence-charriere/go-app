package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
)

func main() {
	app.Run(&mac.Driver{
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
		Title:           "Drag and Drop",
		TitlebarHidden:  true,
		MinWidth:        1000,
		Width:           1280,
		MinHeight:       550,
		Height:          768,
		BackgroundColor: "#21252b",
		DefaultURL:      "/DragDrop",
	})
}
