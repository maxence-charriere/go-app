# app

[![Build Status](https://travis-ci.org/murlokswarm/app.svg?branch=master)](https://travis-ci.org/murlokswarm/app)
[![Go Report Card](https://goreportcard.com/badge/github.com/murlokswarm/app)](https://goreportcard.com/report/github.com/murlokswarm/app)
[![Coverage Status](https://coveralls.io/repos/github/murlokswarm/app/badge.svg?branch=master)](https://coveralls.io/github/murlokswarm/app?branch=master)
[![awesome-go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#readme)
[![GoDoc](https://godoc.org/github.com/murlokswarm/app?status.svg)](https://godoc.org/github.com/murlokswarm/app)

A multi-platform UI framework that uses
[Go](https://golang.org), [HTML](https://en.wikipedia.org/wiki/HTML5) and
[CSS](https://en.wikipedia.org/wiki/Cascading_Style_Sheets).

![ui demo](https://github.com/murlokswarm/app/wiki/assets/ui-demo-large.gif)

## Table of Contents

* [Install](#install)
* [Supported platforms](#support)
* [Hello world](#hello)
* [Architecture](#architecture)
* [Goapp](#goapp)
* [Documentation](#doc)
* [Donate](#donate)

<a name="install"></a>

## Install

```sh
# Install:
go get -u -v github.com/murlokswarm/app/...

# Update:
goapp update -v
```

<a name="support"></a>

## Supported platforms

|Platform|Status|
|:-|:-:|
|[MacOS](https://godoc.org/github.com/murlokswarm/app/drivers/mac#Driver)|✔|
|[Web](https://godoc.org/github.com/murlokswarm/app/drivers/web#Driver)|✔|
|Windows|[🔨](https://github.com/murlokswarm/app/issues/141)|
|Linux|✖|

<a name="hello"></a>

## Hello world

### Setup

```sh
# Go to your repository:
cd YOUR_REPO

# Init the repo:
goapp mac init
```

### Code

```go
// YOUR_REPO/main.go

// Hello compo.
type Hello struct {
    Name string
}

func (h *Hello) Render() string {
    return `
<div class="Hello">
    <h1>
        Hello
        {{if .Name}}
            {{.Name}}
        {{else}}
            world
        {{end}}!
    </h1>
    <input value="{{.Name}}" placeholder="Write a name..." onchange="Name" autofocus>
</div>
    `
}

func main() {
    app.Import(&Hello{})

    // Use mac driver with Hello compo.
    app.Run(&mac.Driver{
        URL: "/hello",
    })
}
```

### Build and run

```sh
# Build and run with debug mode:
goapp mac run -d
```

View [full example](https://github.com/murlokswarm/app/tree/master/examples/hello).

<a name="architecture"></a>

## Architecture

![ui architecture](https://github.com/murlokswarm/app/wiki/assets/architecture.png)

### Elem

An [elem](https://godoc.org/github.com/murlokswarm/app#Elem) represents an UI
element to be displayed. Some can be
[customized with HTML](https://godoc.org/github.com/murlokswarm/app#ElemWithCompo)
content:

* [Windows](https://godoc.org/github.com/murlokswarm/app#NewWindow)
* [Pages](https://godoc.org/github.com/murlokswarm/app#NewPage)
* [Context menus](https://godoc.org/github.com/murlokswarm/app#NewContextMenu)
* [Menubar](https://godoc.org/github.com/murlokswarm/app#MenuBar)
* [Status menu](https://godoc.org/github.com/murlokswarm/app#NewStatusMenu)
* [Dock](https://godoc.org/github.com/murlokswarm/app#Dock)

Others are simple:

* [Notifications](https://godoc.org/github.com/murlokswarm/app#NewNotification)
* [FilePanel](https://godoc.org/github.com/murlokswarm/app#NewFilePanel)
* [SaveFilePanel](https://godoc.org/github.com/murlokswarm/app#NewSaveFilePanel)
* [Share](https://godoc.org/github.com/murlokswarm/app#NewShare)

### Compo

A [compo](https://godoc.org/github.com/murlokswarm/app#Compo) represents an
independent and reusable piece of UI. It exposes an HTML representation of the
UI that can be customized by the
[template syntax](https://golang.org/pkg/text/template/) defined in the Go
standard library. Compos are loaded into
[elems](https://godoc.org/github.com/murlokswarm/app#ElemWithCompo) that support
HTML customization.

### Driver

A [driver](https://godoc.org/github.com/murlokswarm/app#Driver) represents the
app backend. It exposes Go operations to create/modify the UI and calls their
platform specific implementations.

<a name="goapp"></a>

## Goapp

Goapp is a CLI tool to build and run apps built with the app package.

Depending on the platform, apps must be packaged in order to be deployed and
distributed. Packaged applications are usually not managed by a terminal, which
can be an issue when we want to monitor the logs or stop their execution with
system signals.

Goapp can package apps and allows to run them while keeping logs and managing
their lyfecycle within the terminal.

Examples:

```sh
goapp -h         # Help.
goapp mac -h     # Help for MasOS commands.
goapp mac run -h # Help for MasOS run command.

goapp mac run    # Run MacOS .app.
goapp mac run -d # Run MacOS .app with debug.

goapp web run    # Run a web server.
goapp web run -b # Run a web server and launch the main page in the default browser.
```

<a name="doc"></a>

## Documentation

* [Godoc](https://godoc.org/github.com/murlokswarm/app)
  * [mac](https://godoc.org/github.com/murlokswarm/app/drivers/mac)
  * [web](https://godoc.org/github.com/murlokswarm/app/drivers/web)
* [Wiki](https://github.com/murlokswarm/app/wiki)
  * [Getting started with MacOS](https://github.com/murlokswarm/app/wiki/Getting-started-with-MacOS)
  * [Getting started with web](https://github.com/murlokswarm/app/wiki/Getting-started-with-web)
* [Examples](https://github.com/murlokswarm/app/tree/master/examples)
  * [hello](https://github.com/murlokswarm/app/tree/master/examples/hello)
  * [nav](https://github.com/murlokswarm/app/tree/master/examples/nav)
  * [menu](https://github.com/murlokswarm/app/tree/master/examples/menu)
  * [status menu](https://github.com/murlokswarm/app/tree/master/examples/statusmenu)
  * [dock](https://github.com/murlokswarm/app/tree/master/examples/dock)
  * [drag and drop](https://github.com/murlokswarm/app/tree/master/examples/dragdrop)
  * [actions/events](https://github.com/murlokswarm/app/tree/master/examples/action-event)
  * [test](https://github.com/murlokswarm/app/tree/master/examples/test)

<a name="donate"></a>

## Donate

If this project helps you build awesome UI, you can help me grow my cryptos :)

|Crypto|Address|
|-|-|
|[Ethereum (ETH)](https://www.coinbase.com/addresses/5b483b8df2ba04096454ea62)|0x789D63B8869783a15bbFb43331a192DdeC4bDE53|
|[Bitcoin (BTC)](https://www.coinbase.com/addresses/5b483f32bec71f034450c264)|3PRMM9fj7yq9gHxgk2svewWF9BkzzGPa1b|
