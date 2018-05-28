package app_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	app.Import(&tests.Foo{})

	defer func() { recover() }()
	app.Import(tests.NoPointerCompo{})
}

func TestApp(t *testing.T) {
	var d app.Driver
	var newPage func(c app.PageConfig) (app.Page, error)

	output := &bytes.Buffer{}
	app.Loggers = []app.Logger{app.NewLogger(output, output, true, true)}

	app.Import(&tests.Foo{})
	app.Import(&tests.Bar{})

	ctx, cancel := context.WithCancel(context.Background())

	onRun := func() {
		rd := app.RunningDriver()
		require.NotNil(t, rd, "driver not set")
		assert.NotEmpty(t, app.Name())
		assert.Equal(t, filepath.Join("resources", "hello", "world"), app.Resources("hello", "world"))
		assert.Equal(t, filepath.Join("storage", "hello", "world"), app.Storage("hello", "world"))

		// Window:
		win, err := app.NewWindow(app.WindowConfig{
			DefaultURL: "tests.foo",
		})
		require.NoError(t, err)

		compo := win.Component()
		require.NotNil(t, compo)
		app.Render(compo)

		var win2 app.Window
		win2, err = app.WindowByComponent(compo)
		require.NoError(t, err)
		assert.Equal(t, win.ID(), win2.ID())

		var nav app.Navigator
		nav, err = app.NavigatorByComponent(compo)
		require.NoError(t, err)
		assert.Equal(t, win.ID(), nav.ID())

		// Page:
		var page app.Page
		page, err = newPage(app.PageConfig{
			DefaultURL: "tests.foo",
		})
		require.NoError(t, err)

		compo = page.Component()
		require.NotNil(t, compo)
		app.Render(compo)

		var page2 app.Page
		page2, err = app.PageByComponent(compo)
		require.NoError(t, err)
		assert.Equal(t, page.ID(), page2.ID())

		nav, err = app.NavigatorByComponent(compo)
		require.NoError(t, err)
		assert.Equal(t, page.ID(), nav.ID())

		// Menu:
		var menu app.Menu
		menu, err = app.NewContextMenu(app.MenuConfig{
			DefaultURL: "tests.bar",
		})
		require.NoError(t, err)

		compo = menu.Component()
		require.NotNil(t, compo)
		app.Render(compo)

		var elem app.ElementWithComponent
		elem, err = app.ElementByComponent(compo)
		require.NoError(t, err)
		assert.Equal(t, menu.ID(), elem.ID())

		_, err = app.NavigatorByComponent(compo)
		assert.Error(t, err)

		_, err = app.WindowByComponent(compo)
		assert.Error(t, err)

		_, err = app.PageByComponent(compo)
		assert.Error(t, err)

		// File panels:
		err = app.NewFilePanel(app.FilePanelConfig{})
		require.NoError(t, err)

		err = app.NewSaveFilePanel(app.SaveFilePanelConfig{})
		require.NoError(t, err)

		// Share:
		err = app.NewShare("Hello world")
		require.NoError(t, err)

		// Notifications:
		err = app.NewNotification(app.NotificationConfig{})
		require.NoError(t, err)

		// Menubar:
		_, err = app.MenuBar()
		require.NoError(t, err)

		// Status menu:
		var statusMenu app.StatusMenu
		statusMenu, err = app.NewStatusMenu(app.StatusMenuConfig{})
		require.NoError(t, err)

		err = statusMenu.Load("tests.bar")
		require.NoError(t, err)

		compo = statusMenu.Component()
		require.NotNil(t, compo)
		app.Render(compo)

		elem, err = app.ElementByComponent(compo)
		require.NoError(t, err)
		assert.Equal(t, statusMenu.ID(), elem.ID())

		err = statusMenu.SetText("test")
		assert.NoError(t, err)

		err = statusMenu.SetIcon(filepath.Join("tests", "resources", "logo.png"))
		assert.NoError(t, err)

		err = statusMenu.Close()
		assert.NoError(t, err)

		// Dock:
		var dockTile app.DockTile
		dockTile, err = app.Dock()
		require.NoError(t, err)

		err = dockTile.Load("tests.bar")
		require.NoError(t, err)

		compo = dockTile.Component()
		require.NotNil(t, compo)
		app.Render(compo)

		elem, err = app.ElementByComponent(compo)
		require.NoError(t, err)
		assert.Equal(t, dockTile.ID(), elem.ID())

		err = dockTile.SetBadge("42")
		assert.NoError(t, err)

		err = dockTile.SetIcon(filepath.Join("tests", "resources", "logo.png"))
		assert.NoError(t, err)

		// CSS resources:
		assert.Len(t, app.CSSResources(), 0)

		os.MkdirAll(app.Resources("css", "sub"), 0777)
		os.Create(app.Resources("css", "test.css"))
		os.Create(app.Resources("css", "test.scss"))
		os.Create(app.Resources("css", "sub", "sub.css"))
		defer os.RemoveAll(app.Resources())

		assert.Contains(t, app.CSSResources(), app.Resources("css", "test.css"))
		assert.NotContains(t, app.CSSResources(), app.Resources("css", "test.scss"))
		assert.Contains(t, app.CSSResources(), app.Resources("css", "sub", "sub.css"))

		app.CallOnUIGoroutine(func() {
			t.Log("CallOnUIGoroutine")
		})

		cancel()
	}

	dtest := &test.Driver{
		Ctx:   ctx,
		OnRun: onRun,
	}
	d = dtest

	newPage = func(c app.PageConfig) (app.Page, error) {
		err := app.NewPage(c)
		if err != nil {
			return nil, err
		}
		return dtest.Page, nil
	}

	err := app.Run(d, app.Logs())
	require.NoError(t, err)

	t.Log(output.String())
}

