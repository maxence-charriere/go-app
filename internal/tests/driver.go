package tests

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

// DriverSetup is the definition of a function that creates a driver.
type DriverSetup func() app.Driver

// TestDriver is a test suite that ensure that all driver implementations behave
// the same.
func TestDriver(t *testing.T, setup DriverSetup) {
	var d app.Driver

	onRun := func() {
		assert.NotEmpty(t, d.AppName())
		assert.NotEmpty(t, d.Resources())
		assert.NotEmpty(t, d.Storage())

		tmp := d.NewWindow(app.WindowConfig{URL: "tests.Hello"})
		if tmp.Err() == nil {
			c := tmp.Compo()
			assert.NotNil(t, c)
			d.Render(c)
			assertElem(t, tmp)
		}

		d.Render(&Hello{})

		t.Run("elem by compo", func(t *testing.T) { testElemByCompo(t, d) })

		w := d.NewWindow(app.WindowConfig{})
		assertElem(t, w)
		t.Run("window", func(t *testing.T) { testWindow(t, w) })

		c := d.NewController(app.ControllerConfig{})
		assertElem(t, c)
		t.Run("controller", func(t *testing.T) { testController(t, c) })

		cm := d.NewContextMenu(app.MenuConfig{})
		assertElem(t, cm)
		t.Run("context menu", func(t *testing.T) { testMenu(t, cm) })

		fp := d.NewFilePanel(app.FilePanelConfig{})
		assertElem(t, fp)

		sfp := d.NewSaveFilePanel(app.SaveFilePanelConfig{})
		assertElem(t, sfp)

		s := d.NewShare("")
		assertElem(t, s)

		n := d.NewNotification(app.NotificationConfig{})
		assertElem(t, n)

		mb := d.MenuBar()
		assertElem(t, mb)
		t.Run("menu bar", func(t *testing.T) { testMenu(t, mb) })

		sm := d.NewStatusMenu(app.StatusMenuConfig{})
		assertElem(t, sm)
		t.Run("status menu", func(t *testing.T) { testStatusMenu(t, sm) })

		dt := d.DockTile()
		assertElem(t, dt)
		t.Run("dock", func(t *testing.T) { testDock(t, dt) })

		d.Stop()
	}

	f := app.NewFactory()
	f.RegisterCompo(&Hello{})
	f.RegisterCompo(&World{})
	f.RegisterCompo(&Menu{})

	ui := make(chan func(), 32)
	defer close(ui)

	e := app.NewEventRegistry(ui)

	s := &app.Subscriber{
		Events: e,
	}
	defer s.Subscribe(app.Running, onRun).Close()

	d = setup()
	d.Run(app.DriverConfig{
		UI:      ui,
		Factory: f,
		Events:  e,
	})

}

func testElemByCompo(t *testing.T, d app.Driver) {
	tests := []struct {
		scenario string
		elem     app.View
	}{
		{
			scenario: "window",
			elem:     d.NewWindow(app.WindowConfig{URL: "tests.Hello"}),
		},
		{
			scenario: "context menu",
			elem:     d.NewContextMenu(app.MenuConfig{URL: "tests.Menu"}),
		},
		{
			scenario: "status menu",
			elem:     d.NewStatusMenu(app.StatusMenuConfig{URL: "tests.Menu"}),
		},
		{
			scenario: "dock",
			elem: func() app.DockTile {
				dt := d.DockTile()
				dt.Load("tests.Menu")
				return dt
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			e := test.elem
			assertElem(t, e)

			if e.Err() == app.ErrNotSupported {
				return
			}

			c := e.Compo()
			assert.NotNil(t, c)

			ebc := d.ElemByCompo(c)
			if e.Err() == nil {
				assert.Equal(t, e.ID(), ebc.ID())
			}
		})
	}
}

func testView(t *testing.T, v app.View) {
	assert.NotEmpty(t, v.ID())

	v.Load("tests.Unknown")
	assert.Error(t, v.Err())

	v.Load("tests.Hello")
	assertElem(t, v)

	c := v.Compo()
	if v.Err() == app.ErrNotSupported {
		assert.Nil(t, c)
	} else {
		assertElem(t, v)
		assert.NotNil(t, c)
	}

	assert.True(t, v.Contains(c))
	assert.False(t, v.Contains(&Hello{}))

	v.Render(c)
	assertElem(t, v)

	v.Render(&Hello{})
	assert.Error(t, v.Err())
}

func testViewNav(t *testing.T, v app.View, lazy bool) {
	v.Reload()
	assert.Error(t, v.Err())

	v.Load("tests.Hello")
	assertElem(t, v)

	if lazy {
		v.CanPrevious()
		assert.NoError(t, v.Err())

		v.Previous()
		assert.NoError(t, v.Err())

		v.CanNext()
		assert.NoError(t, v.Err())

		v.Next()
		assert.NoError(t, v.Err())
		return
	}

	assert.False(t, v.CanPrevious())
	assert.False(t, v.CanNext())

	v.Previous()
	assert.Error(t, v.Err())

	v.Next()
	assert.Error(t, v.Err())

	v.Load("tests.World")
	assert.True(t, v.CanPrevious())
	assert.False(t, v.CanNext())

	v.Previous()
	assertElem(t, v)
	assert.False(t, v.CanPrevious())
	assert.True(t, v.CanNext())

	v.Next()
	assertElem(t, v)
	assert.True(t, v.CanPrevious())
	assert.False(t, v.CanNext())
}

func assertElem(t *testing.T, e app.Elem) {
	if e.Err() == app.ErrNotSupported {
		return
	}
	assert.NoError(t, e.Err())
}
