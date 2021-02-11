# Getting started

**go-app** is a package to build [progressive web apps (PWA)](https://developers.google.com/web/progressive-web-apps/) with [Go programming language](https://golang.org) and [WebAssembly](https://webassembly.org).

This document is about how to get started by showing how to write and build a simple `hello world!` PWA.

## Prerequisite

Using this package requires a fully functional Go environment installed with a **Go** version equal to or greater than **1.14**. Instructions about how to install and set up Go can be found on [golang.org](https://golang.org/doc/install).

Go installation can be checked with the following command in a terminal:

```bash
go version
```

## Install

Create a Go package for your PWA and change directory to the newly created location:

```bash
# Replace YOUR_PACKAGE by a you own package name.
# Eg. $GOPATH/src/src/hello
mkdir -p $GOPATH/src/YOUR_PACKAGE
cd $GOPATH/src/YOUR_PACKAGE
```

Then Initialize the **go module** and download the **go-app** package.

```bash
go mod init
go get -u github.com/maxence-charriere/go-app/v8/pkg/app
```

## User interface

Create the `app.go` file that will contain the [user interface](/architecture#ui) and write the following code:

```go
// +build wasm

// The UI is running only on a web browser. Therefore, the build instruction
// above is to compile the code below only when the program is built for the
// WebAssembly (wasm) architecture.

package main

import "github.com/maxence-charriere/go-app/v8/pkg/app"

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

// The main function is the entry point of the UI. It is where components are
// associated with URL paths and where the UI is started.
func main() {
	app.Route("/", &hello{}) // hello component is associated with URL path "/".
	app.Run()                // Launches the PWA.
}
```

## Server

Create the `main.go` file that will contain the [server](/architecture#server) and write the following code:

```go
// +build !wasm

// The server is a classic Go program that can run on various architecture but
// not on WebAssembly. Therefore, the build instruction above is to exclude the
// code below from being built on the wasm architecture.

package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

// The main function is the entry of the server. It is where the HTTP handler
// that serves the UI is defined and where the server is started.
//
// Note that because main.go and app.go are built for different architectures,
// this main() function is not in conflict with the one in
// app.go.
func main() {
	// app.Handler is a standard HTTP handler that serves the UI and its
	// resources to make it work in a web browser.
	//
	// It implements the http.Handler interface so it can seamlessly be used
	// with the Go HTTP standard library.
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Build and run

At this point the package should have the following content:

```bash
.
├── app.go
├── go.mod
├── go.sum
└── main.go

0 directories, 4 files
```

The last thing to do is to build the app. The build is done in 2 steps:

```bash
# Build the wasm program that contains the user interface.
GOARCH=wasm GOOS=js go build -o web/app.wasm
```

```bash
# Build the server that serves the wasm program and its resources:
go build
```

Note that when building the UI, the build output is explicitly set to `web/app.wasm`. The reason is that the [app.handler](/reference#Handler) tells the browser to load the UI from the `/web/app.wasm` path.

Once the UI and the server built, the package should have the following content:

```bash
.
├── app.go
├── go.mod
├── go.sum
├── hello
├── main.go
└── web
    └── app.wasm

1 directory, 6 files
```

Launch the server:

```bash
# ./SERVER_NAME
./hello
```

Finally, [navigate to the app](http://localhost:8000) into you web browser.

## Tips

The build process can be simplified by writing a makefile:

```makefile
build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build

run: build
	./hello
```

Then run:

```bash
make run
```

## Next

- [Understand go-app architecture](/architecture)
- [How to create a component](/components)
- [API reference](/reference)
