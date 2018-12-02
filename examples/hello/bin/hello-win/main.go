// +build windows

package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/win"
	"github.com/murlokswarm/app/examples/hello"
)

func main() {
	app.Import(&hello.Hello{})

	app.Run(&win.Driver{
		URL: "/hello.Hello",
	})
}
