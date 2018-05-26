package main

import (
	"os"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
	"github.com/murlokswarm/app/examples/test"
)

func main() {
	app.Loggers = []app.Logger{
		app.NewLogger(os.Stdout, os.Stderr, true),
	}

	app.Import(&test.Webview{})
	app.Import(&test.Menu{})

	app.Run(&web.Driver{
		DefaultURL: "/test.Webview",
	}, app.Logs())
}