func TestAppError(t *testing.T) {
	var d app.Driver

	output := &bytes.Buffer{}
	app.Loggers = []app.Logger{app.NewLogger(output, output, true, true)}

	ctx, cancel := context.WithCancel(context.Background())

	onRun := func() {
		defer cancel()

		app.Render(nil)

		// Window:
		win, err := d.NewWindow(app.WindowConfig{})
		require.NoError(t, err)

		err = win.Load("")
		assert.Error(t, err)

		err = win.Render(nil)
		assert.Error(t, err)

		err = win.Reload()
		assert.Error(t, err)

		err = win.Previous()
		assert.Error(t, err)

		err = win.Next()
		assert.Error(t, err)

		err = win.Close()
		assert.Error(t, err)

		err = win.Move(0, 0)
		assert.Error(t, err)

		err = win.Center()
		assert.Error(t, err)

		err = win.Resize(0, 0)
		assert.Error(t, err)

		err = win.Focus()
		assert.Error(t, err)

		err = win.ToggleFullScreen()
		assert.Error(t, err)

		err = win.ToggleMinimize()
		assert.Error(t, err)

		// Menu:
		var menu app.Menu
		menu, err = d.NewContextMenu(app.MenuConfig{})
		require.NoError(t, err)

		err = menu.Load("")
		assert.Error(t, err)

		err = menu.Render(nil)
		assert.Error(t, err)

		// Status Bar:
		var statusMenu app.StatusMenu
		statusMenu, err = app.NewStatusMenu(app.StatusMenuConfig{
			Text: "test",
		})
		require.NoError(t, err)

		err = statusMenu.Load("")
		assert.Error(t, err)

		err = statusMenu.Render(nil)
		assert.Error(t, err)

		err = statusMenu.SetText("")
		assert.Error(t, err)

		err = statusMenu.SetIcon("")
		assert.Error(t, err)

		err = statusMenu.Close()
		assert.Error(t, err)

		// Dock tile:
		var dockTile app.DockTile
		dockTile, err = app.Dock()
		require.NoError(t, err)

		err = dockTile.Load("")
		assert.Error(t, err)

		err = dockTile.Render(nil)
		assert.Error(t, err)

		err = dockTile.SetIcon("")
		assert.Error(t, err)

		err = dockTile.SetBadge("")
		assert.Error(t, err)
	}

	dtest := &test.Driver{
		Ctx:             ctx,
		SimulateErr:     true,
		SimulateElemErr: true,
		OnRun:           onRun,
	}
	d = app.Logs()(dtest)

	err := app.Run(d, app.Logs())
	assert.Error(t, err)

	_, err = app.NewWindow(app.WindowConfig{})
	assert.Error(t, err)

	err = app.NewPage(app.PageConfig{})
	assert.Error(t, err)

	_, err = app.NewContextMenu(app.MenuConfig{})
	assert.Error(t, err)

	_, err = app.ElementByComponent(nil)
	assert.Error(t, err)

	_, err = app.NavigatorByComponent(nil)
	assert.Error(t, err)

	_, err = app.WindowByComponent(nil)
	assert.Error(t, err)

	_, err = app.PageByComponent(nil)
	assert.Error(t, err)

	err = app.NewFilePanel(app.FilePanelConfig{})
	assert.Error(t, err)

	err = app.NewSaveFilePanel(app.SaveFilePanelConfig{})
	assert.Error(t, err)

	err = app.NewShare(nil)
	assert.Error(t, err)

	err = app.NewNotification(app.NotificationConfig{})
	assert.Error(t, err)

	_, err = app.MenuBar()
	assert.Error(t, err)

	_, err = app.NewStatusMenu(app.StatusMenuConfig{})
	assert.Error(t, err)

	_, err = app.Dock()
	assert.Error(t, err)

	dtest.SimulateErr = false
	err = app.Run(d, app.Logs())
	require.NoError(t, err)

	t.Log(output.String())
}
