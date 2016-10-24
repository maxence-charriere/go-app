# app
[![Build Status](https://travis-ci.org/murlokswarm/app.svg?branch=master)](https://travis-ci.org/murlokswarm/app)
[![Go Report Card](https://goreportcard.com/badge/github.com/murlokswarm/app)](https://goreportcard.com/report/github.com/murlokswarm/app)
[![Coverage Status](https://coveralls.io/repos/github/murlokswarm/app/badge.svg?branch=master)](https://coveralls.io/github/murlokswarm/app?branch=master)
[![GoDoc](https://godoc.org/github.com/murlokswarm/app?status.svg)](https://godoc.org/github.com/murlokswarm/app)

Package to build multiplatform apps with Go, HTML and CSS.

The idea is to use a web browser to handle just the UI part... 
Go for all the rest...

## Why?
A web browser is present on almost every platform. 
It's a part that never stops to evolve. 
Today it embeds enough power to handle beautiful UIs with smooth animations on desktop or the latest phones/tablets.

Go, is a simple, fast and well-built programming language. 
Plus, it is thinked from the ground to gracefully handle dependencies, tests and documentation.

## Install
```
// To be completed.
```

## How it works?
The way it works looks like a little bit like [React.js](https://facebook.github.io/react/).

### 1. Create a component
You create a component which satisfies the [markup.Componer](https://godoc.org/github.com/murlokswarm/markup#Componer) interface.
```go
type Componer interface {
	Render() string
}
```

### 2. Mount the component in a context
Then you mount the component into a [context](https://godoc.org/github.com/murlokswarm/app#Contexter) (Can be a window, a menu, a dock, etc...).
```go
type Contexter interface {
    ID() uid.ID

    Mount(c markup.Componer) // <- Call for mount a component.

    Resize(width float64, height float64)

    Move(x float64, y float64)

    SetIcon(path string)
}
```


### 3. Update a component
Finally, when updates to the component rendering are required, call [app.Render](https://godoc.org/github.com/murlokswarm/app#Render) on the component to be updated.

# Example
```go
// To be completed.
```

# How to
```go
// To be completed.
```

# Linked packages
[markup](https://github.com/murlokswarm/markup)
