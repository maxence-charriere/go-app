// +build darwin,amd64

package main

import (
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/examples/test"
)

func main() {
	app.Import(&test.Webview{})
	app.Import(&test.Menu{})

	app.Run(&mac.Driver{
		Bundle: mac.Bundle{
			Icon:        "logo.png",
			FilePickers: mac.FileReadWriteAccess,
		},

		MenubarConfig: mac.MenuBarConfig{
			OnPreference: func() {
				app.Logf("Preferences clicked")
			},
		},
		DockURL: "/test.Menu",

		OnRun: func() {
			app.Log("OnRun")
			app.Resources()
			app.Storage()

			// testWindow(true)
			testWindow(false)

			app.NewStatusMenu(app.StatusMenuConfig{
				Text: "only text",
				URL:  "/test.Menu",
			})

			app.NewStatusMenu(app.StatusMenuConfig{
				Icon: app.Resources("logo.png"),
				URL:  "/test.Menu",
			})

			statMenu := app.NewStatusMenu(app.StatusMenuConfig{
				Text: "text + ",
				Icon: app.Resources("logo.png"),
				URL:  "/test.Menu",
			})

			go func() {
				time.Sleep(time.Second)
				statMenu.SetText("")

				time.Sleep(time.Second)
				statMenu.SetText("Hello")

				time.Sleep(time.Second)
				statMenu.SetIcon("")

				time.Sleep(time.Second)
				statMenu.SetIcon(app.Resources("logo.png"))

				time.Sleep(time.Second)
				statMenu.SetText("")
				statMenu.SetIcon("")

				time.Sleep(time.Second)
				statMenu.SetIcon(app.Resources("logo.png"))
				statMenu.SetText("Bye bye")

				time.Sleep(time.Second)
				statMenu.Close()
			}()
		},
		OnFocus: func() {
			app.Log("OnFocus")
		},
		OnBlur: func() {
			app.Log("OnBlur")
		},
		OnReopen: func(hasVisibleWindows bool) {
			app.Log("OnReopen hasVisibleWIndow:", hasVisibleWindows)
			if hasVisibleWindows {
				return
			}
			testWindow(false)
		},
		OnQuit: func() bool {
			app.Log("OnQuit")
			return true
			// return false
		},
		OnExit: func() {
			app.Log("OnExit")
		},
	})
}

func testWindow(close bool) {
	win := app.NewWindow(app.WindowConfig{
		Title:    "test window",
		X:        42,
		Y:        42,
		Width:    1024,
		MinWidth: 400,
		// MaxWidth:  300,
		Height:    600,
		MinHeight: 400,
		// FixedSize: true,
		// CloseHidden:    true,
		// MinimizeHidden: true,
		TitlebarHidden:    true,
		FrostedBackground: true,
		URL:               "/test.Webview",

		OnMove: func(x, y float64) {
			app.Logf("Window moved to x:%v y:%v\n", x, y)
		},
		OnResize: func(width, height float64) {
			app.Logf("Window resized to width:%v height:%v\n", width, height)
		},
		OnFocus: func() {
			app.Log("Window focused")
		},
		OnBlur: func() {
			app.Log("Window blured")
		},
		OnFullScreen: func() {
			app.Log("Window full screen")
		},
		OnExitFullScreen: func() {
			app.Log("Window exit full screen")
		},
		OnMinimize: func() {
			app.Log("Window minimized")
		},
		OnDeminimize: func() {
			app.Log("Window deminimized")
		},
		OnClose: func() bool {
			app.Log("Window close")
			return true
			// return false
		},
	})

	x, y := win.Position()
	app.Logf("win.Positon() x:%v, x:%v\n", x, y)

	app.Logf("win.Move(x:%v, y: %v)\n", 42, 42)
	win.Move(42, 42)

	app.Log("win.Center()")
	win.Center()

	width, height := win.Size()
	app.Logf("win.Size() width:%v, height:%v\n", width, height)

	// app.Logf("win.Resize(x:%v, y: %v)\n", 1340, 720)
	// win.Resize(1340, 720)

	// app.Log("win.ToggleFullScreen()")
	// win.ToggleFullScreen()

	// go func() {
	// 	app.Log("win.ToggleFullScreen()")
	// 	time.Sleep(2 * time.Second)
	// 	win.ToggleFullScreen()
	// }()

	// app.Log("win.ToggleMinimize()")
	// win.ToggleMinimize()

	// app.Log("win.ToggleMinimize()")
	// win.ToggleMinimize()

	// win.Focus()

	if close {
		app.Log("win.Close()")
		win.Close()
	}

	app.Log("Window tests OK")
}
