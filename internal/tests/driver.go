package tests

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

// DriverSetup is the definition of a function that creates a driver.
type DriverSetup func(onRun func()) app.Driver

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

		p := d.NewPage(app.PageConfig{})
		assertElem(t, p)
		t.Run("page", func(t *testing.T) { testPage(t, p) })

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

	d = setup(onRun)
	d.Run(f)

}

func testElemByCompo(t *testing.T, d app.Driver) {
	tests := []struct {
		scenario string
		elem     app.ElemWithCompo
	}{
		{
			scenario: "window",
			elem:     d.NewWindow(app.WindowConfig{URL: "tests.Hello"}),
		},
		{
			scenario: "page",
			elem:     d.NewPage(app.PageConfig{URL: "tests.Hello"}),
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

func testElemWithCompo(t *testing.T, e app.ElemWithCompo) {
	assert.NotEmpty(t, e.ID())

	e.Load("tests.Unknown")
	assert.Error(t, e.Err())

	e.Load("tests.Hello")
	assertElem(t, e)

	c := e.Compo()
	if e.Err() == app.ErrNotSupported {
		assert.Nil(t, c)
	} else {
		assertElem(t, e)
		assert.NotNil(t, c)
	}

	assert.True(t, e.Contains(c))
	assert.False(t, e.Contains(&Hello{}))

	e.Render(c)
	assertElem(t, e)

	e.Render(&Hello{})
	assert.Error(t, e.Err())
}

func testNavigator(t *testing.T, n app.Navigator, lazy bool) {
	n.Reload()
	assert.Error(t, n.Err())

	n.Load("tests.Hello")
	assertElem(t, n)

	if lazy {
		n.CanPrevious()
		assert.NoError(t, n.Err())

		n.Previous()
		assert.NoError(t, n.Err())

		n.CanNext()
		assert.NoError(t, n.Err())

		n.Next()
		assert.NoError(t, n.Err())
		return
	}

	assert.False(t, n.CanPrevious())
	assert.False(t, n.CanNext())

	n.Previous()
	assert.Error(t, n.Err())

	n.Next()
	assert.Error(t, n.Err())

	n.Load("tests.World")
	assert.True(t, n.CanPrevious())
	assert.False(t, n.CanNext())

	n.Previous()
	assertElem(t, n)
	assert.False(t, n.CanPrevious())
	assert.True(t, n.CanNext())

	n.Next()
	assertElem(t, n)
	assert.True(t, n.CanPrevious())
	assert.False(t, n.CanNext())
}

func assertElem(t *testing.T, e app.Elem) {
	if e.Err() == app.ErrNotSupported {
		return
	}
	assert.NoError(t, e.Err())
}
