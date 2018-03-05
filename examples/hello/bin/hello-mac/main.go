package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/examples/hello"
)

func main() {
	app.Import(&hello.Hello{})

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
		Title:           "hello world",
		TitlebarHidden:  true,
		Width:           1280,
		Height:          768,
		BackgroundColor: "#21252b",
		DefaultURL:      "/hello.Hello",
	})
}
