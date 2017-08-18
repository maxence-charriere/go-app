package app

import "github.com/google/uuid"

// Element is the interface that describes an app element.
type Element interface {
	// ID returns the element identifier.
	ID() uuid.UUID
}

type Navigator interface {
	Element

	// Navigate navigates to the specified URL.
	// Calls with an URL which contains a component name will load the named
	// component.
	// e.g. /hello
	Navigate(url string) error

	// CanPrevious indicates if navigation to previous page is possible.
	CanPrevious() bool

	// Previous navigates to the previous page.
	// It returns an error if there is no previous page to navigate.
	Previous() error

	// CanNext indicates if navigation to next page is possible.
	CanNext() bool

	// Next navigates to the next page.
	// It returns an error if there is no next page to navigate.
	Next() error
}

type Window interface {
	Navigator

	Position() (x, y float64)

	Move(x, y float64)

	Size() (width, height float64)

	Resize(width, height float64)

	Focus()

	Close()
}

type WindowConfig struct{}

type MenuBar interface {
	Navigator
}

type Dock interface {
	Navigator

	SetIcon(name string)

	SetBadge(v interface{})
}
