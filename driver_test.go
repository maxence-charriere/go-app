package app

import (
	"fmt"

	"github.com/murlokswarm/log"
)

type AbstractDriver struct {
	dock    Docker
	appMenu Contexter
}

func (d *AbstractDriver) Run() {
	log.Info("Running app")
}

func (d *AbstractDriver) NewElement(elem interface{}) Elementer {
	switch elem.(type) {
	case Window:
		return newWindowContext()

	default:
		return newTestContext(fmt.Sprintf("%T", elem))
	}
}

func (d *AbstractDriver) MenuBar() Contexter {
	return d.appMenu
}

func (d *AbstractDriver) Dock() Docker {
	return d.dock
}

func (d *AbstractDriver) Storage() Storer {
	return StorageTest("")
}

func (d *AbstractDriver) JavascriptBridge() string {
	return "alert('bridge not implemented');"
}

func (d *AbstractDriver) Share() Sharer {
	return &ShareTest{}
}

func (d *AbstractDriver) OpenFileChooser(fc FileChooser) {
	return
}

func init() {
	RegisterDriver(&AbstractDriver{
		dock:    newDockContext(),
		appMenu: newTestContext("appMenu"),
	})
}
