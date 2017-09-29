package main

import (
	"log"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
)

func main() {
	app.Run(&mac.Driver{
		OnRun: func() {
			log.Println("OnRun")
			log.Println("app.Resources():", app.Resources())
			log.Println("app.Storage():", app.Storage())

			win := app.NewWindow(app.WindowConfig{
				Title:  "test window",
				X:      42,
				Y:      42,
				Width:  1024,
				Height: 600,

				OnMove: func(x, y float64) {
					log.Printf("Window moved to x:%v y:%v", x, y)
				},
				OnResize: func(width, height float64) {
					log.Printf("Window resized to width:%v height:%v", width, height)
				},
				OnFocus: func() {
					log.Println("Window focused")
				},
				OnBlur: func() {
					log.Println("Window blured")
				},
			})

			x, y := win.Position()
			log.Printf("win.Positon() x:%v, x:%v", x, y)

			log.Printf("win.Move(x:%v, y: %v)", 42, 42)
			win.Move(42, 42)

			log.Println("win.Center()")
			win.Center()

			width, height := win.Size()
			log.Printf("win.Size() width:%v, height:%v", width, height)

			log.Printf("win.Resize(x:%v, y: %v)", 1340, 720)
			win.Resize(1340, 720)

			win.Focus()

			log.Println("all tests OK")
		},
		OnFocus: func() {
			log.Println("OnFocus")
		},
		OnBlur: func() {
			log.Println("OnBlur")
		},
		OnReopen: func(hasVisibleWindows bool) {
			log.Println("OnReopen hasVisibleWIndow:", hasVisibleWindows)
		},
		OnQuit: func() bool {
			log.Println("OnQuit")
			return true
		},
		OnExit: func() {
			log.Println("OnExit")
		},
	})
}
