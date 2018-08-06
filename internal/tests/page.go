package tests

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

func testPage(t *testing.T, p app.Page) {
	// app.Elem
	called := false
	p.WhenWindow(func(w app.Window) {
		called = true
	})
	assert.False(t, called)

	called = false
	p.WhenPage(func(p app.Page) {
		called = true
	})
	assert.True(t, called)

	called = false
	p.WhenNavigator(func(n app.Navigator) {
		called = true
	})
	assert.True(t, called)

	called = false
	p.WhenMenu(func(m app.Menu) {
		called = true
	})
	assert.False(t, called)

	called = false
	p.WhenDockTile(func(d app.DockTile) {
		called = true
	})
	assert.False(t, called)

	called = false
	p.WhenStatusMenu(func(s app.StatusMenu) {
		called = true
	})
	assert.False(t, called)

	p.WhenErr(func(err error) {
		t.Log(err)
	})

	t.Run("navigator", func(t *testing.T) {
		testNavigator(t, p, true)
	})

	t.Run("compo", func(t *testing.T) {
		testElemWithCompo(t, p)
	})

	p.URL()
	assertElem(t, p)

	p.Referer()
	assertElem(t, p)

	p.Close()
	assertElem(t, p)
}
