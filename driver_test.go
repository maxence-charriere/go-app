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

func (d *AbstractDriver) MenuBar() (menu Contexter, ok bool) {
	menu = d.appMenu
	ok = true
	return
}

func (d *AbstractDriver) Dock() (dock Docker, ok bool) {
	dock = d.dock
	ok = true
	return
}

func (d *AbstractDriver) Resources() string {
	return "resources"
}

func (d *AbstractDriver) Storage() string {
	return ""
}

func (d *AbstractDriver) JavascriptBridge() string {
	return "alert('bridge not implemented');"
}

func init() {
	RegisterDriver(&AbstractDriver{
		dock:    newDockContext(),
		appMenu: newTestContext("appMenu"),
	})
}
