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
		&Open{},
		&Window{},
		&Menu{},
		&ContextMenu{},
		&TestMenu{},
	)

	defaultURL := "menu"

	app.NewSubscriber().
		Subscribe(app.Closed, func() { app.Log("goodbye") })

	app.Run(
		&mac.Driver{
			URL: defaultURL,
			MenubarConfig: app.MenuBarConfig{
				CustomURLs: []string{"testmenu"},
			},
			DockURL:   "testmenu",
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
			URL:       defaultURL,
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
