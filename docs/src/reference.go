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
		Submenu(newGodocMenu()).
		OverlayMenu(Menu()).
		Content(&godoc{})
}

type godocMenu struct {
	app.Compo

	rawHTML string
	loading bool
	err     error
}

func newGodocMenu() app.UI {
	return &godocMenu{}
}

func (m *godocMenu) OnMount(ctx app.Context) {
	m.loading = true
	m.err = nil
	m.Update()

	go m.loadMenu()
}

func (m *godocMenu) loadMenu() {
	var html string
	var err error

	defer app.Dispatch(func() {
		if err != nil {
			m.err = err
		}

		m.rawHTML = html
		m.loading = false
		m.Update()
	})

	path := "/web/documents/godoc-index.html"
	res, err := http.Get(path)
	if err != nil {
		err = errors.New("getting reference menu failed").Wrap(err)
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		app.Log("%s", errors.New("reading reference menu failed").
			Tag("path", path).
			Wrap(err))
		return
	}

	html = string(b)
}

func (m *godocMenu) Render() app.UI {
	return app.Aside().
		Class("pane").
		Class("godoc-menu").
		Body(
			app.H1().Text("Index"),
			newLoader().
				Description("reference menu").
				Err(m.err).
				Loading(m.loading),
			app.Section().Body(
				app.Raw(m.rawHTML),
			),
		)
}

type godoc struct {
	app.Compo

	loading     bool
	err         error
	rawHTML     string
	closeToggle func()
}

func (d *godoc) OnMount(ctx app.Context) {
	d.loading = true
	d.err = nil
	d.Update()

	go d.load()
}

func (d *godoc) load() {
	var html string
	var err error

	defer app.Dispatch(func() {
		if err != nil {
			d.err = err
		}

		d.rawHTML = html
		d.loading = false
		d.Update()

		app.Dispatch(d.setupToggle)
		app.Dispatch(d.scrollToSection)
	})

	path := "/web/documents/godoc.html"

	res, err := http.Get(path)
	if err != nil {
		err = errors.New("retrieving reference failed").
			Tag("path", path).
			Wrap(err)
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.New("reading reference failed").
			Tag("path", path).
			Wrap(err)
		return
	}

	html = string(b)
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

func (d *godoc) scrollToSection() {
	app.Window().ScrollToID(app.Window().URL().Fragment)
}

func (d *godoc) OnDismount() {
	if d.closeToggle != nil {
		d.closeToggle()
	}
}

func (d *godoc) Render() app.UI {
	return app.Main().
		Class("pane").
		Class("godoc").
		Body(
			newLoader().
				Description("reference").
				Err(d.err).
				Loading(d.loading),
			app.Raw(d.rawHTML),
		)
}
