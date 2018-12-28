package main

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/drivers/web"
	"github.com/murlokswarm/app/drivers/win"
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
			URL:       defaultURL,
			URLScheme: "goapp-demo",
			SupportedFiles: []mac.FileType{
				{
					Name: "image",
					Role: mac.Viewer,
					UTIs: []string{"public.jpeg"},
				},
			},
		},
		&win.Driver{
			URL: defaultURL,
			DefaultWindow: app.WindowConfig{
				BackgroundColor: "#ff0000",
			},
			URLScheme: "goapp-demo",
			SupportedFiles: []win.FileType{
				{
					Name: "murlok",
					Help: "A test extension for goapp",
					Icon: "like.png",
					Extensions: []win.FileExtension{
						{Ext: ".murlok"},
					},
				},
			},
		},
		&web.Driver{
			URL: defaultURL,
		},
	)
}
