package tests

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

func testMenu(t *testing.T, m app.Menu) {
	assert.NotEmpty(t, m.ID())

	called := false
	m.WhenWindow(func(w app.Window) {
		called = true
	})
	assert.False(t, called)

	called = false
	m.WhenPage(func(p app.Page) {
		called = true
	})
	assert.False(t, called)

	called = false
	m.WhenNavigator(func(n app.Navigator) {
		called = true
	})
	assert.False(t, called)

	called = false
	m.WhenMenu(func(m app.Menu) {
		called = true
	})
	assert.True(t, called)

	m.WhenErr(func(err error) {
		t.Log(err)
	})

	m.Load("tests.Unknown")
	assert.Error(t, m.Err())

	m.Load("tests.Menu")
	assertElem(t, m)

	c := m.Compo()
	if m.Err() == app.ErrNotSupported {
		assert.Nil(t, c)
	} else {
		assertElem(t, m)
		assert.NotNil(t, c)
	}

	assert.True(t, m.Contains(c))
	assert.False(t, m.Contains(&Menu{}))

	m.Render(c)
	assertElem(t, m)

	m.Render(&Menu{})
	assert.Error(t, m.Err())

	assert.NotEmpty(t, m.Type())
}

func testStatusMenu(t *testing.T, m app.StatusMenu) {
	t.Run("menu", func(t *testing.T) { testMenu(t, m) })

	called := false
	m.WhenDockTile(func(d app.DockTile) {
		called = true
	})
	assert.False(t, called)

	called = false
	m.WhenStatusMenu(func(s app.StatusMenu) {
		called = true
	})
	assert.True(t, called)

	m.SetIcon(filepath.Join(resourcesDir(), "resources", "unknown"))
	assert.Error(t, m.Err())

	m.SetIcon(filepath.Join(resourcesDir(), "resources", "logo.png"))
	assertElem(t, m)

	m.SetText("test")
	assertElem(t, m)

	m.Close()
	assertElem(t, m)

}

func testDock(t *testing.T, d app.DockTile) {
	t.Run("menu", func(t *testing.T) { testMenu(t, d) })

	called := false
	d.WhenDockTile(func(d app.DockTile) {
		called = true
	})
	assert.True(t, called)

	called = false
	d.WhenStatusMenu(func(s app.StatusMenu) {
		called = true
	})
	assert.False(t, called)

	d.SetIcon(filepath.Join(resourcesDir(), "resources", "unknown"))
	assert.Error(t, d.Err())

	d.SetIcon(filepath.Join(resourcesDir(), "resources", "logo.png"))
	assertElem(t, d)

	d.SetBadge("test")
	assertElem(t, d)

}

func resourcesDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}
