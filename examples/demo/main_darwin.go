// +build !js

package main

import (
	"net/url"

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

			OnRun: func() {
				newWindow("main", entryCompo)
			},

			OnReopen: func(hasVisibleWindow bool) {
				if !hasVisibleWindow {
					newWindow("main", entryCompo)
				}
			},

			OnFilesOpen: func(filenames []string) {
				app.Log("opened from:", filenames)
			},

			OnURLOpen: func(u *url.URL) {
				app.Log("app opened with:", u)
			},
		})
	}
}
