package app

import "github.com/google/uuid"

// Element is the interface that describes an app element.
type Element interface {
	// ID returns the element identifier.
	ID() uuid.UUID
}

// Navigator is the interface that describes an element that supports
// navigation.
type Navigator interface {
	Element

	// Navigate navigates to the specified URL.
	// Calls with an URL which contains a component name will load the named
	// component.
	// e.g. /hello will load the imported component named hello.
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

// Window is the interface that describes a window.
type Window interface {
	Navigator

	// Position returns the window position.
	Position() (x, y float64)

	// Move moves the window to the position (x, y).
	Move(x, y float64)

	// Size returns the window size.
	Size() (width, height float64)

	// Resize resizes the window to width x height.
	Resize(width, height float64)

	// Focus gives the focus to the window.
	// The window will be put in front, above the other elements.
	Focus()

	// Close closes the window.
	Close()
}

// WindowConfig is a struct that describes a window.
type WindowConfig struct {
	Title           string
	X               float64
	Y               float64
	Width           float64
	MinWidth        float64
	MaxWidth        float64
	Height          float64
	MinHeight       float64
	MaxHeight       float64
	BackgroundColor string
	Borderless      bool
	DisableResize   bool
	Mac             MacWindowConfig

	OnMinimize       func()
	OnDeminimize     func()
	OnFullScreen     func()
	OnExitFullScreen func()
	OnMove           func(x, y float64)
	OnResize         func(width float64, height float64)
	OnFocus          func()
	OnBlur           func()
	OnClose          func() bool
}

// MacWindowConfig is a struct that describes window fields specific to MacOS.
type MacWindowConfig struct {
	BackgroundVibrancy Vibrancy
	HideCloseButton    bool
	HideMinimizeButton bool
	HideTitleBar       bool
}

// Vibrancy represents a constant that define Apple's frost glass effects.
type Vibrancy uint8

// Constants to specify vibrancy effects to use in Apple application elements.
const (
	VibeNone Vibrancy = iota
	VibeLight
	VibeDark
	VibeTitlebar
	VibeSelection
	VibeMenu
	VibePopover
	VibeSidebar
	VibeMediumLight
	VibeUltraDark
)

// Menu is the interface that describes a menu.
type Menu Navigator

// DockTile is the interface that describes a dock tile.
type DockTile interface {
	Navigator

	// SetIcon set the dock tile icon with the named file.
	// It returns an error if the file doesn't exist or if it is not a supported
	// image.
	SetIcon(name string) error

	// SetBadge set the dock tile badge with the string representation of v.
	SetBadge(v interface{})
}

// FilePanelConfig is a struct that describes a file panel.
type FilePanelConfig struct {
	MultipleSelection bool
	IgnoreDirectories bool
	IgnoreFiles       bool
	OnSelect          func(filenames []string)
}

// PopupNotificationConfig is a struct that describes a popup notification.
type PopupNotificationConfig struct {
	Message      string
	ComponentURL string
}
