package app

// DockTile is the interface that describes a dock tile.
type DockTile interface {
	Menu

	// SetIcon set the dock tile icon with the named file.
	// It returns an error if the file doesn't exist or if it is not a supported
	// image.
	SetIcon(name string) error

	// SetBadge set the dock tile badge with the string representation of the
	// value.
	SetBadge(v interface{}) error
}

// NewDockTileWithLogs returns a decorated version of the given dock tile that
// logs all the operations.
func NewDockTileWithLogs(d DockTile) DockTile {
	return &dockTileWithLogs{
		menuWithLogs: menuWithLogs{
			name: "dock",
			base: d,
		},
		base: d,
	}
}

type dockTileWithLogs struct {
	menuWithLogs
	base DockTile
}

func (d *dockTileWithLogs) SetIcon(name string) error {
	Logf("%s %s: set icon with %s", d.name, d.base.ID(), name)

	err := d.base.SetIcon(name)
	if err != nil {
		Errorf("%s %s: set icon failed: %s", d.name, d.base.ID(), err)
	}
	return err
}

func (d *dockTileWithLogs) SetBadge(v interface{}) error {
	Logf("%s %s: set badge with %+v", d.name, d.base.ID(), v)

	err := d.base.SetBadge(v)
	if err != nil {
		Errorf("%s %s: set badge failed: %s", d.name, d.base.ID(), err)
	}
	return err
}

// NewConcurrentDockTile returns a decorated version of the given dock tile that
//  is safe for concurrent operations.
func NewConcurrentDockTile(d DockTile) DockTile {
	return &concurrentDockTile{
		concurrentMenu: concurrentMenu{
			base: d,
		},
		base: d,
	}
}

type concurrentDockTile struct {
	concurrentMenu
	base DockTile
}

func (d *concurrentDockTile) SetIcon(name string) error {
	d.mutex.Lock()
	err := d.base.SetIcon(name)
	d.mutex.Unlock()
	return err
}

func (d *concurrentDockTile) SetBadge(v interface{}) error {
	d.mutex.Lock()
	err := d.base.SetBadge(v)
	d.mutex.Unlock()
	return err
}
