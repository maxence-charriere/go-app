# app
[![Build Status](https://travis-ci.org/murlokswarm/app.svg?branch=master)](https://travis-ci.org/murlokswarm/app)
[![Go Report Card](https://goreportcard.com/badge/github.com/murlokswarm/app)](https://goreportcard.com/report/github.com/murlokswarm/app)
[![Coverage Status](https://coveralls.io/repos/github/murlokswarm/app/badge.svg?branch=master)](https://coveralls.io/github/murlokswarm/app?branch=master)
[![GoDoc](https://godoc.org/github.com/murlokswarm/app?status.svg)](https://godoc.org/github.com/murlokswarm/app)

Package to build multi platform apps using **Go**, **HTML** and **CSS**.

## Table of Contents
- [Install](#install)
	- [MacOS](#macOS)
- [Hello world](#hello)
- [Documentation](#doc)
- [Examples](#examples)

<a name="install"></a>

## Install

<a name="macOS"></a>

### MacOS 10.12 and above
```bash
# MacOS libraries.
xcode-select --install

# Get package.
go get -u github.com/murlokswarm/mac
```

<a name="hello"></a>

## Hello world
![hello](https://github.com/murlokswarm/app/wiki/assets/hello.gif)

### Component:
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
**Components** are structs that implement the 
[app.Component](https://godoc.org/github.com/murlokswarm/app#Component) 
interface.

**Render** method returns a string that contains **HTML5**.
It can be templated following Go standard template syntax:
- [text/template](https://golang.org/pkg/text/template/)
- [html/template](https://golang.org/pkg/html/template/)

**HTML events** like ```onchange``` are mapped to the targetted component 
**field** or **method**.
Here, ```onchange``` is mapped to the field ```Name```.

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

**app.Run** starts the app. It takes a 
[app.Driver](https://godoc.org/github.com/murlokswarm/app#Driver) as argument. 
Here we use the
[MacOS driver](https://godoc.org/github.com/murlokswarm/app/drivers/mac#Driver) 
implementation. 
Other drivers will be released in the futur.

When creating the window, we set the ```DefaultURL``` to our component struct 
name ```/Hello``` in order to have it loaded when the window shows.

### CSS
```css
body {
    font-family: 'Helvetica Neue', 'Segoe UI';
    color: white;
}

h1 {
    font-weight: 300;
    font-size: 42pt;
    max-width: calc(100% - 80px);
    padding: 0 40px;
}

input {
    width: 265px;
    padding: 8pt;
    margin-bottom: 5%;
    border: 0;
    border-left: 2px solid silver;
    outline: none;
    font-size: 14px;
    background: transparent;
    color: white;
}

input:focus {
    border-left-color: deepskyblue;
}

.Hello {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    background-image: url('../space.jpg');
    background-size: cover;
    background-position: center;
    height: 100%;
    width: 100%;
}
```

Because, we want a stylish Hello world, here we define the CSS that will give
us some cool look.

CSS files must be located in ```[PACKAGE PATH]/resources/css/``` directory.
All **.css** files within that directory will be included.

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
- [test](https://github.com/murlokswarm/app/tree/master/examples/test)

From community:
- [grocid/mistlur](https://github.com/grocid/mistlur) - use v1
- [grocid/passdesktop](https://github.com/grocid/passdesktop) - use v1

