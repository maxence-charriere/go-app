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
