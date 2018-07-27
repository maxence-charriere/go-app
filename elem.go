package app

import (
	"time"
)

// Elem is the interface that describes an app element.
type Elem interface {
	// ID returns the element identifier.
	ID() string

	// WhenWindow calls the given func when the element is a window.
	WhenWindow(func(Window))

	// WhenPage calls the given func when the element is a page.
	WhenPage(func(Page))

	// WhenPage calls the given func when the element supports navigation.
	WhenNavigator(func(Navigator))

	// WhenStatusMenu calls the given func when the element is a menu.
	WhenMenu(func(Menu))

	// WhenDockTile calls the given func when the element is a dock tile.
	WhenDockTile(func(DockTile))

	// WhenStatusMenu calls the given func when the element is a status menu.
	WhenStatusMenu(func(StatusMenu))

	// WhenErr call the given func when the element is in an error state.
	WhenErr(func(err error))

	// Err returns the error that prevent the element to work.
	Err() error
}

// ElemWithCompo is the interface that describes an element that hosts
// components.
type ElemWithCompo interface {
	Elem

	// Load loads the page specified by the URL.
	// URL can be formated as fmt package functions.
	// Calls with an URL which contains a component name will load the named
	// component.
	// e.g. hello will load the component named hello.
	// It returns an error if the component is not imported.
	Load(url string, v ...interface{}) error

	// Compo returns the loaded component.
	Compo() Compo

	// Contains reports whether the component is mounted in the element.
	Contains(Compo) bool

	// Render renders the component.
	Render(Compo) error

	// LastFocus returns the last time when the element was focused.
	LastFocus() time.Time
}

// Navigator is the interface that describes an element that supports
// navigation.
type Navigator interface {
	ElemWithCompo

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
