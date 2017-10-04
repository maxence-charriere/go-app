package app

import "github.com/murlokswarm/app/markup"

// Driver is the interface that describes a backend for app rendering.
type Driver interface {
	// Run runs the application with the components resistered in the component
	// builder b.
	Run(b markup.CompoBuilder) error

	// Render renders component c.
	Render(c markup.Component) error

	// Context returns the element where component c is mounted.
	// It returns an error if c is not mounted.
	Context(c markup.Component) (ElementWithComponent, error)

	// NewContextMenu creates and displays the context menu described in
	// configuration c.
	NewContextMenu(c MenuConfig) Menu

	// Resources returns the location of the resources directory.
	Resources() string

	// Logs returns the application logger.
	Logs() Logger

	// CallOnUIGoroutine calls func f and ensure it's called from the UI
	// goroutine.
	CallOnUIGoroutine(f func())
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

	// NewWindow creates and displays the window described in configuration c.
	NewWindow(c WindowConfig) Window
}

// DriverWithMenuBar is the interface that describes a driver with a menu bar.
type DriverWithMenuBar interface {
	Driver

	// MenuBar returns the menu bar.
	MenuBar() Menu
}

// DriverWithDock is the interface that describes a driver with a dock.
type DriverWithDock interface {
	Driver

	// Dock returns the dock tile.
	Dock() DockTile
}

// DriverWithShare is the interface that describes a driver with sharing
// support.
type DriverWithShare interface {
	Driver

	// Share shares the value v.
	Share(v interface{})
}

// DriverWithFilePanels is the interface that describes a driver able to open a
// file panel.
type DriverWithFilePanels interface {
	Driver

	// NewFilePanel creates and displays the file panel described in
	// configuration c.
	NewFilePanel(c FilePanelConfig) Element
}

// DriverWithPopupNotifications is the interface that describes a driver able to
// display popup notifications.
type DriverWithPopupNotifications interface {
	// NewPopupNotification creates and displays the popup notification
	// described in configuration c.
	NewPopupNotification(c PopupNotificationConfig) Element
}
