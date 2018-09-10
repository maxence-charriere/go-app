// +build darwin,amd64

package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/examples/hello"
)

func main() {
	app.Import(&hello.Hello{})

	app.Run(&mac.Driver{
		URL: "/hello.Hello",
	})
}
