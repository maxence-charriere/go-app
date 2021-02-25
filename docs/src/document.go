package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

type document struct {
	app.Compo

	Ipath string

	path        string
	description string
	document    string
	loading     bool
	err         error
}

func newDocument(path string) *document {
	return &document{Ipath: path}
}

func (d *document) Description(t string) *document {
	d.description = t
	return d
}

func (d *document) OnNav(ctx app.Context) {
	ctx.Dispatch(func() {
		d.scrollToFragment(ctx)
	})
}

func (d *document) load(ctx app.Context) {
	if d.Ipath == "" {
		return
	}
	d.path = d.Ipath

	d.loading = true
	d.err = nil
	d.Update()

	get := func() {
		var doc string
		var err error

		defer ctx.Dispatch(func() {
			if err != nil {
				d.err = err
			}

			d.document = doc
			d.loading = false
			d.Update()
			d.highlightCode(ctx)
			d.scrollToFragment(ctx)
		})

		doc, err = d.get(d.Ipath)
	}

	if app.IsServer {
		get()
		return
	}
	go get()
}

func (d *document) get(path string) (string, error) {
	res, err := http.Get(path)
	if err != nil {
		return "", errors.New("getting document failed").Wrap(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("reading document failed").Wrap(err)
	}

	return fmt.Sprintf("<div>%s</div>", parseMarkdown(b)), nil
}

func (d *document) highlightCode(ctx app.Context) {
	ctx.Dispatch(func() {
		app.Window().Get("Prism").Call("highlightAll")
	})
}

func (d *document) scrollToFragment(ctx app.Context) {
	ctx.Dispatch(func() {
		if ctx.Page.URL().Fragment == "" {
			fmt.Println("gonna scroll trop")
			app.Window().ScrollToID("top-doc")
			return
		}
		app.Window().ScrollToID(app.Window().URL().Fragment)
	})
}

func (d *document) Render() app.UI {
	if d.Ipath != d.path {
		d.Defer(d.load)
	}

	return app.Main().
		Class("pane").
		Body(
			app.Div().
				ID("top-doc").
				Class("document").
				Body(
					newLoader().
						Description(d.description).
						Err(d.err).
						Loading(d.loading),
					app.Raw(d.document),
					issue(d.Ipath),
					support(),
				),
		)
}
