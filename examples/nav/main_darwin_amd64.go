package main

import (
	"fmt"
	"log"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
)

func main() {
	app.Run(&mac.Driver{
		OnRun: func() {
			log.Println("OnRun")
			fmt.Println("app.Resources():", app.Resources())
			fmt.Println("app.Storage():", app.Storage())

			win := app.NewWindow(app.WindowConfig{
				Title:  "test window",
				X:      42,
				Y:      42,
				Width:  1024,
				Height: 600,

				OnMove: func(x, y float64) {
					fmt.Printf("window moved to x:%v y:%v\n", x, y)
				},
			})

			x, y := win.Position()
			fmt.Printf("win.Positon() x:%v, x:%v\n", x, y)

			fmt.Printf("win.Move(x:%v, y: %v)\n", 42, 42)
			win.Move(42, 42)

			fmt.Println("win.Center()")
			win.Center()

			fmt.Println("all tests OK")
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
