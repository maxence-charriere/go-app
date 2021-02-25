package main

import (
	"io/ioutil"
	"net/http"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/maxence-charriere/go-app/v8/pkg/errors"
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
			app.Article().
				Class("hspace-out-stretch").
				Body(newGodoc()),
		)
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
}

func (i *godocIndex) init(ctx app.Context) {
	i.unfocusLink(ctx)
	if i.rawHTML != "" {
		i.focusLink(ctx)
		return
	}

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

	loading     bool
	err         error
	rawHTML     string
	closeToggle func()
}

func newGodoc() *godoc {
	return &godoc{}
}

func (d *godoc) OnMount(ctx app.Context) {
	d.loading = true
	d.err = nil
	d.Update()

	go d.load(ctx)
}

func (d *godoc) load(ctx app.Context) {
	var html string
	var err error

	defer ctx.Dispatch(func() {
		if err != nil {
			d.err = err
		}

		d.rawHTML = html
		d.loading = false
		d.Update()

		ctx.Dispatch(d.setupToggle)
		ctx.Dispatch(d.scrollToSection)
	})

	path := "/web/godoc.html"

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
