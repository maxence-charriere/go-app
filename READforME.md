# app

[![Build Status](https://travis-ci.org/murlokswarm/app.svg?branch=master)](https://travis-ci.org/murlokswarm/app)
[![Go Report Card](https://goreportcard.com/badge/github.com/murlokswarm/app)](https://goreportcard.com/report/github.com/murlokswarm/app)
[![Coverage Status](https://coveralls.io/repos/github/murlokswarm/app/badge.svg?branch=master)](https://coveralls.io/github/murlokswarm/app?branch=master)
[![GoDoc](https://godoc.org/github.com/murlokswarm/app?status.svg)]

A multi-platform UI framework that uses
[Go](https://golang.org), [HTML](https://en.wikipedia.org/wiki/HTML5) and
[CSS](https://en.wikipedia.org/wiki/Cascading_Style_Sheets).

![hello](https://github.com/murlokswarm/app/wiki/assets/app.gif)

## Table of Contents

* [Install](#install)
* [Supported platforms](#support)
* [Hello world](#hello)
* [Architecture](#architecture)
* [Build](#build)
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

<a name="architecture"></a>

## Architecture

![hello](https://github.com/murlokswarm/app/wiki/assets/architecture.png)

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

<a name="build"></a>

## Build

<a name="doc"></a>

## Documentation

<a name="donate"></a>

## Donate
