// +build !js

package web

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/murlokswarm/app/appjs"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "Web"
}

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

	d.fileHandler = http.FileServer(http.Dir("resources"))
	d.Server.Handler = d
	// http.Handle("/", d)

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
	if len(req.URL.Path) == 0 || req.URL.Path == "/" {
		req.URL.Path = d.DefaultURL
	}

	if compoName := app.ComponentNameFromURL(req.URL); d.factory.Registered(compoName) {
		d.handleComponent(res, req)
		return
	}

	filename := filepath.Join("resources", req.URL.Path)
	if fi, err := os.Stat(filename); err != nil || fi.IsDir() {
		req.URL.Path = d.NotFoundURL
		d.handleNotFound(res, req)
		return
	}

	d.fileHandler.ServeHTTP(res, req)
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

	config.Javascripts = append(config.Javascripts, "goapp.js")
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
	enc := html.NewEncoder(&b, markup)

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
	panic("not implemented")
}

// Close shutdown the server.
func (d *Driver) Close() {
	d.Server.Shutdown(context.Background())
}
