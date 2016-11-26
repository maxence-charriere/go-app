![app](https://github.com/murlokswarm/app/blob/master/resources/logo/mur_logo-alt2-gris-small.png)
# app
[![Build Status](https://travis-ci.org/murlokswarm/app.svg?branch=master)](https://travis-ci.org/murlokswarm/app)
[![Go Report Card](https://goreportcard.com/badge/github.com/murlokswarm/app)](https://goreportcard.com/report/github.com/murlokswarm/app)
[![Coverage Status](https://coveralls.io/repos/github/murlokswarm/app/badge.svg?branch=master)](https://coveralls.io/github/murlokswarm/app?branch=master)
[![GoDoc](https://godoc.org/github.com/murlokswarm/app?status.svg)](https://godoc.org/github.com/murlokswarm/app)

Package to build multiplatform apps with Go, HTML and CSS.

## Install
1. Install Golang:
    - [golang.org](https://golang.org/doc/install)
    - [Homebrew (MacOS)](http://www.golangbootcamp.com/book/get_setup)

2. Get a driver:
    - MacOS: ```go get -u github.com/murlokswarm/mac```
    - IOS: Not available yet
    - Android: Not available yet
    - Windows: Not available yet

3.  Get Xcode  (If you develop on MacOS): 
    - [Xcode](https://itunes.apple.com/us/app/xcode/id497799835?mt=12)

## Getting started
![hello](https://github.com/murlokswarm/examples/blob/master/mac/hello/capture-1.png)

### Import a driver
```Go
import (
	_ "github.com/murlokswarm/mac"
)
```

### Create a component
```Go
// Hello implements app.Componer interface.
type Hello struct {
	Greeting string
}

// Render returns the HTML markup that describes the appearance of the
// component.
// It supports standard HTML and extends it slightly to handle other component
// declaration or Golang callbacks.
// Can be templated following rules from https://golang.org/pkg/text/template.
func (h *Hello) Render() string {
	return `
<div class="WindowLayout">    
    <div class="HelloBox">
        <h1>
            Hello,
            <span>{{if .Greeting}}, {{html .Greeting}}{{end}}</span>
        </h1>
        <input type="text" placeholder="What is your name?" _onchange="OnInputChange" />
    </div>
</div>
    `
}

// OnInputChange is the handler called when an onchange event occurs.
// In the HTML markup, a Go component method is target by prefixing the event with "_".
// eg _onchange.
func (h *Hello) OnInputChange(arg app.ChangeArg) {
	h.Greeting = arg.Value // Changing the greeting.
	app.Render(h)          // Tells the app to update the rendering of the component.
}

func init() {
	// Registers the Hello component.
	// Allows the app to create a Hello component when it finds its declaration
	// into a HTML markup.
	app.RegisterComponent(&Hello{})
}
```

### Write the main
```go
func main() {
    // When app is launched
	app.OnLaunch = func() {
		// Creates a window context.
		win := app.NewWindow(app.Window{
			Title:          "Hello World",
			Width:          1280,
			Height:         720,
			TitlebarHidden: true,
		})

		hello := &Hello{} // Creates a hello component.
		win.Mount(hello)  // Mounts the hello component into the window context.
	}

	app.Run() // Runs the app.
}

```

### Style your component
```css
/* In resources/css/hello.css.*/

body {
    background-image: url("../bg1.jpg");
    background-size: cover;
    background-position: center;
    color: white;
}

.WindowLayout {
    position: absolute;
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
}

.HelloBox {
    padding: 20pt;
}

h1 {
    font-weight: 300;
}

input {
    width: 100%;
    padding: 5pt;
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
```

[Full example](https://github.com/murlokswarm/examples/tree/master/mac/hello)

## Documentation
- [Wiki](https://github.com/murlokswarm/app/wiki)
- [GoDoc](https://godoc.org/github.com/murlokswarm/app)
