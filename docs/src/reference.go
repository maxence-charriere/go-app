package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type reference struct {
	app.Compo
}

func newReference() app.Composer {
	return &reference{}
}

func (r *reference) Render() app.UI {
	return newPage().
		Index(newGodocIndex()).
		Content(
			app.Div().
				Class("hspace-out-stretch").
				Body(newGodoc()),
		).
		IssueTitle("API reference")
}

type godocIndex struct {
	app.Compo

	rawHTML   string
	fragment  string
	isLoading bool
	err       error
}

func newGodocIndex() app.UI {
	return &godocIndex{}
}

func (i *godocIndex) OnPreRender(ctx app.Context) {
	i.init(ctx)
}

func (i *godocIndex) OnNav(ctx app.Context) {
	i.init(ctx)

	fragment := ctx.Page.URL().Fragment
	if fragment == "" {
		fragment = "top"
	}
	ctx.ScrollTo(fragment)
}

func (i *godocIndex) init(ctx app.Context) {
	i.unfocusLink(ctx)
	if i.rawHTML != "" {
		i.focusLink(ctx)
		return
	}

	ctx.Page.SetTitle("API reference for building a Progressive Web App with Go")
	ctx.Page.SetDescription("The API reference for building a Progressive Web App (PWA) with the go-app Go (Golang) package. ")

	i.load(ctx)
}

func (i *godocIndex) load(ctx app.Context) {
	i.isLoading = true
	i.err = nil
	i.Update()

	ctx.Async(func() {
		html, err := get(ctx, "/web/godoc-index.html")

		i.Defer(func(ctx app.Context) {
			i.rawHTML = string(html)
			i.err = err
			i.isLoading = false

			fragment := i.linkID(ctx.Page.URL().Fragment)
			if fragment == "" {
				fragment = "top"
			}

			i.Update()
			i.Defer(i.focusLink)
			ctx.ScrollTo(fragment)
		})
	})
}

func (i *godocIndex) Render() app.UI {
	return app.Div().
		Class("godoc-index").
		Body(
			app.Raw(i.rawHTML),
			newLoader().
				Title("Loading").
				Description("index").
				Size(48).
				Err(i.err).
				Loading(i.isLoading),
		)
}

func (i *godocIndex) focusLink(ctx app.Context) {
	fragment := app.Window().URL().Fragment
	if fragment == "" {
		return
	}

	link := app.Window().GetElementByID(i.linkID(fragment))
	if !link.Truthy() {
		return
	}

	link.Set("className", "focus")
	i.fragment = fragment
}

func (i *godocIndex) unfocusLink(ctx app.Context) {
	if i.fragment == "" {
		return
	}

	link := app.Window().GetElementByID(i.linkID(i.fragment))
	if !link.Truthy() {
		return
	}

	link.Set("className", "")
	i.fragment = ""
}

func (i *godocIndex) linkID(fragment string) string {
	return "src-" + fragment
}

type godoc struct {
	app.Compo

	isLoading   bool
	err         error
	rawHTML     string
	closeToggle func()
}

func newGodoc() *godoc {
	return &godoc{}
}

func (d *godoc) OnPreRender(ctx app.Context) {
	d.init(ctx)
}

func (d *godoc) OnNav(ctx app.Context) {
	d.init(ctx)
}

func (d *godoc) init(ctx app.Context) {
	if d.rawHTML != "" {
		return
	}
	d.load(ctx)
}

func (d *godoc) load(ctx app.Context) {
	d.isLoading = true
	d.err = nil
	d.Update()

	ctx.Async(func() {
		html, err := get(ctx, "/web/godoc.html")

		d.Defer(func(ctx app.Context) {
			d.rawHTML = string(html)
			d.err = err
			d.isLoading = false

			fragment := ctx.Page.URL().Fragment
			if fragment == "" {
				fragment = "top"
			}

			d.Update()
			d.Defer(d.setupToggle)
			ctx.ScrollTo(fragment)
		})
	})
}

func (d *godoc) setupToggle(ctx app.Context) {
	onToggle := app.FuncOf(d.onToggle)

	pkgOverview := app.Window().GetElementByID("pkg-overview")
	if !pkgOverview.Truthy() {
		return
	}
	pkgOverview.Call("addEventListener", "click", onToggle)

	pkgIndex := app.Window().GetElementByID("pkg-index")
	if !pkgIndex.Truthy() {
		return
	}
	pkgIndex.Call("addEventListener", "click", onToggle)

	d.closeToggle = func() {
		pkgOverview.Call("removeEventListener", "click", onToggle)
		pkgIndex.Call("removeEventListener", "click", onToggle)
		onToggle.Release()
	}

	if w, _ := app.Window().Size(); w >= 720 {
		pkgIndex.Set("className", "toggle")
	}
}

func (d *godoc) onToggle(src app.Value, args []app.Value) interface{} {
	d.Defer(func(ctx app.Context) {
		switch src.Get("className").String() {
		case "toggleVisible":
			src.Set("className", "toggle")

		case "toggle":
			src.Set("className", "toggleVisible")
		}
	})

	return nil
}

func (d *godoc) OnDismount() {
	if d.closeToggle != nil {
		d.closeToggle()
	}
}

func (d *godoc) Render() app.UI {
	return app.Div().
		Class("godoc").
		Body(
			app.Raw(d.rawHTML),
			newLoader().
				Class("page-loader").
				Class("fill").
				Title("Loading").
				Description("API reference").
				Err(d.err).
				Loading(d.isLoading),
		)
}
