package core

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	testMenu(t, &Menu{})
}

func TestDockTile(t *testing.T) {
	d := &DockTile{}
	testMenu(t, d)

	whenDockTileCalled := false
	d.WhenDockTile(func(m app.DockTile) {
		whenDockTileCalled = true
	})
	assert.True(t, whenDockTileCalled)

	d.SetIcon("")
	assert.Error(t, d.Err())

	d.SetBadge("")
	assert.Error(t, d.Err())
}

func TestStatusMenu(t *testing.T) {
	s := &StatusMenu{}
	testMenu(t, s)

	whenStatusMenuCalled := false
	s.WhenStatusMenu(func(m app.StatusMenu) {
		whenStatusMenuCalled = true
	})
	assert.True(t, whenStatusMenuCalled)

	s.SetIcon("")
	assert.Error(t, s.Err())

	s.SetText("")
	assert.Error(t, s.Err())

	s.Close()
	assert.Error(t, s.Err())
}

func testMenu(t *testing.T, m app.Menu) {
	whenMenuCalled := false
	m.WhenMenu(func(m app.Menu) {
		whenMenuCalled = true
	})
	assert.True(t, whenMenuCalled)

	m.Load("")
	assert.Error(t, m.Err())

	assert.Nil(t, m.Compo())
	assert.False(t, m.Contains(nil))

	m.Render(nil)
	assert.Error(t, m.Err())

	assert.NotEmpty(t, m.Type())
}
