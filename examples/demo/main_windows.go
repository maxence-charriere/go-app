// +build !js

package main

import (
	"net/url"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/web"
	"github.com/murlokswarm/app/drivers/win"
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
		app.Run(&win.Driver{
			Settings: win.Settings{
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

			OnRun: func() {
				newWindow("main", entryCompo, false)
			},

			OnFilesOpen: func(filenames []string) {
				app.Log("opened from:", filenames)
			},

			OnURLOpen: func(u *url.URL) {
				app.Log("opened with", u)
			},
		})
	}
}
