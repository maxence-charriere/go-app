// +build !js

package web

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"os"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/murlokswarm/app/internal/file"
	"github.com/murlokswarm/app/internal/html"
)

func init() {
	app.Loggers = []app.Logger{
		app.NewLogger(os.Stdout, os.Stderr, true, true),
	}
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f *app.Factory) error {
	d.factory = f

	if len(d.NotFoundURL) == 0 {
		d.NotFoundURL = "web.NotFound"
	}

	if d.Server == nil {
		d.Server = &http.Server{
			Addr: ":7042",
		}
	}

	http.Handle("/", d)

	fileHandler := http.FileServer(http.Dir("resources"))
	fileHandler = http.StripPrefix("/resources/", fileHandler)
	fileHandler = newGzipHandler(fileHandler)
	http.Handle("/resources/", fileHandler)

	if d.OnServerRun != nil {
		d.OnServerRun()
	}

	errC := make(chan error)
	go func() {
		err := d.Server.ListenAndServe()
		if err == http.ErrServerClosed {
			err = nil
		}
		errC <- err
	}()

	return <-errC
}

// ServeHTTP is the http.Handler that route wether to serve a page or a
// resource.
func (d *Driver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" || len(req.URL.Path) == 0 {
		req.URL.Path = d.URL
	}

	if n := core.CompoNameFromURL(req.URL); d.factory.IsCompoRegistered(n) {
		d.handleCompo(res, req)
		return
	}

	req.URL.Path = d.NotFoundURL
	d.handleNotFound(res, req)
}

func (d *Driver) handleCompo(res http.ResponseWriter, req *http.Request) {
	c, err := d.factory.NewCompo(core.CompoNameFromURL(req.URL))
	if err != nil {
		http.NotFound(res, req)
		return
	}

	htmlConf := app.HTMLConfig{}
	if configurator, ok := c.(app.Configurator); ok {
		htmlConf = configurator.Config()
	}

	if len(htmlConf.CSS) == 0 {
		htmlConf.CSS = file.CSS(d.Resources("css"))
	}

	htmlConf.Javascripts = append(htmlConf.Javascripts, d.Resources("goapp.js"))
	page := dom.Page(htmlConf, "console.log")

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(page))
}

func (d *Driver) handleNotFound(res http.ResponseWriter, req *http.Request) {
	compo, err := d.factory.NewCompo(core.CompoNameFromURL(req.URL))
	if err != nil {
		http.NotFound(res, req)
		return
	}

	markup := html.NewMarkup(d.factory)

	var root app.Tag
	if root, err = markup.Mount(compo); err != nil {
		http.NotFound(res, req)
		return
	}

	var b bytes.Buffer
	enc := html.NewEncoder(&b, markup, false)

	if err = enc.Encode(root); err != nil {
		http.NotFound(res, req)
		return
	}

	page := html.NewPage(html.PageConfig{
		Title:        "not found",
		DefaultCompo: template.HTML(b.String()),
	})

	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte(page))
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	return "go webapp"
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(p ...string) string {
	return ""
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	app.Log("CallOnUIGoroutine is not supported on server side")
}

// Stop shutdown the server.
func (d *Driver) Stop() {
	d.Server.Shutdown(context.Background())
}
