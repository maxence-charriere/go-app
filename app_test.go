package app_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/tests"
)

func TestImport(t *testing.T) {
	app.Import(&tests.Foo{})

	defer func() { recover() }()
	app.Import(tests.NoPointerCompo{})
}

func TestApp(t *testing.T) {
	var d app.Driver
	ctx, cancel := context.WithCancel(context.Background())

	onRun := func() {
		if rd := app.RunningDriver(); rd == nil {
			t.Fatal("driver is not set")
		}

		if name := app.Name(); name != "Driver unit tests" {
			t.Error("app name is not test:", name)
		}

		if resources := app.Resources("hello", "world"); resources != "resources/hello/world" {
			t.Error("resources is not resources/hello/world:", resources)
		}

		if storage := app.Storage("hello", "world"); storage != "storage/hello/world" {
			t.Error("storage is not storage/hello/world:", storage)
		}

		win, err := app.NewWindow(app.WindowConfig{
			DefaultURL: "tests.foo",
		})
		if err != nil {
			t.Error(err)
		}

		winCompo := win.Component()
		if winCompo == nil {
			t.Error("component is nil")
		}

		app.Render(winCompo)

		var win2 app.Window
		if win2, err = app.WindowByComponent(winCompo); err != nil {
			t.Error(err)
		}

		if win != win2 {
			t.Error("win and win2 are different")
		}

		if _, err = app.WindowByComponent(&tests.Foo{}); err == nil {
			t.Error("error is nil")
		}

		var page app.Page
		if page, err = app.NewPage(app.PageConfig{
			DefaultURL: "tests.foo",
		}); err != nil {
			t.Error(err)
		}

		pageCompo := page.Component()
		if pageCompo == nil {
			t.Error("component is nil")
		}

		app.Render(pageCompo)

		var page2 app.Page
		if page2, err = app.PageByComponent(pageCompo); err != nil {
			t.Error(err)
		}

		if page != page2 {
			t.Error("page and page2 are different")
		}

		if _, err := app.NewPage(app.PageConfig{
			DefaultURL: "/ErrorTest",
		}); err == nil {
			t.Error("error is nil")
		}

		if _, err = app.PageByComponent(&tests.Foo{}); err == nil {
			t.Error("error is nil")
		}

		if _, err = app.PageByComponent(winCompo); err == nil {
			t.Error("error is nil")
		}

		var menu app.Menu
		if menu, err = app.NewContextMenu(app.MenuConfig{
			DefaultURL: "tests.bar",
		}); err != nil {
			t.Error(err)
		}

		menuCompo := menu.Component()
		if menuCompo == nil {
			t.Error("component is nil")
		}

		if _, err = app.WindowByComponent(menuCompo); err == nil {
			t.Error("error is nil")
		}

		if _, err = app.ElementByComponent(menuCompo); err != nil {
			t.Error(err)
		}

		err = app.NewFilePanel(app.FilePanelConfig{})
		if err != nil && !app.NotSupported(err) {
			t.Error(err)
		}

		err = app.NewSaveFilePanel(app.SaveFilePanelConfig{})
		if err != nil && !app.NotSupported(err) {
			t.Error(err)
		}

		err = app.NewShare("Hello world")
		if err != nil && !app.NotSupported(err) {
			t.Error(err)
		}

		err = app.NewNotification(app.NotificationConfig{})
		if err != nil && !app.NotSupported(err) {
			t.Error(err)
		}

		app.MenuBar()
		app.Dock()

		t.Run("css resources", testCSSResources)

		app.CallOnUIGoroutine(func() {
		})

		cancel()
	}

	d = &test.Driver{
		Ctx:   ctx,
		OnRun: onRun,
	}

	app.Import(&tests.Foo{})
	app.Import(&tests.Bar{})

	if err := app.Run(d); err != nil {
		t.Fatal(err)
	}
}

func testCSSResources(t *testing.T) {
	defer os.RemoveAll(app.Resources())

	os.MkdirAll(app.Resources("css"), 0777)
	if f1, err := os.Create(app.Resources("css", "test.css")); err == nil {
		defer f1.Close()
	}
	if f2, err := os.Create(app.Resources("css", "test.scss")); err == nil {
		defer f2.Close()
	}

	os.MkdirAll(app.Resources("css", "sub"), 0777)
	if f3, err := os.Create(app.Resources("css", "sub", "sub.css")); err == nil {
		defer f3.Close()
	}

	css := app.CSSResources()
	expected := []string{
		app.Resources("css", "sub", "sub.css"),
		app.Resources("css", "test.css"),
	}

	if !reflect.DeepEqual(css, expected) {
		t.Error("expected:", expected)
		t.Error("current :", css)
	}
}
