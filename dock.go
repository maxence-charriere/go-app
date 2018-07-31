package app

import "fmt"

// DockTile is the interface that describes a dock tile.
// Accept only components that contain menu and menuitem tags.
type DockTile interface {
	Menu

	// SetIcon set the dock tile icon with the named file.
	// The icon should be a .png file.
	SetIcon(name string) error

	// SetBadge set the dock tile badge with the string representation of the
	// value.
	SetBadge(v interface{}) error
}

type dockWithLogs struct {
	DockTile
}

func (d *dockWithLogs) Load(url string, v ...interface{}) {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("dock tile is loading %s", parsedURL)
	})

	d.DockTile.Load(url, v...)
	if d.Err() != nil {
		Log("dock tile failed to load %s: %s",
			parsedURL,
			d.Err(),
		)
	}
}

func (d *dockWithLogs) Render(c Compo) {
	WhenDebug(func() {
		Debug("dock tile is rendering %T", c)
	})

	d.DockTile.Render(c)
	if d.Err() != nil {
		Log("dock tile failed to render %T: %s",
			c,
			d.Err(),
		)
	}
}

func (d *dockWithLogs) SetIcon(name string) error {
	WhenDebug(func() {
		Debug("dock tile is setting its icon to %s", name)
	})

	err := d.DockTile.SetIcon(name)
	if err != nil {
		Log("dock tile failed to set its icon: %s", err)
	}
	return err
}

func (d *dockWithLogs) SetBadge(v interface{}) error {
	WhenDebug(func() {
		Debug("dock tile is setting its badge to %d", v)
	})

	err := d.DockTile.SetBadge(v)
	if err != nil {
		Log("dock tile failed to set its badge: %s", err)
	}
	return err
}
