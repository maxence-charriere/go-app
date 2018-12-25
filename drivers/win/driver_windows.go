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
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
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

		logsCancel = func() {
			cancel()
			logsWriter.Close()
		}
	}

	logger := core.ToWriter(os.Stderr)
	app.Logger = core.WithPrompt(logger)
	app.EnableDebug(len(debug) != 0)
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(c app.DriverConfig) error {
	if len(goappBuild) != 0 {
		return d.runGoappBuild()
	}

	if len(dev) != 0 {
		defer func() {
			err := recover()
			if err != nil {
				app.Log(errors.Errorf("%v", err))
			}

			fmt.Println("press 'Enter' to exit")
			b := make([]byte, 1)
			os.Stdin.Read(b)
		}()
	}

	if err := loadDLL(); err != nil {
		return errors.Wrap(err, "loading goapp.dll failed")
	}
	defer closeDLL()

	d.ui = c.UI
	d.factory = c.Factory
	d.events = c.Events
	d.elems = core.NewElemDB()
	d.winRPC.Handler = winCall
	driver = d

	d.goRPC.Handle("driver.OnRun", d.onRun)
	d.goRPC.Handle("driver.OnFilesOpen", d.onFilesOpen)
	d.goRPC.Handle("driver.OnURLOpen", d.onURLOpen)
	d.goRPC.Handle("driver.OnClose", d.onClose)
	d.goRPC.Handle("driver.Log", d.log)

	d.goRPC.Handle("windows.OnResize", handleWindow(onWindowResize))
	d.goRPC.Handle("windows.OnFocus", handleWindow(onWindowFocus))
	d.goRPC.Handle("windows.OnBlur", handleWindow(onWindowBlur))
	d.goRPC.Handle("windows.OnFullScreen", handleWindow(onWindowFullScreen))
	d.goRPC.Handle("windows.OnExitFullScreen", handleWindow(onWindowExitFullScreen))
	d.goRPC.Handle("windows.OnClose", handleWindow(onWindowClose))
	d.goRPC.Handle("windows.OnCallback", handleWindow(onWindowCallback))
	d.goRPC.Handle("windows.OnNavigate", handleWindow(onWindowNavigate))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.stop = cancel

	if err := initBridge(); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				if logsCancel != nil {
					logsCancel()
				}

				wg.Done()
				return

			case fn := <-d.ui:
				fn()
			}
		}
	}()

	wg.Wait()
	return nil
}

func (d *Driver) runGoappBuild() error {
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(goappBuild, b, 0777)
}

func (d *Driver) configureDefaultWindow() {
	if d.DefaultWindow == (app.WindowConfig{}) {
		d.DefaultWindow = app.WindowConfig{
			Title:     d.AppName(),
			MinWidth:  480,
			Width:     1280,
			MinHeight: 480,
			Height:    768,
			URL:       d.URL,
		}
	}

	if len(d.DefaultWindow.URL) == 0 {
		d.DefaultWindow.URL = d.URL
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
	e.(app.ElemWithCompo).Render(c)
}

// ElemByCompo satisfies the app.Driver interface.
func (d *Driver) ElemByCompo(c app.Compo) app.Elem {
	return d.elems.GetByCompo(c)
}

// NewWindow satisfies the app.Driver interface.
func (d *Driver) NewWindow(c app.WindowConfig) app.Window {
	return newWindow(c)
}

// UI satisfies the app.Driver interface.
func (d *Driver) UI(f func()) {
	d.ui <- f
}

func (d *Driver) log(in map[string]interface{}) interface{} {
	msg := in["Msg"].(string)
	app.Log(msg)
	return nil
}

func (d *Driver) onRun(in map[string]interface{}) interface{} {
	d.configureDefaultWindow()

	if len(d.URL) != 0 {
		app.NewWindow(d.DefaultWindow)
	}

	d.events.Emit(app.Running)
	return nil
}

func (d *Driver) onFilesOpen(in map[string]interface{}) interface{} {
	d.events.Emit(app.OpenFilesRequested, bridge.Strings(in["Filenames"]))
	return nil
}

func (d *Driver) onURLOpen(in map[string]interface{}) interface{} {
	if u, err := url.Parse(in["URL"].(string)); err != nil {
		d.events.Emit(app.OpenURLRequested, u)
	}

	return nil
}

func (d *Driver) onClose(in map[string]interface{}) interface{} {
	d.events.Emit(app.Closed)

	d.UI(func() {
		d.stop()
	})

	return nil
}
