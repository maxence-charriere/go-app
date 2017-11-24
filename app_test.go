package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
)

type Component app.ZeroCompo

func (c *Component) Render() string {
	return `<div>Hello</div>`
}

type InvalidComponent app.ZeroCompo

func (c InvalidComponent) Render() string {
	return ``
}

func TestApp(t *testing.T) {
	driver := &test.Driver{
		Test: t,
	}

	tests := []struct {
		scenario string
		function func(t *testing.T)
	}{
		{
			scenario: "imports a component",
			function: testImport,
		},
		{
			scenario: "imports invalid component panics",
			function: testImportInvalidComponent,
		},
		{
			scenario: "runs faulty driver panics",
			function: testRunDriverError,
		},
		{
			scenario: "get running driver when app is not running panics",
			function: testRunningDriverPanic,
		},
		{
			scenario: "run",
			function: func(t *testing.T) { testRun(t, driver) },
		},
		{
			scenario: "call run while app is running panics",
			function: testRunMultiple,
		},
		{
			scenario: "import component while app is running panics",
			function: testImportWhenDriverRuns,
		},
		{
			scenario: "returns the running driver",
			function: func(t *testing.T) { testRunningDriver(t, driver) },
		},
		{
			scenario: "renders a component",
			function: func(t *testing.T) { testRender(t, driver) },
		},
		{
			scenario: "context returns an element",
			function: func(t *testing.T) { testContext(t, driver) },
		},
		{
			scenario: "context returns an error",
			function: testContextError,
		},
		{
			scenario: "creates a context menu",
			function: testNewContextMenu,
		},
		{
			scenario: "resources returns a filepath",
			function: testResources,
		},
		{
			scenario: "call on ui goroutine",
			function: func(t *testing.T) { testCallOnUIGoroutine(t, driver) },
		},
		{
			scenario: "storage returns a filepath",
			function: testStorage,
		},
		{
			scenario: "creates a window",
			function: testNewWindow,
		},
		{
			scenario: "returns the menu bar",
			function: testMenuBar,
		},
		{
			scenario: "returns the dock tile",
			function: testDock,
		},
		{
			scenario: "shares",
			function: testShare,
		},
		{
			scenario: "creates a file panel",
			function: testNewFilePanel,
		},
		{
			scenario: "creates a popup notification",
			function: testNewPopupNotification,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, test.function)
	}
}

func testImport(t *testing.T) {
	app.Import(&Component{})
}

func testImportInvalidComponent(t *testing.T) {
	defer func() { recover() }()

	app.Import(InvalidComponent{})
	t.Error("no panic")
}

func testRunDriverError(t *testing.T) {
	defer func() { recover() }()

	app.Run(&test.Driver{
		RunSouldErr: true,
	})
	t.Error("no panic")
}

func testRunningDriverPanic(t *testing.T) {
	defer func() { recover() }()

	app.RunningDriver()
	t.Error("no panic")
}

func testRun(t *testing.T, driver *test.Driver) {
	app.Run(driver)
}

func testRunMultiple(t *testing.T) {
	defer func() { recover() }()

	app.Run(&test.Driver{})
	t.Error("no panic")
}

func testImportWhenDriverRuns(t *testing.T) {
	defer func() { recover() }()

	app.Import(&Component{})
	t.Error("no panic")
}

func testRunningDriver(t *testing.T, driver *test.Driver) {
	if app.RunningDriver() != driver {
		t.Error("driver is not the running driver")
	}
}

func testRender(t *testing.T, driver *test.Driver) {
	var compo app.Component
	driver.OnWindowLoad = func(w app.Window, c app.Component) {
		compo = c
	}
	defer func() {
		driver.OnWindowLoad = nil
	}()

	window := driver.NewWindow(app.WindowConfig{
		DefaultURL: "app.component",
	})
	defer window.Close()

	app.Render(compo)
}

func testContext(t *testing.T, driver *test.Driver) {
	var compo app.Component
	driver.OnWindowLoad = func(w app.Window, c app.Component) {
		compo = c
	}
	defer func() {
		driver.OnWindowLoad = nil
	}()

	window := driver.NewWindow(app.WindowConfig{
		DefaultURL: "app_test.component",
	})
	defer window.Close()

	ctx, err := app.Context(compo)
	if err != nil {
		t.Fatal(err)
	}
	if ctx != window {
		t.Fatal("returned context is not the window")
	}
}

func testContextError(t *testing.T) {
	_, err := app.Context(&Component{})
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testNewContextMenu(t *testing.T) {
	if menu := app.NewContextMenu(app.MenuConfig{}); menu == nil {
		t.Fatal("menu is nil")
	}
}

func testResources(t *testing.T) {
	resources := app.Resources()
	if len(resources) == 0 {
		t.Fatal("resources is empty")
	}
	t.Log(resources)
}

func testCallOnUIGoroutine(t *testing.T, d *test.Driver) {
	done := make(chan struct{})

	go func() {
		f := <-d.UIchan
		f()
	}()

	app.CallOnUIGoroutine(func() {
		done <- struct{}{}
	})
	<-done
}

func testStorage(t *testing.T) {
	if !app.SupportsStorage() {
		t.Fatal("storage is not supported")
	}

	storage := app.Storage()
	if len(storage) == 0 {
		t.Fatal("storage is empty")
	}
	t.Log(storage)
}

func testNewWindow(t *testing.T) {
	if !app.SupportsWindows() {
		t.Fatal("windows are no supported")
	}

	if window := app.NewWindow(app.WindowConfig{}); window == nil {
		t.Fatal("window is nil")
	}
}

func testMenuBar(t *testing.T) {
	if !app.SupportsMenuBar() {
		t.Fatal("menu bar is not supported")
	}

	if menubar := app.MenuBar(); menubar == nil {
		t.Fatal("menu bar is nil")
	}
}

func testDock(t *testing.T) {
	if !app.SupportsDock() {
		t.Fatal("dock is not supported")
	}

	if dock := app.Dock(); dock == nil {
		t.Fatal("dock is nil")
	}
}

func testShare(t *testing.T) {
	if !app.SupportsShare() {
		t.Fatal("share is not supported")
	}

	app.Share(42)
}

func testNewFilePanel(t *testing.T) {
	if !app.SupportsFilePanels() {
		t.Fatal("file panels are not supported")
	}

	if panel := app.NewFilePanel(app.FilePanelConfig{}); panel == nil {
		t.Fatal("pannel is nil")
	}
}

func testNewPopupNotification(t *testing.T) {
	if !app.SupportsPopupNotifications() {
		t.Fatal("popup notifications are not supported")
	}

	if popup := app.NewPopupNotification(app.PopupNotificationConfig{}); popup == nil {
		t.Fatal("popup is nil")
	}
}
