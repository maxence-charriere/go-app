# app
[![Build Status](https://travis-ci.org/murlokswarm/app.svg?branch=master)](https://travis-ci.org/murlokswarm/app)
[![Go Report Card](https://goreportcard.com/badge/github.com/murlokswarm/app)](https://goreportcard.com/report/github.com/murlokswarm/app)
[![Coverage Status](https://coveralls.io/repos/github/murlokswarm/app/badge.svg?branch=master)](https://coveralls.io/github/murlokswarm/app?branch=master)
[![GoDoc](https://godoc.org/github.com/murlokswarm/app?status.svg)](https://godoc.org/github.com/murlokswarm/app)
[![Contribute PATREON](https://img.shields.io/badge/Contribute-PATREON-f76959.svg)](https://www.patreon.com/maxencecharriere)

Package to build MacOS and Web apps using **Go**, **HTML** and **CSS**.

## Table of Contents
- [Install](#install)
- [Hello world](#hello)
- [Drivers](#drivers)
- [Documentation](#doc)
- [Examples](#examples)

<a name="install"></a>

## Install

```bash
# MacOS 10.12 and above.
xcode-select --install
```

```bash
# Get package.
go get -u -v github.com/murlokswarm/app/...
```

<a name="hello"></a>

## Hello world

![hello](https://github.com/murlokswarm/app/wiki/assets/hello.gif)

### Main
```go
func main() {
	app.Run(&mac.Driver{
		OnRun: func() {
			newWindow()
		},

		OnReopen: func(hasVisibleWindow bool) {
			if !hasVisibleWindow {
				newWindow()
			}
		},
	})
}

func newWindow() {
	app.NewWindow(app.WindowConfig{
		Title:           "hello world",
		TitlebarHidden:  true,
		Width:           1280,
		Height:          768,
		BackgroundColor: "#21252b",
		DefaultURL:      "/Hello",
	})
}
```

[app.Run](https://godoc.org/github.com/murlokswarm/app#Run) starts the app. 
It takes an 
[app.Driver](https://godoc.org/github.com/murlokswarm/app#Driver) as argument. 
Here we use the
[MacOS driver](https://godoc.org/github.com/murlokswarm/app/drivers/mac#Driver) 
implementation.
See [other drivers](#drivers).

When creating the window, we set the ```DefaultURL``` to our Hello component 
struct name: ```/Hello```.
It will make the component loaded when the window is displayed.

### Component
```go
func init() {
	app.Import(&Hello{})
}

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
	<input value="{{.Name}}" placeholder="Say something..." onchange="Name" autofocus>
</div>
	`
}
```
Components are structs that implement the 
[app.Component](https://godoc.org/github.com/murlokswarm/app#Component) 
interface.

Render method returns a string that contains HTML5.
It can be templated following Go standard template syntax:
- [text/template](https://golang.org/pkg/text/template/)
- [html/template](https://golang.org/pkg/html/template/)

HTML events like ```onchange``` are mapped to the targetted component 
field or method.
Here, ```onchange``` is mapped to the field ```Name```.

### CSS

```css
/* [PACKAGE PATH]/resources/css/hello.css */

.Hello {
    /* Your CSS */
}
```

Because, we want a stylish Hello world, we define the CSS that will give us some 
cool look.


CSS files are located in ```[PACKAGE PATH]/resources/css/``` directory.
All .css files within that directory will be included.

See the 
[full example](https://github.com/murlokswarm/app/tree/master/examples/hello).

<a name="drivers"></a>

## Drivers
A driver contains specific code that allows the app package to work on multiple 
platforms.

- [MacOS](https://godoc.org/github.com/murlokswarm/app/drivers/mac)
- [Web](https://godoc.org/github.com/murlokswarm/app/drivers/web) - *run on the top of [gopherjs](https://github.com/gopherjs/gopherjs)*

Other drivers will come in the future.

<a name="doc"></a>

## Documentation
- [GoDoc](https://godoc.org/github.com/murlokswarm/app)
- [v1 to v2 migration guide](https://github.com/murlokswarm/app/wiki/V1ToV2)

<a name="examples"></a>

## Examples
From package:
- [hello](https://github.com/murlokswarm/app/tree/master/examples/hello)
- [nav](https://github.com/murlokswarm/app/tree/master/examples/nav)
- [menu](https://github.com/murlokswarm/app/tree/master/examples/menu)
- [dock](https://github.com/murlokswarm/app/tree/master/examples/dock)
- [drag and drop](https://github.com/murlokswarm/app/tree/master/examples/dragdrop)
- [test](https://github.com/murlokswarm/app/tree/master/examples/test)

From community:
- [grocid/mistlur](https://github.com/grocid/mistlur) - use v1
- [grocid/passdesktop](https://github.com/grocid/passdesktop) - use v1

