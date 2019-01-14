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
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/win/uwp"
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
				Timeout: time.Millisecond * 100,
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
		return d.build()
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

	d.Elems = core.NewElemDB()
	d.Events = c.Events
	d.Factory = c.Factory
	d.Platform, d.Go = uwp.RPC(d.UI)
	d.JSToPlatform = "window.external.notify"
	d.OpenDefaultBrowser = openDefaultBrowser
	d.NewWindowFunc = newWindow
	d.ResourcesFunc = d.resources
	// d.StorageFunc = d.storage
	d.UIChan = c.UI
	driver = d

	disconnect := uwp.Connect()
	defer disconnect()

	d.Go.Handle("driver.OnRun", d.onRun)
	d.Go.Handle("driver.OnFilesOpen", d.onFilesOpen)
	d.Go.Handle("driver.OnURLOpen", d.onURLOpen)
	d.Go.Handle("driver.OnClose", d.onClose)
	d.Go.Handle("driver.Log", d.log)

	d.Go.Handle("windows.OnResize", handleWindow(onWindowResize))
	d.Go.Handle("windows.OnFocus", handleWindow(onWindowFocus))
	d.Go.Handle("windows.OnBlur", handleWindow(onWindowBlur))
	d.Go.Handle("windows.OnFullScreen", handleWindow(onWindowFullScreen))
	d.Go.Handle("windows.OnExitFullScreen", handleWindow(onWindowExitFullScreen))
	d.Go.Handle("windows.OnClose", handleWindow(onWindowClose))
	d.Go.Handle("windows.OnCallback", handleWindow(onWindowCallback))
	d.Go.Handle("windows.OnNavigate", handleWindow(onWindowNavigate))

	d.Go.Handle("menus.OnClose", handleMenu(onMenuClose))
	d.Go.Handle("menus.OnCallback", handleMenu(onMenuCallback))

	ctx, cancel := context.WithCancel(context.Background())
	d.stop = cancel
	defer cancel()

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

			case fn := <-d.UIChan:
				fn()
			}
		}
	}()

	wg.Wait()
	return nil
}

func (d *Driver) build() error {
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(goappBuild, b, 0777)
}

func (d *Driver) configureDefaultWindow() {
	if d.DefaultWindow == (app.WindowConfig{}) {
		d.DefaultWindow = app.WindowConfig{
			Title: d.AppName(),
			URL:   d.URL,
		}
	}

	if len(d.DefaultWindow.URL) == 0 {
		d.DefaultWindow.URL = d.URL
	}
}

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) app.Menu {
	m := newMenu(c, "context menu")
	if m.Err() != nil {
		return m
	}

	err := d.Platform.Call("driver.SetContextMenu", nil, struct {
		ID string
	}{
		ID: m.ID(),
	})

	m.SetErr(err)
	return m
}

func (d *Driver) log(in map[string]interface{}) {
	msg := in["Msg"].(string)
	app.Log(msg)
}

func (d *Driver) resources() string {
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		app.Log(err)
	}

	return filepath.Join(appdir, "Resources")
}

func (d *Driver) onRun(in map[string]interface{}) {
	d.configureDefaultWindow()

	if len(d.URL) != 0 {
		app.NewWindow(d.DefaultWindow)
	}

	d.Events.Emit(app.Running)
}

func (d *Driver) onFilesOpen(in map[string]interface{}) {
	d.Events.Emit(app.OpenFilesRequested, core.ConvertToStringSlice(in["Filenames"]))
}

func (d *Driver) onURLOpen(in map[string]interface{}) {
	if u, err := url.Parse(in["URL"].(string)); err == nil {
		d.Events.Emit(app.OpenURLRequested, u)
	}
}

func (d *Driver) onClose(in map[string]interface{}) {
	d.Events.Emit(app.Closed)

	d.UI(func() {
		d.stop()
	})
}

func openDefaultBrowser(url string) error {
	return exec.Command("start", url).Run()
}
