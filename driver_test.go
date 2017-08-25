package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/murlokswarm/app/markup"
	"github.com/pkg/errors"
)

type driverTest struct {
	Test     *testing.T
	elements ElementStore
}

func (d *driverTest) Run(b markup.CompoBuilder) error {
	d.Test.Logf("driver.Run")

	d.elements = NewElementStore()
	return nil
}

func (d *driverTest) Render(c markup.Component) error {
	d.Test.Logf("driver.Render: %T", c)

	elem, ok := d.elements.ElementByComponent(c)
	if !ok {
		panic(errors.Errorf("no element contain component %#v", c))
	}
	return elem.Render(c)
}

func (d *driverTest) Context(c markup.Component) (e ElementWithComponent, err error) {
	d.Test.Logf("driver.Context: %T", c)
	return
}

func (d *driverTest) Resources() string {
	d.Test.Logf("driver.Resources")

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(wd, "resources")
}
