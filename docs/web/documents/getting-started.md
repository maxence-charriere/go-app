## Intro

**go-app** is a package to build [progressive web apps (PWA)](https://developers.google.com/web/progressive-web-apps/) with the [Go programming language](https://golang.org) and [WebAssembly](https://webassembly.org).

You are about to learn how to get started with this package by **building and running an app that displays `Hello World`**.

## Prerequisite

Using this package **requires a fully functional Go environment installed with [Go 1.18](https://golang.org/doc/install) or greater**.

Go installation can be checked with the following command in a terminal:

```bash
go version
```

## Install

Create a Go package for your PWA and change directory to the newly created location:

```bash
mkdir -p $GOPATH/src/github.com/YOUR_GITHUB_ID/hello
cd $GOPATH/src/github.com/YOUR_GITHUB_ID/hello
```

Then Initialize the **go module** and download the **go-app** package.

```bash
go mod init
go get -u github.com/maxence-charriere/go-app/v10/pkg/app
```

## Code

Here is the code used to create a progressive web app that displays a simple Hello World.

### Hello component

```go
package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type hello struct {
	app.Compo
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *hello) Render() app.UI {
	return app.H1().Text("Hello World!")
}
```

### Main

```go
// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// The first thing to do is to associate the hello component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", func() app.Composer { return &hello{} })

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for precompiling
	// instructions.
	app.RunWhenOnBrowser()

	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
```

## Build and run

Running a progressive app with **go-app** requires 2 Go programs:

- A client that runs in a web browser
- A server that serves the client and its resources

At this point, the package has the following content:

```bash
.
├── go.mod
├── go.sum
└── main.go

0 directories, 4 files
```

### Build the client

```bash
GOARCH=wasm GOOS=js go build -o web/app.wasm
```

Note that the build output is explicitly set to `web/app.wasm`. The reason why is that the [Handler](/reference#Handler) expects the client to be a **static resource** located at the `/web/app.wasm` path.

### Build the server

```bash
go build
```

### Run the app

Now the client and server built, the package has the following content:

```bash
.
├── go.mod
├── go.sum
├── hello
├── main.go
└── web
    └── app.wasm

1 directory, 6 files
```

The server is ran with the following command:

```bash
./hello
```

The app is now accessible from a web browser at http://localhost:8000.

### Use a Makefile

The build process can be simplified by writing a makefile:

```makefile
build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build

run: build
	./hello
```

It can now be built and run with this single command:

```bash
make run
```

## Next

- [Architecture](/architecture)
- [Reference](/reference)
- [Build a component](/components)
