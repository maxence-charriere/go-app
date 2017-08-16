package app

import "github.com/murlokswarm/app/markup"

// Driver is the interface that describes a backend for app rendering.
type Driver interface {
	// Run runs the application with the components resistered in the component
	// builder b.
	Run(b markup.CompoBuilder) error

	// Resources returns the location of the resources directory.
	Resources() string
}

// DriverWithStorage is the interface that describes a driver which supports
// storage.
// A storage is the local directory where an application can save files.
// This is the ideal location to save dynamic contents and local database files.
type DriverWithStorage interface {
	Driver

	// Storage returns the location of the storage directory.
	Storage() string
}

// DriverWithWindows is the interface that describes a driver able to create
// windows.
type DriverWithWindows interface {
	Driver

	// NewWindow creates and displays a window described by configuration c.
	NewWindow(c WindowConfig) Window
}

// DriverWithMenuBar is the interface that describes a driver with a menu bar.
type DriverWithMenuBar interface {
	Driver

	// MenuBar returns the menu bar.
	MenuBar() MenuBar
}

// DriverWithDock is the interface that describes a driver with a dock.
type DriverWithDock interface {
	Driver

	// Dock returns the dock.
	Dock() Dock
}
