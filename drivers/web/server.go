// +build !js

package web

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"path"

	"github.com/murlokswarm/app/appjs"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f app.Factory) error {
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
		req.URL.Path = d.DefaultURL
	}

	if compoName := app.ComponentNameFromURL(req.URL); d.factory.Registered(compoName) {
		d.handleComponent(res, req)
		return
	}

	req.URL.Path = d.NotFoundURL
	d.handleNotFound(res, req)
}

func (d *Driver) handleComponent(res http.ResponseWriter, req *http.Request) {
	compo, err := d.factory.New(app.ComponentNameFromURL(req.URL))
	if err != nil {
		http.NotFound(res, req)
		return
	}

	var config html.PageConfig
	if page, ok := compo.(html.Page); ok {
		config = page.PageConfig()
	}

	if len(config.CSS) == 0 {
		config.CSS = app.CSSResources()
	}

	config.Javascripts = append(config.Javascripts, d.Resources("goapp.js"))
	config.AppJS = appjs.AppJS("console.log")

	page := html.NewPage(config)

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(page))
}

func (d *Driver) handleNotFound(res http.ResponseWriter, req *http.Request) {
	compo, err := d.factory.New(app.ComponentNameFromURL(req.URL))
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
		Title:            "not found",
		DefaultComponent: template.HTML(b.String()),
	})

	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte(page))
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	return "goapp"
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(p ...string) string {
	resources := path.Join(p...)
	resources = path.Join("resources", resources)
	return resources
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(p ...string) string {
	return ""
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Component) error {
	return app.NewErrNotSupported("render")
}

// ElementByComponent satisfies the app.Driver interface.
func (d *Driver) ElementByComponent(c app.Component) (app.ElementWithComponent, error) {
	return nil, app.NewErrNotSupported("element by component")
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	app.Error("CallOnUIGoroutine is not supported on server side")
}

// Close shutdown the server.
func (d *Driver) Close() {
	d.Server.Shutdown(context.Background())
}
