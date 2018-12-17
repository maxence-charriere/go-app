// +build !js

package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac"
	"github.com/murlokswarm/app/drivers/web"
)

func main() {
	app.Import(
		&NavPane{},
		&Hello{},
		&Open{},
		&Window{},
	)

	entryCompo := "open"

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
				newWindow("main", entryCompo, false)

				app.NewMsg("app-open").WithValue(appOpen{
					From: "OnRun()",
					Time: time.Now(),
				}).Post()
			},

			OnReopen: func(hasVisibleWindow bool) {
				if !hasVisibleWindow {
					newWindow("main", entryCompo, false)
				}

				app.NewMsg("app-open").WithValue(appOpen{
					From: fmt.Sprintf("OnReopen(%v)", hasVisibleWindow),
					Time: time.Now(),
				}).Post()
			},

			OnFilesOpen: func(filenames []string) {
				app.NewMsg("app-open").WithValue(appOpen{
					From: fmt.Sprintf("OnFilesOpen(%v)", app.Pretty(filenames)),
					Time: time.Now(),
				}).Post()
			},

			OnURLOpen: func(u *url.URL) {
				app.NewMsg("app-open").WithValue(appOpen{
					From: fmt.Sprintf("OnURLOpen(%s)", u),
					Time: time.Now(),
				}).Post()
			},

			OnQuit: func() {
				app.Log("Goodbye")
			},
		})
	}
}
