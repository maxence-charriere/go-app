package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Menu is the interface that describes a menu.
type Menu ElementWithComponent

// MenuConfig is a struct that describes a menu.
type MenuConfig struct {
	DefaultURL string

	OnClose func()
}

// NewMenuWithLogs returns a decorated version of the given menu that logs
// all the operations.
// Uses the default logger.
func NewMenuWithLogs(m Menu, name string) Menu {
	return &menuWithLogs{
		name: name,
		base: m,
	}
}

type menuWithLogs struct {
	name string
	base Menu
}

func (m *menuWithLogs) ID() uuid.UUID {
	id := m.base.ID()
	Logf("%s id is %s", m.name, id)
	return id
}

func (m *menuWithLogs) Load(url string, v ...interface{}) error {
	fmtURL := fmt.Sprintf(url, v...)
	Logf("%s %s: loading %s", m.name, m.base.ID(), fmtURL)

	err := m.base.Load(url, v...)
	if err != nil {
		Errorf("%s %s: loading %s failed: %s", m.name, m.base.ID(), fmtURL, err)
	}
	return err
}

func (m *menuWithLogs) Component() Component {
	c := m.base.Component()
	Logf("%s %s: mounted component is %T", m.name, m.base.ID(), c)
	return c
}

func (m *menuWithLogs) Contains(c Component) bool {
	ok := m.base.Contains(c)
	Logf("%s %s: contains %T is %v", m.name, m.base.ID(), c, ok)
	return ok
}

func (m *menuWithLogs) Render(c Component) error {
	Logf("%s %s: rendering component %T", m.name, m.base.ID(), c)

	err := m.base.Render(c)
	if err != nil {
		Errorf("%s %s: rendering %T failed: %s", m.name, m.base.ID(), c, err)
	}
	return err
}

func (m *menuWithLogs) LastFocus() time.Time {
	focused := m.base.LastFocus()
	Logf("%s %s: last focus at %v", m.name, m.base.ID(), focused)
	return focused
}
