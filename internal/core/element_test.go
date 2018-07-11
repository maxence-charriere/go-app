package core

import (
	"testing"

	"github.com/murlokswarm/app"
)

func TestElementWithComponentBase(t *testing.T) {
	e := &ElementWithComponentBase{}
	e.WhenWindow(func(w app.Window) {})
	e.WhenPage(func(p app.Page) {})
	e.WhenDockTile(func(d app.DockTile) {})
	e.WhenStatusMenu(func(s app.StatusMenu) {})
}
