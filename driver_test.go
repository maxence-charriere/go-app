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

func (d *AbstractDriver) NewContext(ctx interface{}) Contexter {
	switch ctx.(type) {
	case Window:
		return newWindowCtx()

	default:
		return NewZeroContext(fmt.Sprintf("%T", ctx))
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

func init() {
	RegisterDriver(&AbstractDriver{
		dock:    newDockCtx(),
		appMenu: NewZeroContext("appMenu"),
	})
}
