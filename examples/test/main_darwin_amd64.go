package main

import (
	"fmt"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/log"
)

func main() {
	app.Run(&mac.Driver{
		Logger: &log.Logger{Debug: true},

		OnRun: func() {
			fmt.Println("OnRun")
			fmt.Println("app.Resources():", app.Resources())
			fmt.Println("app.Storage():", app.Storage())

			testWindow(true)
			testWindow(false)
		},
		OnFocus: func() {
			fmt.Println("OnFocus")
		},
		OnBlur: func() {
			fmt.Println("OnBlur")
		},
		OnReopen: func(hasVisibleWindows bool) {
			fmt.Println("OnReopen hasVisibleWIndow:", hasVisibleWindows)
			if hasVisibleWindows {
				return
			}
			testWindow(false)
		},
		OnQuit: func() bool {
			fmt.Println("OnQuit")
			return true
		},
		OnExit: func() {
			fmt.Println("OnExit")
		},
	})
}

func testWindow(close bool) {
	win := app.NewWindow(app.WindowConfig{
		Title:  "test window",
		X:      42,
		Y:      42,
		Width:  1024,
		Height: 600,

		OnMove: func(x, y float64) {
			fmt.Printf("Window moved to x:%v y:%v\n", x, y)
		},
		OnResize: func(width, height float64) {
			fmt.Printf("Window resized to width:%v height:%v\n", width, height)
		},
		OnFocus: func() {
			fmt.Println("Window focused")
		},
		OnBlur: func() {
			fmt.Println("Window blured")
		},
		OnClose: func() bool {
			fmt.Println("Window close")
			return true
		},
	})

	x, y := win.Position()
	fmt.Printf("win.Positon() x:%v, x:%v\n", x, y)

	fmt.Printf("win.Move(x:%v, y: %v)\n", 42, 42)
	win.Move(42, 42)

	fmt.Println("win.Center()")
	win.Center()

	width, height := win.Size()
	fmt.Printf("win.Size() width:%v, height:%v\n", width, height)

	fmt.Printf("win.Resize(x:%v, y: %v)\n", 1340, 720)
	win.Resize(1340, 720)

	win.Focus()

	if close {
		fmt.Println("win.Close()")
		win.Close()
	}

	fmt.Println("Window tests OK")
}
