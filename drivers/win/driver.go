// +build windows

// Package win is the driver to be used for applications that will run on
// Windows.
package win

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/logs"
	"github.com/pkg/errors"
)

var (
	driver        *Driver
	goappBuild    = os.Getenv("GOAPP_BUILD")
	debug         = os.Getenv("GOAPP_DEBUG") == "true"
	goappLogsAddr = os.Getenv("GOAPP_LOGS_ADDR")
	goappLogs     *logs.GoappClient
)

func init() {
	runtime.LockOSThread()

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

	// The func called right after app.Run.
	OnRun func()

	factory  *app.Factory
	elems    *core.ElemDB
	winRPC   bridge.PlatformRPC
	goRPC    bridge.GoRPC
	uichan   chan func()
	stopchan chan error
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f *app.Factory) error {
	if len(goappBuild) != 0 {
		return d.runGoappBuild()
	}

	defer func() {
		recover()
		time.Sleep(time.Second * 30)
	}()

	if err := loadDLL(); err != nil {
		return errors.Wrap(err, "loading goapp.dll failed")
	}
	defer closeDLL()

	d.factory = f
	d.elems = core.NewElemDB()
	d.winRPC.Handler = winCall

	d.uichan = make(chan func(), 256)
	defer close(d.uichan)

	d.uichan = make(chan func(), 256)
	d.stopchan = make(chan error)
	aliveTicker := time.NewTicker(time.Minute)
	defer close(d.uichan)
	defer close(d.stopchan)
	defer aliveTicker.Stop()

	driver = d

	if err := initBridge(); err != nil {
		return err
	}

	d.onRun()

	for {
		select {
		case err := <-d.stopchan:
			return err

		case fn := <-d.uichan:
			fn()

		case <-aliveTicker.C:
		}
	}
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Compo) {
	e := d.ElemByCompo(c)
	if e.Err() == nil {
		e.(app.ElemWithCompo).Render(c)
	}
}

// ElemByCompo satisfies the app.Driver interface.
func (d *Driver) ElemByCompo(c app.Compo) app.Elem {
	return d.elems.GetByCompo(c)
}

// NewWindow satisfies the app.Driver interface.
func (d *Driver) NewWindow(c app.WindowConfig) app.Window {
	return newWindow(c)
}

// CallOnUIGoroutine satisfies the app.Driver interface.
func (d *Driver) CallOnUIGoroutine(f func()) {
	d.uichan <- f
}

func (d *Driver) runGoappBuild() error {
	b, err := json.MarshalIndent(d.Settings, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(goappBuild, b, 0777)
}

func (d *Driver) onRun() {
	if d.OnRun == nil {
		d.OnRun = d.newMainWindow
	}

	d.OnRun()
}

func (d *Driver) newMainWindow() {
	app.NewWindow(app.WindowConfig{
		Title:          d.AppName(),
		TitlebarHidden: true,
		MinWidth:       480,
		Width:          1280,
		MinHeight:      480,
		Height:         768,
		URL:            d.URL,
	})
}
