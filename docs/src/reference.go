package main

import (
	"io/ioutil"
	"net/http"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

func Reference() app.UI {
	return &reference{}
}

type reference struct {
	app.Compo
}

func (r *reference) Render() app.UI {
	return app.Shell().
		Class("app-background").
		Menu(Menu()).
		Submenu(GodocMenu()).
		OverlayMenu(Menu()).
		Content(Godoc())
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
	return app.Div().
		Class("godoc-menu").
		Body(
			app.H2().Text("Reference"),
			app.If(m.rawHTML != "",
				app.Raw(m.rawHTML),
			),
		)
}

func Godoc() app.UI {
	return &godoc{}
}

type godoc struct {
	app.Compo

	rawHTML string
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
}

func (d *godoc) Render() app.UI {
	return app.Div().
		Class("godoc").
		Body(
			app.H1().Text("Reference"),
			app.Raw(d.rawHTML),
		)
}
