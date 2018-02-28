package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
)

func main() {
	// app.Run(&mac.Driver{
	// 	OnRun: func() {
	// 		newWindow()
	// 	},

	// 	OnReopen: func(hasVisibleWindow bool) {
	// 		if !hasVisibleWindow {
	// 			newWindow()
	// 		}
	// 	},
	// })

	app.Run(&web.Driver{
		DefaultURL: "/Hello",
	})
}

// func newWindow() {
// 	app.NewWindow(app.WindowConfig{
// 		Title:           "hello world",
// 		TitlebarHidden:  true,
// 		Width:           1280,
// 		Height:          768,
// 		BackgroundColor: "#21252b",
// 		DefaultURL:      "/Hello",
// 	})
// }
