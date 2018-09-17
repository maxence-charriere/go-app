// +build windows

// Package win is the driver to be used for applications that will run on
// Windows.
package win

import (
	"fmt"
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Driver is the app.Driver implementation for Windows.
type Driver struct {
	core.Driver

	// The URL of the component to load in the main window.
	URL string
}

// Run satisfies the app.Driver interface.
func (d *Driver) Run(f *app.Factory) error {
	fmt.Println("hello")
	time.Sleep(time.Second * 10)
	return nil
}
