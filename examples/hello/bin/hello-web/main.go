package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
	"github.com/murlokswarm/app/examples/hello"
)

func main() {
	app.Import(&hello.Hello{})

	app.Run(&web.Driver{
		URL: "/hello.Hello",
	})
}
