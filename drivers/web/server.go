// +build !js

package web

import (
	"context"
	"net/http"
	"os"
	"runtime"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/murlokswarm/app/internal/file"
	"github.com/murlokswarm/app/internal/logs"
)

func init() {
	logger := logs.ToWriter(os.Stderr)

	if runtime.GOOS == "windows" {
		logger = logs.WithPrompt(logger)
	} else {
		logger = logs.WithColoredPrompt(logger)
	}

	app.Logger = logger
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f *app.Factory) error {
	d.factory = f

	if len(d.NotFoundURL) == 0 {
		d.NotFoundURL = "/web.NotFound"
	}

	if d.Server == nil {
		d.Server = &http.Server{
			Addr: ":7042",
		}
	}

	if addr := os.Getenv("GOAPP_SERVER_ADDR"); len(addr) != 0 {
		d.Server.Addr = addr
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

	if n := core.CompoNameFromURL(req.URL); !d.factory.IsCompoRegistered(n) {
		req.URL.Path = d.NotFoundURL
		d.handle(res, req, http.StatusNotFound)
		return
	}

	d.handle(res, req, http.StatusOK)
}

func (d *Driver) handle(res http.ResponseWriter, req *http.Request, status int) {
	app.WhenDebug(func() {
		app.Logf("serving %s", req.URL)
	})

	compoName := core.CompoNameFromURL(req.URL)

	c, err := d.factory.NewCompo(compoName)
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
	page := dom.Page(htmlConf, "console.log", compoName)

	res.WriteHeader(status)
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
	app.Logf("CallOnUIGoroutine is not supported on server side")
}

// Stop shutdown the server.
func (d *Driver) Stop() {
	d.Server.Shutdown(context.Background())
}
