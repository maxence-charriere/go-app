![app](https://github.com/murlokswarm/app/blob/master/assets/logo/mur_logo-alt2-gris-github.png)
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
    - Android: Contribution welcome
    - Windows: Contribution welcome

3.  Get Xcode if you develop for MacOS (mandatory for Apple frameworks): 
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
type Hello struct {
	Greeting string
}

func (h *Hello) Render() string {
	return `
<div class="WindowLayout">    
    <div class="HelloBox">
        <h1>
            Hello,
            <span>{{if .Greeting}}{{html .Greeting}}{{else}}World{{end}}</span>
        </h1>
        <input type="text" placeholder="What is your name?" _onchange="OnInputChange" />
    </div>
</div>
    `
}

func (h *Hello) OnInputChange(arg app.ChangeArg) {
	h.Greeting = arg.Value
	app.Render(h)
}

func init() {
	app.RegisterComponent(&Hello{})
}
```

### Write the main
```go
func main() {
	app.OnLaunch = func() {
		win := app.NewWindow(app.Window{
			Title:          "Hello World",
			Width:          1280,
			Height:         720,
			TitlebarHidden: true,
		})

		hello := &Hello{}
		win.Mount(hello)
	}

	app.Run()
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
    overflow: hidden;
}

.WindowLayout {
    position: fixed;
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
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
