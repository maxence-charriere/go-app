// +build !js

package web

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/murlokswarm/app/internal/file"
)

func init() {
	logger := core.ToWriter(os.Stderr)

	if runtime.GOOS == "windows" {
		logger = core.WithPrompt(logger)
	} else {
		logger = core.WithColoredPrompt(logger)
	}

	app.Logger = logger
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(c app.DriverConfig) error {
	d.factory = c.Factory

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
	app.WhenDebug(func() {
		app.Logf("serving %s", req.URL)
	})

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
		htmlConf.CSS = file.Filenames(d.Resources("css"), ".css")
	}

	if len(htmlConf.Javascripts) == 0 {
		htmlConf.Javascripts = file.Filenames(d.Resources("js"), ".js")
	}

	htmlConf.Javascripts = append(htmlConf.Javascripts, d.Resources("goapp.js"))

	page := dom.Page{
		Title:         htmlConf.Title,
		Metas:         htmlConf.Metas,
		Icon:          d.Icon,
		CSS:           cleanWindowsPath(htmlConf.CSS),
		Javascripts:   cleanWindowsPath(htmlConf.Javascripts),
		GoRequest:     "console.log", // Overloaded in client.go.
		RootCompoName: compoName,
	}

	res.WriteHeader(status)
	res.Write([]byte(page.String()))
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

func cleanWindowsPath(paths []string) []string {
	c := make([]string, len(paths))

	for i, p := range paths {
		c[i] = strings.Replace(p, `\`, "/", -1)
	}

	return c
}
