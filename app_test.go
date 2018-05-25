package app_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
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
	ctx, cancel := context.WithCancel(context.Background())

	onRun := func() {
		rd := app.RunningDriver()
		require.NotNil(t, rd, "driver not set")
		assert.Equal(t, "Driver unit tests", app.Name())
		assert.Equal(t, filepath.Join("resources", "hello", "world"), app.Resources("hello", "world"))
		assert.Equal(t, filepath.Join("storage", "hello", "world"), app.Storage("hello", "world"))

		testWindow(t)
		testPage(t, newPage)
		testMenu(t)

		err := app.NewFilePanel(app.FilePanelConfig{})
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
		t.Run("css no resources", testCSSResourcesNoResources)

		app.CallOnUIGoroutine(func() {
		})

		cancel()
	}

	dtest := &test.Driver{
		Ctx:   ctx,
		OnRun: onRun,
	}
	d = app.Logs()(dtest)

	app.Import(&tests.Foo{})
	app.Import(&tests.Bar{})

	newPage = func(c app.PageConfig) (app.Page, error) {
		err := app.NewPage(c)
		if err != nil {
			return nil, err
		}
		return dtest.Page, nil
	}

	if err := app.Run(d); err != nil {
		t.Fatal(err)
	}
}

func testWindow(t *testing.T) {
	win, err := app.NewWindow(app.WindowConfig{
		DefaultURL: "tests.foo",
	})
	if err != nil {
		t.Fatal(err)
	}

	compo := win.Component()
	if compo == nil {
		t.Error("component is nil")
	}

	app.Render(compo)

	var win2 app.Window
	if win2, err = app.WindowByComponent(compo); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, win.ID(), win2.ID())

	if _, err = app.NavigatorByComponent(compo); err != nil {
		t.Fatal(err)
	}

	if _, err = app.WindowByComponent(&tests.Foo{}); err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)

	if _, err = app.NavigatorByComponent(&tests.Foo{}); err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testPage(t *testing.T, newPage func(c app.PageConfig) (app.Page, error)) {
	page, err := newPage(app.PageConfig{
		DefaultURL: "tests.foo",
	})
	if err != nil {
		t.Fatal(err)
	}

	compo := page.Component()
	if compo == nil {
		t.Fatal("component is nil")
	}

	app.Render(compo)

	var page2 app.Page
	if page2, err = app.PageByComponent(compo); err != nil {
		t.Fatal(err)
	}

	if page != page2 {
		t.Fatal("page and page2 are different")
	}

	if _, err = app.NavigatorByComponent(compo); err != nil {
		t.Fatal(err)
	}

	if _, err = newPage(app.PageConfig{
		DefaultURL: "/ErrorTest",
	}); err == nil {
		t.Error("error is nil")
	}
	t.Log(err)

	if _, err = app.PageByComponent(&tests.Foo{}); err == nil {
		t.Error("error is nil")
	}
	t.Log(err)
}

func testMenu(t *testing.T) {
	menu, err := app.NewContextMenu(app.MenuConfig{
		DefaultURL: "tests.bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	compo := menu.Component()
	if compo == nil {
		t.Fatal("component is nil")
	}

	if _, err = app.ElementByComponent(compo); err != nil {
		t.Fatal(err)
	}

	if _, err = app.NavigatorByComponent(compo); err == nil {
		t.Fatal("error is nil")
	}

	if _, err = app.WindowByComponent(compo); err == nil {
		t.Fatal("error is nil")
	}

	if _, err = app.PageByComponent(compo); err == nil {
		t.Fatal("error is nil")
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

func testCSSResourcesNoResources(t *testing.T) {
	if len(app.CSSResources()) != 0 {
		t.Error("resources found")
	}
}
