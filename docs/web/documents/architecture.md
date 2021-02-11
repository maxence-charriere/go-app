# Architecture

[Progressive web apps](https://developers.google.com/web/progressive-web-apps) created with **go-app** are built with the following architecture:

![architecture diagram](/web/images/architecture.png)

## Web browser

The web browser is where the [app](#app) is running. Here is a list of well-known web browsers:

- [Chrome](https://www.google.com/chrome)
- [Safari](https://www.apple.com/safari)
- [Firefox](https://www.mozilla.org/firefox)
- [Electron (Chromium embedded)](https://www.electronjs.org/)

When a user navigates to the app domain, the web browser requests to the [server](#server) a webpage and its associated resources such as the [app](#app), [images, and styles](#static-resources). Then runs the app once all resources are gathered.

## Server

The server is the service that provides the [app](#app) and its required [resources](#static-resources) to [web browsers](#web-browser). It is where app metadata are defined, and where styles and scripts are linked. Servers are implemented by using the [app.Handler](/reference#Handler), as a standard [http.Handler](https://golang.org/pkg/net/http/#Handler):

```go
// +build !wasm

// The build instruction above ensures that this code is only used to build the
// server.

package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

func main() {
	// Standard http.Handler that handles all paths.
	http.Handle("/", &app.Handler{
		Name:        "Hello",                   // Metadata for PWA name.
		Title:       "Hello",                   // Metadata for page title.
		Description: "An Hello World! example", // Metadata for page description.
		Styles: []string{
			"/web/hello.css", // Inlude .css file.
		},
	})

	// Launches the server.
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

Servers are built with a standard build command:

```bash
go build
```

They can be deployed on multiple supports such as:

- Localhost
- Dedicated instances ([EC2](https://aws.amazon.com/ec2), [Google App Engine](https://cloud.google.com/appengine), ...)
- Containers ([Docker](https://www.docker.com/), [Kubernetes](https://kubernetes.io/))
- Cloud function ([AWS lambda](https://aws.amazon.com/lambda/), [Google Cloud functions](https://cloud.google.com/functions/), ...)

## App

The app is the program that contains the user interface that is displayed on the [web browser](#web-browser). It is built as a [WebAssembly](https://webassembly.org) binary that is served by the [server](#server) as a [static resource](/static-resources).

It contains at least one [component](/components): a customized, independent, and reusable UI element written in pure Go, that is associated with a URL path.

```go
// +build wasm

// The build instruction above ensures that this code is only used to build the
// app.

package main

import "github.com/maxence-charriere/go-app/v8/pkg/app"

// hello is a component that displays a simple "Hello World!". A component is
// created by embedding app.Compo into a struct.
type hello struct {
	app.Compo
}

// The Render method is where the component appearance is defined.
func (h *hello) Render() app.UI {
	return app.H1().Text("Hello World!")
}

func main() {
	app.Route("/", &hello{}) // hello component is associated with URL path "/".
	app.Run()                // Launches the app in the web browser.
}
```

Apps are built by specifying the `wasm` architecture and the `js` operatng system when using the build command:

```bash
GOARCH=wasm GOOS=js go build -o web/app.wasm
```

The app must be named `app.wasm` and be located in the `web` directory: a directory that is by default relative to the [server](#server) binary and where all [static resources](/static-resources) are located.

## Static resources

[Static resources](/static-resources) represent resources that are not dynamically generated such as:

- Styles (\*.css)
- Scripts (\*.js)
- Images
- Sounds
- Documents

Like the [app](#app), they are served by the [server](#server) to [web browsers](#web-browser) and are located in a directory named `web`, by default relative to the server binary.

```bash
.
├── app.go          # App source.
├── hello           # Server.
├── main.go         # Server source.
└── web             # Web directory containing all static resources.
    ├── app.wasm    # App.
    └── hello.css   # Style.
```

## Next

- [How to create a component](/components)
- [Deal with static resources](/static-resources)
- [API reference](/reference)
