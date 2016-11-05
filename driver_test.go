package app

import (
	"fmt"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
)

type AbstractDriver struct {
	dock    Contexter
	appMenu Contexter
}

func (d *AbstractDriver) Run() {
	log.Info("Running app")
}

func (d *AbstractDriver) NewContext(ctx interface{}) Contexter {
	switch ctx.(type) {
	case Window:
		return NewZeroContext("window")

	default:
		return NewZeroContext(fmt.Sprintf("%T", ctx))
	}
}

func (d *AbstractDriver) Render(target uid.ID, HTML string) (err error) {
	log.Infof("rendering %v:\n\033[32m%v\033[00m", target, HTML)
	return
}

func (d *AbstractDriver) AppMenu() Contexter {
	return d.appMenu
}

func (d *AbstractDriver) Dock() Contexter {
	return d.dock
}

func (d *AbstractDriver) Resources() ResourceLocation {
	return "resources"
}

func (d *AbstractDriver) JavascriptBridge() string {
	return "alert('bridge not implemented');"
}

func init() {
	RegisterDriver(&AbstractDriver{
		dock:    NewZeroContext("dock"),
		appMenu: NewZeroContext("appMenu"),
	})
}
