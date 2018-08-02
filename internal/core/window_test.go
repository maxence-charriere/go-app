package core

import (
	"testing"

	"github.com/murlokswarm/app"

	"github.com/stretchr/testify/assert"
)

func TestWindow(t *testing.T) {
	w := &Window{}

	whenWinCalled := false
	w.WhenWindow(func(w app.Window) {
		whenWinCalled = true
	})
	assert.True(t, whenWinCalled)

	whenNavCalled := false
	w.WhenNavigator(func(n app.Navigator) {
		whenNavCalled = true
	})
	assert.True(t, whenNavCalled)

	w.Load("")
	assert.Error(t, w.Err())

	assert.Nil(t, w.Compo())
	assert.False(t, w.Contains(nil))

	w.Render(nil)
	assert.Error(t, w.Err())

	w.Reload()
	assert.Error(t, w.Err())

	assert.False(t, w.CanPrevious())

	w.Previous()
	assert.Error(t, w.Err())

	assert.False(t, w.CanNext())

	w.Next()
	assert.Error(t, w.Err())

	x, y := w.Position()
	assert.Zero(t, x)
	assert.Zero(t, y)
	assert.Error(t, w.Err())

	w.Move(42, 42)
	assert.Error(t, w.Err())

	w.Center()
	assert.Error(t, w.Err())

	width, height := w.Size()
	assert.Zero(t, width)
	assert.Zero(t, height)
	assert.Error(t, w.Err())

	w.Resize(42, 42)
	assert.Error(t, w.Err())

	w.Focus()
	assert.Error(t, w.Err())

	w.FullScreen()
	assert.Error(t, w.Err())

	w.ExitFullScreen()
	assert.Error(t, w.Err())

	w.Minimize()
	assert.Error(t, w.Err())

	w.Deminimize()
	assert.Error(t, w.Err())

	w.Close()
	assert.Error(t, w.Err())
}
