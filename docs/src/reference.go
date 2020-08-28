package main

import (
	"io/ioutil"
	"net/http"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

type reference struct {
	app.Compo
}

func newReference() app.UI {
	return &reference{}
}

func (r *reference) Render() app.UI {
	return app.Shell().
		Class("app-background").
		Menu(Menu()).
		Submenu(GodocMenu()).
		OverlayMenu(Menu()).
		Content(&godoc{})
}

func GodocMenu() app.UI {
	return &godocMenu{}
}

type godocMenu struct {
	app.Compo

	rawHTML string
}

func (m *godocMenu) OnMount(ctx app.Context) {
	m.loadMenu()
}

func (m *godocMenu) loadMenu() {
	path := "/web/godoc-index.html"

	res, err := http.Get(path)
	if err != nil {
		app.Log("%s", errors.New("retrieving menu content failed").
			Tag("path", path).
			Wrap(err))
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		app.Log("%s", errors.New("reading menu content failed").
			Tag("path", path).
			Wrap(err))
		return
	}

	m.rawHTML = string(b)
	m.Update()
}

func (m *godocMenu) Render() app.UI {
	return app.Aside().
		Class("layout").
		Class("godoc-menu").
		Body(
			app.Div().Class("header"),
			app.Div().
				Class("content").
				Body(
					app.Section().Body(
						app.H1().Text("Table of contents"),
						app.If(m.rawHTML != "",
							app.Raw(m.rawHTML),
						),
					),
				),
		)
}

type godoc struct {
	app.Compo

	rawHTML     string
	closeToggle func()
}

func (d *godoc) OnMount(ctx app.Context) {
	d.loadMenu()
}

func (d *godoc) loadMenu() {
	path := "/web/godoc.html"

	res, err := http.Get(path)
	if err != nil {
		app.Log("%s", errors.New("retrieving content failed").
			Tag("path", path).
			Wrap(err))
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		app.Log("%s", errors.New("reading content failed").
			Tag("path", path).
			Wrap(err))
		return
	}

	d.rawHTML = string(b)
	d.Update()
	app.Dispatch(d.setupToggle)
}

func (d *godoc) setupToggle() {
	onToggle := app.FuncOf(d.onToggle)

	pkgOverview := app.Window().GetElementByID("pkg-overview")
	if !pkgOverview.Truthy() {
		panic(errors.New("pkg-overview elem not found"))
	}
	pkgOverview.Call("addEventListener", "click", onToggle)

	pkgIndex := app.Window().GetElementByID("pkg-index")
	if !pkgIndex.Truthy() {
		panic(errors.New("pkg-index elem not found"))
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
	app.Dispatch(func() {
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
	return app.Main().
		Class("layout").
		Class("godoc").
		Body(
			app.Div().Class("header"),
			app.Div().
				Class("content").
				Body(
					app.Section().Body(
						app.Raw(d.rawHTML),
					),
				),
		)
}
