// +build darwin,amd64

package mac

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/mac/objc"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
	"github.com/pkg/errors"
)

var (
	driver     *Driver
	goappBuild = os.Getenv("GOAPP_BUILD")
	debug      = os.Getenv("GOAPP_DEBUG") == "true"
)

const (
	// PreferencesRequested is the event emitted when the menubar Preferences
	// button is clicked.
	PreferencesRequested app.Event = "app.mac.preferencesRequested"
)

func init() {
	if len(goappBuild) != 0 {
		app.Logger = func(format string, a ...interface{}) {}
		return
	}

	logger := core.ToWriter(os.Stderr)
	app.Logger = core.WithColoredPrompt(logger)
	app.EnableDebug(debug)
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(c app.DriverConfig) error {
	if len(goappBuild) != 0 {
		return d.runGoappBuild()
	}

	if driver != nil {
		return errors.New("running already")
	}

	d.ui = c.UI
	d.factory = c.Factory
	d.events = c.Events
	d.elems = core.NewElemDB()
	d.devID = generateDevID()
	d.macRPC, d.goRPC = objc.RPC(d.UI)
	driver = d

	d.goRPC.Handle("driver.OnRun", d.onRun)
	d.goRPC.Handle("driver.OnFocus", d.onFocus)
	d.goRPC.Handle("driver.OnBlur", d.onBlur)
	d.goRPC.Handle("driver.OnReopen", d.onReopen)
	d.goRPC.Handle("driver.OnFilesOpen", d.onFilesOpen)
	d.goRPC.Handle("driver.OnURLOpen", d.onURLOpen)
	d.goRPC.Handle("driver.OnFileDrop", d.onFileDrop)
	d.goRPC.Handle("driver.OnClose", d.onClose)

	d.goRPC.Handle("windows.OnMove", handleWindow(onWindowMove))
	d.goRPC.Handle("windows.OnResize", handleWindow(onWindowResize))
	d.goRPC.Handle("windows.OnFocus", handleWindow(onWindowFocus))
	d.goRPC.Handle("windows.OnBlur", handleWindow(onWindowBlur))
	d.goRPC.Handle("windows.OnFullScreen", handleWindow(onWindowFullScreen))
	d.goRPC.Handle("windows.OnExitFullScreen", handleWindow(onWindowExitFullScreen))
	d.goRPC.Handle("windows.OnMinimize", handleWindow(onWindowMinimize))
	d.goRPC.Handle("windows.OnDeminimize", handleWindow(onWindowDeminimize))
	d.goRPC.Handle("windows.OnClose", handleWindow(onWindowClose))
	d.goRPC.Handle("windows.OnCallback", handleWindow(onWindowCallback))
	d.goRPC.Handle("windows.OnNavigate", handleWindow(onWindowNavigate))
	d.goRPC.Handle("windows.OnAlert", handleWindow(onWindowAlert))

	d.goRPC.Handle("menus.OnClose", handleMenu(onMenuClose))
	d.goRPC.Handle("menus.OnCallback", handleMenu(onMenuCallback))

	d.goRPC.Handle("controller.OnDirectionChange", handleController(onControllerDirectionChange))
	d.goRPC.Handle("controller.OnButtonPressed", handleController(onControllerButtonPressed))
	d.goRPC.Handle("controller.OnConnected", handleController(onControllerConnected))
	d.goRPC.Handle("controller.OnDisconnected", handleController(onControllerDisconnected))
	d.goRPC.Handle("controller.OnPause", handleController(onControllerPause))
	d.goRPC.Handle("controller.OnClose", handleController(onControllerClose))

	d.goRPC.Handle("filePanels.OnSelect", handleFilePanel(onFilePanelSelect))
	d.goRPC.Handle("saveFilePanels.OnSelect", handleSaveFilePanel(onSaveFilePanelSelect))

	d.goRPC.Handle("notifications.OnReply", handleNotification(onNotificationReply))

	ctx, cancel := context.WithCancel(context.Background())
	d.stop = cancel

	go func() {
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				d.macRPC.Call("driver.Terminate", nil, nil)
				return

			case fn := <-d.ui:
				fn()
			}
		}
	}()

	err := d.macRPC.Call("driver.Run", nil, nil)
	return err
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
			MinHeight: 480,
			URL:       d.URL,
		}
	}

	if len(d.DefaultWindow.URL) == 0 {
		d.DefaultWindow.URL = d.URL
	}
}

// AppName satisfies the app.Driver interface.
func (d *Driver) AppName() string {
	out := struct {
		AppName string
	}{}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		app.Panic(err)
	}

	if len(out.AppName) != 0 {
		return out.AppName
	}

	wd, err := os.Getwd()
	if err != nil {
		app.Panic(errors.Wrap(err, "app name unreachable"))
	}

	return filepath.Base(wd)
}

