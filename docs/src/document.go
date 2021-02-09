package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/maxence-charriere/go-app/v7/pkg/app"
	"github.com/maxence-charriere/go-app/v7/pkg/errors"
)

type document struct {
	app.Compo

	path        string
	description string
	document    string
	loading     bool
	err         error
}

func newDocument(path string) *document {
	return &document{path: path}
}

func (d *document) Description(t string) *document {
	d.description = t
	return d
}

func (d *document) OnPreRender(ctx app.Context) {
	u := *ctx.Page.URL()
	u.Scheme = "http"
	u.Path = d.path

	d.document, d.err = d.get(u.String())
	fmt.Println("on prerender:", d.err)
	d.Update()
}

func (d *document) OnMount(ctx app.Context) {
	d.loading = true
	d.err = nil
	d.Update()

	go d.load(ctx)
}

func (d *document) load(ctx app.Context) {
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

	doc, err = d.get(d.path)
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
		app.Window().ScrollToID(app.Window().URL().Fragment)
	})
}

func (d *document) Render() app.UI {
	return app.Main().
		Class("pane").
		Class("document").
		Body(
			newLoader().
				Description(d.description).
				Err(d.err).
				Loading(d.loading),
			app.Raw(d.document),
			issue(d.path),
			support(),
		)
}

func parseMarkdown(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	return markdown.ToHTML(md, parser, nil)
}
