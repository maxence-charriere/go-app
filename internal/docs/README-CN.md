# app

[![Build Status](https://travis-ci.org/murlokswarm/app.svg?branch=master)](https://travis-ci.org/murlokswarm/app)
[![Go Report Card](https://goreportcard.com/badge/github.com/murlokswarm/app)](https://goreportcard.com/report/github.com/murlokswarm/app)
[![Coverage Status](https://coveralls.io/repos/github/murlokswarm/app/badge.svg?branch=master)](https://coveralls.io/github/murlokswarm/app?branch=master)
[![awesome-go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#readme)
[![GoDoc](https://godoc.org/github.com/murlokswarm/app?status.svg)](https://godoc.org/github.com/murlokswarm/app)

一个使用
[Go](https://golang.org), [HTML](https://en.wikipedia.org/wiki/HTML5),
[CSS](https://en.wikipedia.org/wiki/Cascading_Style_Sheets)构建的多平台UI框架

![ui demo](https://github.com/murlokswarm/app/wiki/assets/ui-demo-large.gif)

[English](../../README.md) | 中文

## 目录

* [安装](#install)
* [支持平台](#support)
* [示例：Hello world](#hello)
* [架构](#architecture)
* [开发工具：Goapp](#goapp)
* [文档](#doc)
* [捐赠](#donate)

<a name="install"></a>

## 安装

```sh
# 安装:
go get -u -v github.com/murlokswarm/app/...

# 更新:
goapp update -v
```

<a name="support"></a>

## 支持平台

|平台|状态|
|:-|:-:|
|[MacOS](https://godoc.org/github.com/murlokswarm/app/drivers/mac#Driver)|✔|
|[Web](https://godoc.org/github.com/murlokswarm/app/drivers/web#Driver)|✔|
|Windows|[🔨](https://github.com/murlokswarm/app/issues/141)|
|Linux|✖|

<a name="hello"></a>

## Hello world

### 创建步骤

```sh
# 进入你的项目目录:
cd YOUR_REPO

# 初始化目录:
goapp mac init
```

### 代码

```go
// 你的项目目录/main.go

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

    // 使用mac驱动去运行Hello组件.
    app.Run(&mac.Driver{
        URL: "/hello",
    })
}
```

### 构建并运行

```sh
# 构建并运行debug模式:
goapp mac run -d
```

查看 [完整示例](https://github.com/murlokswarm/app/tree/master/examples/hello).

<a name="architecture"></a>

## 架构

![ui architecture](https://github.com/murlokswarm/app/wiki/assets/architecture.png)

### 元素

一个 [元素](https://godoc.org/github.com/murlokswarm/app#Elem)代表一个UI组件, 部分元素可以
[使用HTML去自定义](https://godoc.org/github.com/murlokswarm/app#ElemWithCompo)

目录:

* [Windows](https://godoc.org/github.com/murlokswarm/app#NewWindow)
* [Pages](https://godoc.org/github.com/murlokswarm/app#NewPage)
* [Context menus](https://godoc.org/github.com/murlokswarm/app#NewContextMenu)
* [Menubar](https://godoc.org/github.com/murlokswarm/app#MenuBar)
* [Status menu](https://godoc.org/github.com/murlokswarm/app#NewStatusMenu)
* [Dock](https://godoc.org/github.com/murlokswarm/app#Dock)

其余的一些简单示例:

* [Notifications](https://godoc.org/github.com/murlokswarm/app#NewNotification)
* [FilePanel](https://godoc.org/github.com/murlokswarm/app#NewFilePanel)
* [SaveFilePanel](https://godoc.org/github.com/murlokswarm/app#NewSaveFilePanel)
* [Share](https://godoc.org/github.com/murlokswarm/app#NewShare)

### 组件

[组件](https://godoc.org/github.com/murlokswarm/app#Compo)代表一个可以独立、可复用的UI组件. 它暴露的UI的HTML可以通过Go的一些基础库中提供的[模板语法](https://golang.org/pkg/text/template/)去进行自定义。
组件能够在
[元素](https://godoc.org/github.com/murlokswarm/app#ElemWithCompo) 里使用并且支持HTML自定义化。

### 驱动

[驱动](https://godoc.org/github.com/murlokswarm/app#Driver)代表app后台的具体运行方式。它暴露一些`Go`的操作方法去创建/修改UI和调用它们，并且会针对于特定于平台进行实现。
<a name="goapp"></a>

## 官方cli工具:Goapp

Goapp是一个用来构建和运行通过`app`生成pakage的应用的官方命令行工具。

根据平台的不同，必须打包应用程序才能进行部署
和发布。打包的应用程序通常不由终端管理
当我们想要监视日志或停止执行时，可能会出现问题。

`Goapp`可以通过终端打包应用程序以及运行他们,于此同时还能保持日志并管理它们的生命周期。

示例:

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

## 文档

* [Godoc](https://godoc.org/github.com/murlokswarm/app)
  * [mac](https://godoc.org/github.com/murlokswarm/app/drivers/mac)
  * [web](https://godoc.org/github.com/murlokswarm/app/drivers/web)
* [Wiki](https://github.com/murlokswarm/app/wiki)
  * [快速开始：MacOS](https://github.com/murlokswarm/app/wiki/Getting-started-with-MacOS)
  * [快速开始： web](https://github.com/murlokswarm/app/wiki/Getting-started-with-web)
  * [如何使用css？](https://github.com/murlokswarm/app/wiki/CSS)
* [示例](https://github.com/murlokswarm/app/tree/master/examples)
  * [hello](https://github.com/murlokswarm/app/tree/master/examples/hello)
  * [nav](https://github.com/murlokswarm/app/tree/master/examples/nav)
  * [menu](https://github.com/murlokswarm/app/tree/master/examples/menu)
  * [status menu](https://github.com/murlokswarm/app/tree/master/examples/statusmenu)
  * [dock](https://github.com/murlokswarm/app/tree/master/examples/dock)
  * [drag and drop](https://github.com/murlokswarm/app/tree/master/examples/dragdrop)
  * [actions/events](https://github.com/murlokswarm/app/tree/master/examples/action-event)
  * [test](https://github.com/murlokswarm/app/tree/master/examples/test)

<a name="donate"></a>

## 捐赠

如果这个项目可以帮助你建立好的用户界面，你可以通过以下方式赞助我！ :)

|Crypto|Address|
|-|-|
|[Ethereum (ETH)](https://www.coinbase.com/addresses/5b483b8df2ba04096454ea62)|0x789D63B8869783a15bbFb43331a192DdeC4bDE53|
|[Bitcoin (BTC)](https://www.coinbase.com/addresses/5b483f32bec71f034450c264)|3PRMM9fj7yq9gHxgk2svewWF9BkzzGPa1b|
