package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/drivers/test"
	"github.com/murlokswarm/app/markup"
)

type Component markup.ZeroCompo

func (c *Component) Render() string {
	return `<div>Hello</div>`
}

type InvalidComponent markup.ZeroCompo

func (c InvalidComponent) Render() string {
	return ``
}

func TestApp(t *testing.T) {
	d := &test.Driver{
		Test: t,
	}

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should import component",
			test: testImport,
		},
		{
			name: "import invalid component should fail",
			test: testImportInvalidComponent,
		},
		{
			name: "run with driver error should panic",
			test: testRunDriverError,
		},
		{
			name: "get running driver when app is not running should panic",
			test: testRunningDriverPanic,
		},
		{
			name: "should run",
			test: func(t *testing.T) { testRun(t, d) },
		},
		{
			name: "second run should panic",
			test: testRunMultiple,
		},
		{
			name: "import component when driver is running should fail",
			test: testImportWhenDriverRuns,
		},
		{
			name: "should return the running driver",
			test: func(t *testing.T) { testRunningDriver(t, d) },
		},
		{
			name: "should render a component",
			test: func(t *testing.T) { testRender(t, d) },
		},
		{
			name: "render should log an error",
			test: testRenderLogError,
		},
		{
			name: "context should return an element",
			test: func(t *testing.T) { testContext(t, d) },
		},
		{
			name: "context should return an error",
			test: testContextError,
		},
		{
			name: "should create a context menu",
			test: testNewContextMenu,
		},
		{
			name: "resources should return a filepath",
			test: testResources,
		},
		{
			name: "logs should return the logger",
			test: testLogs,
		},
		{
			name: "should call on ui goroutine",
			test: func(t *testing.T) { testCallOnUIGoroutine(t, d) },
		},
		{
			name: "storage should return a filepath",
			test: testStorage,
		},
		{
			name: "should create a window",
			test: testNewWindow,
		},
		{
			name: "should return the menu bar",
			test: testMenuBar,
		},
		{
			name: "should return the dock tile",
			test: testDock,
		},
		{
			name: "should share",
			test: testShare,
		},
		{
			name: "should create a file panel",
			test: testNewFilePanel,
		},
		{
			name: "should create a popup notification",
			test: testNewPopupNotification,
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testImport(t *testing.T) {
	app.Import(&Component{})
}

func testImportWhenDriverRuns(t *testing.T) {
	defer func() { recover() }()

	app.Import(&Component{})
	t.Error("should panic")
}

func testImportInvalidComponent(t *testing.T) {
	defer func() { recover() }()

	app.Import(InvalidComponent{})
	t.Error("should panic")
}

func testRunDriverError(t *testing.T) {
	defer func() { recover() }()

	app.Run(&test.Driver{
		RunSouldErr: true,
	})
	t.Error("should panic")
}

func testRunningDriverPanic(t *testing.T) {
	defer func() { recover() }()

	app.RunningDriver()
	t.Error("should panic")
}

func testRun(t *testing.T, d *test.Driver) {
	app.Run(d)
}

func testRunMultiple(t *testing.T) {
	defer func() { recover() }()

	app.Run(&test.Driver{})
	t.Error("should panic")
}

func testRunningDriver(t *testing.T, d *test.Driver) {
	if app.RunningDriver() != d {
		t.Fatal("running driver should be d")
	}
}

func testRender(t *testing.T, d *test.Driver) {
	var compo markup.Component
	d.OnWindowLoad = func(w app.Window, c markup.Component) {
		compo = c
	}
	defer func() {
		d.OnWindowLoad = nil
	}()

	window := d.NewWindow(app.WindowConfig{
		DefaultURL: "app.component",
	})
	defer window.Close()

	app.Render(compo)
}

func testRenderLogError(t *testing.T) {
	app.Render(&Component{})
}

func testContext(t *testing.T, d *test.Driver) {
	var compo markup.Component
	d.OnWindowLoad = func(w app.Window, c markup.Component) {
		compo = c
	}
	defer func() {
		d.OnWindowLoad = nil
	}()

	window := d.NewWindow(app.WindowConfig{
		DefaultURL: "app_test.component",
	})
	defer window.Close()

	ctx, err := app.Context(compo)
	if err != nil {
		t.Fatal(err)
	}
	if ctx != window {
		t.Fatal("returned context should be the window")
	}
}

func testContextError(t *testing.T) {
	_, err := app.Context(&Component{})
	if err == nil {
		t.Fatal("context should return an error")
	}
	t.Log(err)
}

func testNewContextMenu(t *testing.T) {
	if menu := app.NewContextMenu(app.MenuConfig{}); menu == nil {
		t.Fatal("menu should not be nil")
	}
}

func testResources(t *testing.T) {
	resources := app.Resources()
	if len(resources) == 0 {
		t.Fatal("resources should return a filepath")
	}
	t.Log(resources)
}

func testLogs(t *testing.T) {
	app.Logs().Log("hello")
	app.Logs().Error("world")
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
		t.Fatal("storage should be supported")
	}

	storage := app.Storage()
	if len(storage) == 0 {
		t.Fatal("storage should return a filepath")
	}
	t.Log(storage)
}

func testNewWindow(t *testing.T) {
	if !app.SupportsWindows() {
		t.Fatal("windows should be supported")
	}

	if window := app.NewWindow(app.WindowConfig{}); window == nil {
		t.Fatal("window should not be nil")
	}
}

func testMenuBar(t *testing.T) {
	if !app.SupportsMenuBar() {
		t.Fatal("menu bar should be supported")
	}

	if menubar := app.MenuBar(); menubar == nil {
		t.Fatal("menu bar should not be nil")
	}
}

func testDock(t *testing.T) {
	if !app.SupportsDock() {
		t.Fatal("dock should be supported")
	}

	if dock := app.Dock(); dock == nil {
		t.Fatal("dock should not be nil")
	}
}

func testShare(t *testing.T) {
	if !app.SupportsShare() {
		t.Fatal("share should be supported")
	}

	app.Share(42)
}

func testNewFilePanel(t *testing.T) {
	if !app.SupportsFilePanels() {
		t.Fatal("file panels should be supported")
	}

	if panel := app.NewFilePanel(app.FilePanelConfig{}); panel == nil {
		t.Fatal("pannel should not be nil")
	}
}

func testNewPopupNotification(t *testing.T) {
	if !app.SupportsPopupNotifications() {
		t.Fatal("popup notifications should be supported")
	}

	if popup := app.NewPopupNotification(app.PopupNotificationConfig{}); popup == nil {
		t.Fatal("popup should not be nil")
	}
}
