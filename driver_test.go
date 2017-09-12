package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

type testDriver struct {
	test         *testing.T
	compoBuilder markup.CompoBuilder
	elements     ElementStore
	menubar      Menu
	dock         DockTile
	uichan       chan func()

	onWindowLoad func(w Window, c markup.Component)
}

func (d *testDriver) Run(b markup.CompoBuilder) error {
	d.compoBuilder = b
	d.elements = NewElementStore()

	d.menubar = newTestMenu(d, MenuConfig{})
	d.dock = newDockTile(d)

	d.uichan = make(chan func(), 256)
	return nil
}

func (d *testDriver) Render(c markup.Component) error {
	elem, err := d.elements.ElementByComponent(c)
	if err != nil {
		return errors.Wrap(err, "rendering component")
	}
	return elem.Render(c)
}

func (d *testDriver) Context(c markup.Component) (e ElementWithComponent, err error) {
	if e, err = d.elements.ElementByComponent(c); err != nil {
		err = errors.Wrap(err, "can't get context")
	}
	return
}

func (d *testDriver) NewContextMenu(c MenuConfig) Menu {
	return newTestMenu(d, c)
}

func (d *testDriver) Resources() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(wd, "resources")
}

func (d *testDriver) CallOnUIGoroutine(f func()) {
	d.uichan <- f
}

func (d *testDriver) Storage() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(wd, "storage")
}

func (d *testDriver) NewWindow(c WindowConfig) Window {
	return newTestWindow(d, c)
}

func (d *testDriver) MenuBar() Menu {
	return d.menubar
}

func (d *testDriver) Dock() DockTile {
	return d.dock
}

func (d *testDriver) Share(v interface{}) {
}

func (d *testDriver) NewFilePanel(c FilePanelConfig) Element {
	return newTestElement(d)
}

func (d *testDriver) NewPopupNotification(c PopupNotificationConfig) Element {
	return newTestElement(d)
}
