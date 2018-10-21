package test

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom.v2"
)

// Menu is a test menu that implements the app.Menu interface.
type Menu struct {
	core.Menu

	driver *Driver
	id     string
	dom    dom.Engine
	compo  app.Compo
}

func newMenu(d *Driver, c app.MenuConfig) *Menu {
	m := &Menu{
		driver: d,
		id:     uuid.New().String(),
		dom:    dom.Engine{Factory: d.factory},
	}

	d.elems.Put(m)

	if len(c.URL) != 0 {
		m.Load(c.URL)
	}

	return m
}

// ID satisfies the app.Menu interface.
func (m *Menu) ID() string {
	return m.id
}

// Load satisfies the app.Menu interface.
func (m *Menu) Load(urlFmt string, v ...interface{}) {
	var err error
	defer func() {
		m.SetErr(err)
	}()

	u := fmt.Sprintf(urlFmt, v...)
	n := core.CompoNameFromURLString(u)

	var c app.Compo
	if c, err = m.driver.factory.NewCompo(n); err != nil {
		return
	}

	m.compo = c
	err = m.dom.New(c)
}

// Compo satisfies the app.Menu interface.
func (m *Menu) Compo() app.Compo {
	return m.compo
}

// Contains satisfies the app.Menu interface.
func (m *Menu) Contains(c app.Compo) bool {
	return m.dom.Contains(c)
}

// Render satisfies the app.Menu interface.
func (m *Menu) Render(c app.Compo) {
	m.SetErr(m.dom.Render(c))
}
