package main

import "github.com/murlokswarm/app"

var (
	win app.Windower
)

func main() {
	app.OnLaunch = func() {
		if menuBar, ok := app.MenuBar(); ok {
			menuBar.Mount(&MenuBar{})
		}

		win = newMainWindow()
		win.Mount(&Paris{})
	}

	app.OnReopen = func() {
		if win != nil {
			return
		}
		win = newMainWindow()
		win.Mount(&Paris{})
	}

	app.Run()
}

func newMainWindow() app.Windower {
	return app.NewWindow(app.Window{
		Title:           "nav",
		TitlebarHidden:  true,
		Width:           1280,
		Height:          768,
		BackgroundColor: "#21252b",
		OnClose: func() bool {
			win = nil
			return true
		},
	})
}
