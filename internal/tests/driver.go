package tests

import (
	"os"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/murlokswarm/app/internal/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// NewMinimalDriver creates a driver that have nearly nothing supported.
func NewMinimalDriver(c app.DriverConfig) *core.Driver {
	d := &core.Driver{
		Elems:    core.NewElemDB(),
		Events:   c.Events,
		Factory:  c.Factory,
		Platform: &core.Platform{},
		UIChan:   c.UI,
	}

	d.Platform.Handler = func(call string) error {
		returnID := gjson.Get(call, "ReturnID").Str
		d.Platform.Return(returnID, "", "not supported")
		return nil
	}

	return d
}

// NewDriver creates a driver that covers general test cases.
func NewDriver(c app.DriverConfig) *core.Driver {
	d := &core.Driver{
		Elems:   core.NewElemDB(),
		Events:  c.Events,
		Factory: c.Factory,
		OpenDefaultBrowserFunc: func(string) error {
			return nil
		},
		NewContextMenuFunc: func(d *core.Driver) *core.Menu {
			return &core.Menu{
				DOM:    dom.Engine{Resources: d.Resources},
				Driver: d,
			}
		},
		NewDockTileFunc: func(d *core.Driver) *core.DockTile {
			return &core.DockTile{
				Menu: core.Menu{
					DOM:       dom.Engine{Resources: d.Resources},
					Driver:    d,
					NoDestroy: true,
				},
			}
		},
		NewMenuBarFunc: func(d *core.Driver) *core.Menu {
			return &core.Menu{
				DOM:       dom.Engine{Resources: d.Resources},
				Driver:    d,
				NoDestroy: true,
			}
		},
		NewStatusMenuFunc: func(d *core.Driver) *core.StatusMenu {
			return &core.StatusMenu{
				Menu: core.Menu{
					DOM:    dom.Engine{Resources: d.Resources},
					Driver: d,
				},
			}
		},
		NewWindowFunc: func(d *core.Driver) *core.Window {
			w := &core.Window{
				DOM: dom.Engine{
					Resources: d.Resources,
					AttrTransforms: []dom.Transform{
						dom.JsToGoHandler,
						dom.HrefCompoFmt,
					},
				},
				Driver: d,
			}

			w.SetPosition(42, 42)
			w.SetSize(42, 42)
			w.SetIsFocus(false)
			w.SetIsFullScreen(false)
			w.SetIsMinimized(false)
			return w
		},
		Platform: &core.Platform{},
		UIChan:   c.UI,
	}

	d.Platform.Handler = func(call string) error {
		returnID := gjson.Get(call, "ReturnID").Str
		d.Platform.Return(returnID, "", "")
		return nil
	}

	d.NewDockTile(app.MenuConfig{})
	return d
}

// TestDriver is a test suite that test a app.Driver.
func TestDriver(t *testing.T, d app.Driver, c app.DriverConfig) {
	defer os.RemoveAll(d.Resources())
	err := os.MkdirAll(d.Resources(), 0755)
	require.NoError(t, err)

	err = file.Copy(d.Resources("logo.png"), file.RepoPath("logo.png"))
	require.NoError(t, err)

	c.Factory.RegisterCompo(&Foo{})
	c.Factory.RegisterCompo(&Bar{})
	c.Factory.RegisterCompo(&Menu{})

	sub := app.Subscriber{Events: c.Events}

	sub.Subscribe(app.Running, func() {
		t.Log("app name:", d.AppName())

		notSetElem := d.ElemByCompo(&Foo{})
		assert.Equal(t, app.ErrElemNotSet, notSetElem.Err())

		testDockTile(t, d, d.DockTile())
		testMenu(t, d.MenuBar())
		testMenu(t, d.NewContextMenu(app.MenuConfig{URL: "tests.menu"}))
		testController(t, d.NewController(app.ControllerConfig{}))
		testElem(t, d.NewFilePanel(app.FilePanelConfig{}))
		testElem(t, d.NewNotification(app.NotificationConfig{}))
		testElem(t, d.NewSaveFilePanel(app.SaveFilePanelConfig{}))
		testStatusMenu(t, d, d.NewStatusMenu(app.StatusMenuConfig{URL: "tests.menu"}))
		testElem(t, d.NewShare(nil))
		testWindow(t, d.NewWindow(app.WindowConfig{URL: "tests.foo"}))

		assertSupported(t, d.OpenDefaultBrowser("https://github.com"))
		d.Render(&Foo{})

		t.Log(d.Resources("test"))
		t.Log(d.Storage("test"))
		t.Log(d.Target())

		d.UI(d.Stop)
	})

	err = d.Run(c)
	t.Log(err)
}

func testElem(t *testing.T, e app.Elem) {
	if !assertSupported(t, e.Err()) {
		return
	}

	t.Log("elem id:", e.ID())
	assert.False(t, e.Contains(&Foo{}))
}

func testDockTile(t *testing.T, d app.Driver, dt app.DockTile) {
	if !assertSupported(t, dt.Err()) {
		return
	}

	testMenu(t, dt)

	isDockTile := false
	dt.WhenDockTile(func(app.DockTile) { isDockTile = true })
	assert.True(t, isDockTile)

	dt.SetIcon(d.Resources("logo.png"))
	assertSupported(t, dt.Err())

	dt.SetBadge("hello")
	assertSupported(t, dt.Err())
}

func testController(t *testing.T, c app.Controller) {
	if !assertSupported(t, c.Err()) {
		return
	}

	c.Close()
}

func testMenu(t *testing.T, m app.Menu) {
	if !assertSupported(t, m.Err()) {
		return
	}

	testElem(t, m)

	isMenu := false
	m.WhenMenu(func(app.Menu) { isMenu = true })
	assert.True(t, isMenu)

	isView := false
	m.WhenView(func(app.View) { isView = true })
	assert.True(t, isView)

	t.Log(m.Kind())

	assert.False(t, m.CanPrevious())
	assert.False(t, m.CanNext())

	m.Previous()
	assert.Error(t, m.Err())

	m.Next()
	assert.Error(t, m.Err())

	m.Load("tests.menu")
	assert.NoError(t, m.Err())

	m.Reload()
	assert.NoError(t, m.Err())

	m.Load("tests.menu?idx=1")
	assert.NoError(t, m.Err())

	m.Previous()
	assert.NoError(t, m.Err())

	m.Next()
	assert.NoError(t, m.Err())

	compo := m.Compo()
	assert.NotNil(t, compo)
	assert.True(t, m.Contains(compo))

	m.Render(compo)
	assert.NoError(t, m.Err())

	m.Load("tests.foo")
	assert.Error(t, m.Err())
}

func testStatusMenu(t *testing.T, d app.Driver, s app.StatusMenu) {
	if !assertSupported(t, s.Err()) {
		return
	}

	testMenu(t, s)

	isStatusMenu := false
	s.WhenStatusMenu(func(app.StatusMenu) { isStatusMenu = true })
	assert.True(t, isStatusMenu)

	s.SetIcon(d.Resources("logo.png"))
	assertSupported(t, s.Err())

	s.SetText("hello")
	assertSupported(t, s.Err())

	s.Close()
}

func testWindow(t *testing.T, w app.Window) {
	if !assertSupported(t, w.Err()) {
		return
	}

	testElem(t, w)

	isView := false
	w.WhenView(func(app.View) { isView = true })
	assert.True(t, isView)

	isWindow := false
	w.WhenWindow(func(app.Window) { isWindow = true })
	assert.True(t, isWindow)

	assert.False(t, w.CanPrevious())
	assert.False(t, w.CanNext())

	w.Previous()
	assert.Error(t, w.Err())

	w.Next()
	assert.Error(t, w.Err())

	w.Load("tests.foo")
	assert.NoError(t, w.Err())

	w.Reload()
	assert.NoError(t, w.Err())

	w.Load("tests.foo?idx=1")
	assert.NoError(t, w.Err())

	w.Previous()
	assert.NoError(t, w.Err())

	w.Next()
	assert.NoError(t, w.Err())

	compo := w.Compo()
	assert.NotNil(t, compo)
	assert.True(t, w.Contains(compo))

	w.Render(compo)
	assert.NoError(t, w.Err())

	w.Move(42, 42)
	assertSupported(t, w.Err())

	w.Center()
	assertSupported(t, w.Err())

	w.Resize(42, 42)
	assertSupported(t, w.Err())

	w.Focus()
	assertSupported(t, w.Err())

	w.FullScreen()
	assertSupported(t, w.Err())

	w.ExitFullScreen()
	assertSupported(t, w.Err())

	w.Minimize()
	assertSupported(t, w.Err())

	w.Deminimize()
	assertSupported(t, w.Err())

	x, y := w.Position()
	t.Logf("positon: %v, %v", x, y)

	width, height := w.Size()
	t.Logf("size: %vx%v", width, height)

	t.Log("is focus:", w.IsFocus())
	t.Log("is full screen:", w.IsFullScreen())
	t.Log("is minimized:", w.IsMinimized())

	w.Close()
}

func assertSupported(t *testing.T, err error) bool {
	if err == app.ErrNotSupported {
		return false
	}

	return assert.NoError(t, err)
}
