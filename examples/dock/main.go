package main

import "github.com/murlokswarm/app"

func main() {
	app.OnLaunch = func() {
		if dock, ok := app.Dock(); ok {
			menu := &DockMenu{}
			dock.Mount(menu)
		}
	}
	app.Run()
}
