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

	m.Load("tests.Unknown")
	assert.Error(t, m.Err())

	m.Load("tests.Menu")
	assert.Error(t, m.Err())

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
}

func testStatusMenu(t *testing.T, m app.StatusMenu) {
	t.Run("menu", func(t *testing.T) { testMenu(t, m) })

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
