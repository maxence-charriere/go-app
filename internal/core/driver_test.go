package core

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDriver(t *testing.T) {
	d := &Driver{}
	require.Implements(t, (*app.Driver)(nil), d)

	assert.Error(t, d.Run(nil))
	assert.Empty(t, d.AppName())
	assert.Equal(t, "resources", d.Resources())
	assert.Equal(t, "storage", d.Storage())
	assert.Error(t, d.Render(nil))
	assert.Error(t, d.ElemByCompo(nil).Err())

	w := d.NewWindow(app.WindowConfig{})
	assert.Error(t, w.Err())

	m := d.NewContextMenu(app.MenuConfig{})
	assert.Error(t, m.Err())

	assert.Error(t, d.NewPage(app.PageConfig{}))
	assert.Error(t, d.NewFilePanel(app.FilePanelConfig{}))
	assert.Error(t, d.NewSaveFilePanel(app.SaveFilePanelConfig{}))
	assert.Error(t, d.NewShare(nil))
	assert.Error(t, d.NewNotification(app.NotificationConfig{}))

	mb := d.MenuBar()
	assert.Error(t, mb.Err())

	_, err := d.NewStatusMenu(app.StatusMenuConfig{})
	assert.Error(t, err)

	dt := d.Dock()
	assert.Error(t, dt.Err())

	d.CallOnUIGoroutine(func() {
		t.Log("call from ui goroutine")
	})
}
