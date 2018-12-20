// +build !js

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

	entryCompo := "window"

	switch app.Kind {
	case "web":
		app.Run(&web.Driver{
			URL: entryCompo,
		})

	default:
		app.Run(&mac.Driver{
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

			URL: "window",

			OnQuit: func() {
				app.Log("Goodbye")
			},
		})
	}
}

// func nain() {
// 	app.Import(
// 		&NavPane{},
// 		&Hello{},
// 		&Window{},
// 	)

// 	defer app.NewSubscriber().
// 		Subscribe("files-openened", nil).
// 		Close()

// 	app.Run(
// 		mac.Config{},
// 		web.Config{},
// 	)

// 	// filter right config from target.
// }
