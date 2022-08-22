## Intro

go-app v9 brings new features like removing the need to manually call `compo.Update()` to update what is displayed on the screen. Unfortunately, it comes with some breaking changes.

## Changes

This document has for purpose to help you transition from go-app v8 to v9 by enumerating the API changes.

### General

- This package now requires [Go 1.18](https://golang.org/doc/install)

### Components

- Components are auto-updated when a component lifecycle event, an HTML event, or a dispatch occurs
- `Compo.Defer()` has been removed (now in Context)

### Context

- `Context` is now an interface
  - `ctx.Src` becomes `ctx.Src()`
  - `ctx.JSSrc` becomes `ctx.JSSrc()`
  - `ctx.AppUpdateAvailable` becomes `ctx.AppUpdateAvailable()`
  - `ctx.Page` becomes `ctx.Page()`
- `ctx.Dispatch(func())` becomes `ctx.Dispatch(func(app.Context))`
- UI elements are now always checked to be mounted before being dispatched

## API design decisions

- V9 release focuses mainly on getting the usage of the package more reactive.
- The API is designed to make unrecommended things difficult. A good example would be calling a `Context.Dispatch()` inside a `Render()` method.
- With the exception of `Handler` and JS-related things, everything that you might need is now available from [Context](/reference#Context).
