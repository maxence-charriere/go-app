package app

import (
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
)

type AbstractDriver struct {
}

func (d *AbstractDriver) Run() {
	log.Info("Running app")
}

func (d *AbstractDriver) Render(target uid.ID, HTML string) (err error) {
	log.Infof("rendering %v:\n\033[32m%v\033[00m", target, HTML)
	return
}

func init() {
	RegisterDriver(&AbstractDriver{})
}
