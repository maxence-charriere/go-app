// +build windows

// Package win is the driver to be used for applications that will run on
// Windows.
package win

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/logs"
)

var (
	driver        *Driver
	goappBuild    = os.Getenv("GOAPP_BUILD")
	debug         = os.Getenv("GOAPP_DEBUG") == "true"
	goappLogsAddr = os.Getenv("GOAPP_LOGS_ADDR")
	goappLogs     *logs.GoappClient
)

func init() {
	if len(goappBuild) != 0 {
		app.Logger = func(format string, a ...interface{}) {}
		return
	}

	if len(goappLogsAddr) != 0 {
		app.EnableDebug(debug)
		goappLogs = logs.NewGoappClient(goappLogsAddr, logs.WithPrompt)
		app.Logger = goappLogs.Logger()
		return
	}

	logger := logs.ToWriter(os.Stderr)
	app.Logger = logs.WithPrompt(logger)
}

// Driver is the app.Driver implementation for Windows.
type Driver struct {
	core.Driver

	// The URL of the component to load in the main window.
	URL string

	// The settings used to generate the app package.
	Settings Settings
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f *app.Factory) error {
	if len(goappBuild) != 0 {
		return d.runGoappBuild()
	}

	app.Log("hello world")
	time.Sleep(time.Second * 5)
	return nil
}

func (d *Driver) runGoappBuild() error {
	b, err := json.MarshalIndent(d.Settings, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(goappBuild, b, 0777)
}