// Resources satisfies the app.Driver interface.
func (d *Driver) Resources(path ...string) string {
	out := struct {
		Resources string
	}{}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		app.Panic(err)
	}

	r := filepath.Join(path...)
	return filepath.Join(out.Resources, r)
}

// Storage satisfies the app.Driver interface.
func (d *Driver) Storage(path ...string) string {
	s := filepath.Join(path...)
	return filepath.Join(d.support(), "storage", s)
}

// Render satisfies the app.Driver interface.
func (d *Driver) Render(c app.Compo) {
	e := d.ElemByCompo(c)

	if e.Err() == app.ErrElemNotSet {
		return
	}

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

// NewContextMenu satisfies the app.Driver interface.
func (d *Driver) NewContextMenu(c app.MenuConfig) app.Menu {
	m := newMenu(c, "context menu")
	if m.Err() != nil {
		return m
	}

	err := d.macRPC.Call("driver.SetContextMenu", nil, m.ID())
	m.SetErr(err)
	return m
}

// NewController statisfies the app.Driver interface.
func (d *Driver) NewController(c app.ControllerConfig) app.Controller {
	return newController(c)
}

// NewFilePanel satisfies the app.Driver interface.
func (d *Driver) NewFilePanel(c app.FilePanelConfig) app.Elem {
	return newFilePanel(c)
}

// NewSaveFilePanel satisfies the app.Driver interface.
func (d *Driver) NewSaveFilePanel(c app.SaveFilePanelConfig) app.Elem {
	return newSaveFilePanel(c)
}

// NewShare satisfies the app.Driver interface.
func (d *Driver) NewShare(v interface{}) app.Elem {
	return newSharePanel(v)
}

// NewNotification satisfies the app.Driver interface.
func (d *Driver) NewNotification(c app.NotificationConfig) app.Elem {
	return newNotification(c)
}

// MenuBar satisfies the app.Driver interface.
func (d *Driver) MenuBar() app.Menu {
	return d.menubar
}

// NewStatusMenu satisfies the app.Driver interface.
func (d *Driver) NewStatusMenu(c app.StatusMenuConfig) app.StatusMenu {
	return newStatusMenu(c)
}

// DockTile satisfies the app.Driver interface.
func (d *Driver) DockTile() app.DockTile {
	return d.docktile
}

// UI satisfies the app.Driver interface.
func (d *Driver) UI(f func()) {
	d.ui <- f
}

// Stop satisfies the app.Driver interface.
func (d *Driver) Stop() {
	if err := d.macRPC.Call("driver.Close", nil, nil); err != nil {
		app.Log("stop failed:", err)
		d.stop()
	}
}

func (d *Driver) support() string {
	out := struct {
		Support string
	}{}

	if err := d.macRPC.Call("driver.Bundle", &out, nil); err != nil {
		app.Panic(err)
	}

	// Set up the support directory in case of the app is not bundled.
	if strings.HasSuffix(out.Support, "{appname}") {
		wd, err := os.Getwd()
		if err != nil {
			app.Panic(errors.Wrap(err, "support unreachable"))
		}

		appname := filepath.Base(wd) + "-" + d.devID
		out.Support = strings.Replace(out.Support, "{appname}", appname, 1)
	}

	return out.Support
}

func (d *Driver) onRun(in map[string]interface{}) interface{} {
	d.configureDefaultWindow()
	d.menubar = newMenuBar(d.MenubarConfig)
	d.docktile = newDockTile(app.MenuConfig{URL: d.DockURL})

	if len(d.URL) != 0 {
		app.NewWindow(d.DefaultWindow)
	}

	d.events.Emit(app.Running)
	return nil
}

func (d *Driver) onFocus(in map[string]interface{}) interface{} {
	d.events.Emit(app.Focused)
	return nil
}

func (d *Driver) onBlur(in map[string]interface{}) interface{} {
	d.events.Emit(app.Blurred)
	return nil
}

func (d *Driver) onReopen(in map[string]interface{}) interface{} {
	hasVisibleWindow := in["HasVisibleWindows"].(bool)

	if !hasVisibleWindow && len(d.URL) != 0 {
		app.NewWindow(d.DefaultWindow)
	}

	d.events.Emit(app.Reopened)
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

func (d *Driver) onFileDrop(in map[string]interface{}) interface{} {
	d.droppedFiles = bridge.Strings(in["Filenames"])
	return nil
}

func (d *Driver) onClose(in map[string]interface{}) interface{} {
	d.events.Emit(app.Closed)

	d.UI(func() {
		d.stop()
	})

	return nil
}

func generateDevID() string {
	h := md5.New()
	wd, _ := os.Getwd()
	io.WriteString(h, wd)
	return fmt.Sprintf("%x", h.Sum(nil))
}
