package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/drivers/web"
)

func main() {
	app.Import(
		&NavPane{},
		&Hello{},
		&Window{},
	)

	defaultURL := "window"

	app.NewSubscriber().
		Subscribe(app.Closed, func() { app.Log("goodbye") })

	app.Run(
		&mac.Driver{
			URL: defaultURL,
			Settings: mac.Settings{
				SupportedFiles: []mac.FileType{
					{
						Name: "image",
						Role: mac.Viewer,
						UTIs: []string{"public.jpeg"},
					},
				},
				URLScheme: "goapp-demo",
			},
		},
		&web.Driver{
			URL: defaultURL,
		},
	)
}
