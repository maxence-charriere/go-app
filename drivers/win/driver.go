// +build windows

// Package win is the driver to be used for applications that will run on
// Windows.
package win

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/file"
	"github.com/pkg/errors"
)

var (
	driver     *Driver
	debug      string
	dev        string
	logsURL    string
	logsWriter *file.HTTPWriter
	logsCancel func()
	goappBuild = os.Getenv("GOAPP_BUILD")
)

func init() {
	runtime.LockOSThread()

	if len(goappBuild) != 0 {
		app.Logger = func(format string, a ...interface{}) {}
		return
	}

	if len(logsURL) != 0 {
		app.EnableDebug(len(debug) != 0)

		logsWriter = &file.HTTPWriter{
			URL: logsURL,
			Client: &http.Client{
				Timeout: time.Second,
			},
		}

		cancel, err := file.CaptureOutput(logsWriter)
		if err != nil {
			panic(err)
		}

		logsCancel = cancel
	}

	logger := core.ToWriter(os.Stderr)
	app.Logger = core.WithPrompt(logger)
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

	// The func called when the app is about to exit.
	OnExit func()

	factory *app.Factory
	elems   *core.ElemDB
	winRPC  bridge.PlatformRPC
	goRPC   bridge.GoRPC
	uichan  chan func()
	stop    func()
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f *app.Factory) error {
	if len(goappBuild) != 0 {
		return d.runGoappBuild()
	}

	if len(dev) != 0 {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("press 'Enter' to exit")
				b := make([]byte, 1)
				os.Stdin.Read(b)
			}
		}()
	}

	if err := loadDLL(); err != nil {
		return errors.Wrap(err, "loading goapp.dll failed")
	}
	defer closeDLL()

	d.factory = f
	d.elems = core.NewElemDB()
	d.winRPC.Handler = winCall

	d.goRPC.Handle("driver.OnRun", d.onRun)
	d.goRPC.Handle("driver.OnExit", d.onExit)
	d.goRPC.Handle("driver.Log", d.log)

	d.goRPC.Handle("windows.OnResize", handleWindow(onWindowResize))
	d.goRPC.Handle("windows.OnFocus", handleWindow(onWindowFocus))
	d.goRPC.Handle("windows.OnBlur", handleWindow(onWindowBlur))
	d.goRPC.Handle("windows.OnFullScreen", handleWindow(onWindowFullScreen))
	d.goRPC.Handle("windows.OnExitFullScreen", handleWindow(onWindowExitFullScreen))
	d.goRPC.Handle("windows.OnCallback", handleWindow(onWindowCallback))

	d.uichan = make(chan func(), 256)
	defer close(d.uichan)

	aliveTicker := time.NewTicker(time.Minute)
	defer aliveTicker.Stop()

	driver = d

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.stop = cancel

	if err := initBridge(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			if logsCancel != nil {
				time.Sleep(time.Second)
				logsCancel()
				logsWriter.Close()
			}

			return nil

		case fn := <-d.uichan:
			fn()

		case <-aliveTicker.C:
		}
	}
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(path ...string) string {
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		app.Log(err)
	}

	r := filepath.Join(path...)
	r = filepath.Join(appdir, "Resources", r)
	return r
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

func (d *Driver) log(in map[string]interface{}) interface{} {
	msg := in["Msg"].(string)
	app.Log(msg)
	return nil
}

func (d *Driver) onRun(in map[string]interface{}) interface{} {
	if d.OnRun == nil {
		d.OnRun = d.newMainWindow
	}

	d.OnRun()
	return nil
}

func (d *Driver) onExit(in map[string]interface{}) interface{} {
	if d.OnExit != nil {
		d.OnExit()
	}

	d.stop()
	return nil
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
