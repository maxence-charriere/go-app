package core

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	p := &Page{}

	whenPageCalled := false
	p.WhenPage(func(p app.Page) {
		whenPageCalled = true
	})
	assert.True(t, whenPageCalled)

	whenNavCalled := false
	p.WhenNavigator(func(n app.Navigator) {
		whenNavCalled = true
	})
	assert.True(t, whenNavCalled)

	p.Load("")
	assert.Error(t, p.Err())

	assert.Nil(t, p.Compo())
	assert.False(t, p.Contains(nil))

	p.Render(nil)
	assert.Error(t, p.Err())

	p.Reload()
	assert.Error(t, p.Err())

	assert.False(t, p.CanPrevious())

	p.Previous()
	assert.Error(t, p.Err())

	assert.False(t, p.CanNext())

	p.Next()
	assert.Error(t, p.Err())

	assert.Zero(t, p.URL())
	assert.Zero(t, p.Referer())

	p.Close()
	assert.Error(t, p.Err())
}
