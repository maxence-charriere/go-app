package app

import (
	"time"

	"github.com/google/uuid"
)

// Elem is the interface that describes an app element.
type Elem interface {
	// ID returns the element identifier.
	ID() uuid.UUID

	// WhenWindow calls the given func when the element is a window.
	WhenWindow(func(Window))

	// WhenPage calls the given func when the element is a page.
	WhenPage(func(Page))

	// WhenStatusMenu calls the given func when the element is a menu.
	WhenMenu(func(Menu))

	// WhenDockTile calls the given func when the element is a dock tile.
	WhenDockTile(func(DockTile))

	// WhenStatusMenu calls the given func when the element is a status menu.
	WhenStatusMenu(func(StatusMenu))

	// WhenNotSet call the given func when the element is not set.
	WhenNotSet(func())

	// IsNotSet reports whether the element is set.
	IsNotSet() bool
}

// ElementWithComponent is the interface that describes an element that hosts
// components.
type ElementWithComponent interface {
	Elem

	// Load loads the page specified by the URL.
	// URL can be formated as fmt package functions.
	// Calls with an URL which contains a component name will load the named
	// component.
	// e.g. hello will load the component named hello.
	// It returns an error if the component is not imported.
	Load(url string, v ...interface{}) error

	// Component returns the loaded component.
	Component() Component

	// Contains reports whether the component is mounted in the element.
	Contains(Component) bool

	// Render renders the component.
	Render(Component) error

	// LastFocus returns the last time when the element was focused.
	LastFocus() time.Time
}

// Navigator is the interface that describes an element that supports
// navigation.
type Navigator interface {
	ElementWithComponent

	// Reload reloads the current page.
	Reload() error

	// CanPrevious reports whether load the previous page is possible.
	CanPrevious() bool

	// Previous loads the previous page.
	// It returns an error if there is no previous page to load.
	Previous() error

	// CanNext indicates if loading next page is possible.
	CanNext() bool

	// Next loads the next page.
	// It returns an error if there is no next page to load.
	Next() error
}

// Closer is the interface that describes an element that can be closed.
type Closer interface {
	// Close closes the element and free its allocated resources.
	Close() error
}

// NotificationConfig is a struct that describes a notification.
type NotificationConfig struct {
	Title     string
	Subtitle  string
	Text      string
	ImageName string
	Sound     bool

	OnReply func(reply string) `json:"-"`
}
