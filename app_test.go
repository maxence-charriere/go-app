package app

import (
	"testing"

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
	d := &testDriver{
		test: t,
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
			name: "import component when driver is running should fail",
			test: testImportWhenDriverRuns,
		},
		{
			name: "import invalid component should fail",
			test: testImportInvalidComponent,
		},
		{
			name: "should run",
			test: func(t *testing.T) { testRun(t, d) },
		},
		{
			name: "second run should panic",
			test: testRunPanic,
		},
		{
			name: "should return the running driver",
			test: func(t *testing.T) { testRunningDriver(t, d) },
		},
		{
			name: "running driver when app is not running should panic",
			test: testRunningDriverPanic,
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
	Import(&Component{})
}

func testImportWhenDriverRuns(t *testing.T) {
	defer func() { recover() }()

	driver = &testDriver{}
	defer func() { driver = nil }()

	Import(&Component{})
	t.Error("should panic")
}

func testImportInvalidComponent(t *testing.T) {
	defer func() { recover() }()

	Import(InvalidComponent{})
	t.Error("should panic")
}

func testRun(t *testing.T, d *testDriver) {
	if err := Run(d); err != nil {
		t.Fatal(err)
	}
}

func testRunPanic(t *testing.T) {
	defer func() { recover() }()

	Run(&testDriver{
		test: t,
	})
	t.Error("should panic")
}

func testRunningDriver(t *testing.T, d *testDriver) {
	if RunningDriver() != d {
		t.Fatal("running driver should be d")
	}
}

func testRunningDriverPanic(t *testing.T) {
	d := driver
	driver = nil
	defer func() { driver = d }()
	defer func() { recover() }()

	RunningDriver()
	t.Error("should panic")
}

func testRender(t *testing.T, d *testDriver) {
	var compo markup.Component
	d.onWindowLoad = func(w Window, c markup.Component) {
		compo = c
	}
	defer func() {
		d.onWindowLoad = nil
	}()

	window := d.NewWindow(WindowConfig{
		DefaultURL: "app.component",
	})
	defer window.Close()

	Render(compo)
}

func testRenderLogError(t *testing.T) {
	Render(&Component{})
}

func testContext(t *testing.T, d *testDriver) {
	var compo markup.Component
	d.onWindowLoad = func(w Window, c markup.Component) {
		compo = c
	}
	defer func() {
		d.onWindowLoad = nil
	}()

	window := d.NewWindow(WindowConfig{
		DefaultURL: "app.component",
	})
	defer window.Close()

	ctx, err := Context(compo)
	if err != nil {
		t.Fatal(err)
	}
	if ctx != window {
		t.Fatal("returned context should be the window")
	}
}

func testContextError(t *testing.T) {
	_, err := Context(&Component{})
	if err == nil {
		t.Fatal("context should return an error")
	}
	t.Log(err)
}

func testNewContextMenu(t *testing.T) {
	if menu := NewContextMenu(MenuConfig{}); menu == nil {
		t.Fatal("menu should not be nil")
	}
}

func testResources(t *testing.T) {
	resources := Resources()
	if len(resources) == 0 {
		t.Fatal("resources should return a filepath")
	}
	t.Log(resources)
}

func testCallOnUIGoroutine(t *testing.T, d *testDriver) {
	done := make(chan struct{})

	go func() {
		f := <-d.uichan
		f()
	}()

	CallOnUIGoroutine(func() {
		done <- struct{}{}
	})
	<-done
}

func testStorage(t *testing.T) {
	if !SupportsStorage() {
		t.Fatal("storage should be supported")
	}

	storage := Storage()
	if len(storage) == 0 {
		t.Fatal("storage should return a filepath")
	}
	t.Log(storage)
}

func testNewWindow(t *testing.T) {
	if !SupportsWindows() {
		t.Fatal("windows should be supported")
	}

	if window := NewWindow(WindowConfig{}); window == nil {
		t.Fatal("window should not be nil")
	}
}

func testMenuBar(t *testing.T) {
	if !SupportsMenuBar() {
		t.Fatal("menu bar should be supported")
	}

	if menubar := MenuBar(); menubar == nil {
		t.Fatal("menu bar should not be nil")
	}
}

func testDock(t *testing.T) {
	if !SupportsDock() {
		t.Fatal("dock should be supported")
	}

	if dock := Dock(); dock == nil {
		t.Fatal("dock should not be nil")
	}
}

func testShare(t *testing.T) {
	if !SupportsShare() {
		t.Fatal("share should be supported")
	}

	Share(42)
}

func testNewFilePanel(t *testing.T) {
	if !SupportsFilePanels() {
		t.Fatal("file panels should be supported")
	}

	if panel := NewFilePanel(FilePanelConfig{}); panel == nil {
		t.Fatal("pannel should not be nil")
	}
}

func testNewPopupNotification(t *testing.T) {
	if !SupportsPopupNotifications() {
		t.Fatal("popup notifications should be supported")
	}

	if popup := NewPopupNotification(PopupNotificationConfig{}); popup == nil {
		t.Fatal("popup should not be nil")
	}
}
