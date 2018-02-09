package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Window is the interface that describes a window.
type Window interface {
	ElementWithNavigation

	// Position returns the window position.
	Position() (x, y float64)

	// Move moves the window to the position (x, y).
	Move(x, y float64)

	// Center moves the window to the center of the screen.
	Center()

	// Size returns the window size.
	Size() (width, height float64)

	// Resize resizes the window to width * height.
	Resize(width, height float64)

	// Focus gives the focus to the window.
	// The window will be put in front, above the other elements.
	Focus()

	// ToggleFullScreen takes the window into or out of fullscreen mode.
	ToggleFullScreen()

	// Minimize takes the window into or out of minimized mode
	ToggleMinimize()

	// Close closes the element.
	Close()
}

// WindowConfig is a struct that describes a window.
type WindowConfig struct {
	Title           string          `json:"title"`
	X               float64         `json:"x"`
	Y               float64         `json:"y"`
	Width           float64         `json:"width"`
	MinWidth        float64         `json:"min-width"`
	MaxWidth        float64         `json:"max-width"`
	Height          float64         `json:"height"`
	MinHeight       float64         `json:"min-height"`
	MaxHeight       float64         `json:"max-height"`
	BackgroundColor string          `json:"background-color"`
	NoResizable     bool            `json:"no-resizable"`
	NoClosable      bool            `json:"no-closable"`
	NoMinimizable   bool            `json:"no-minimizable"`
	TitlebarHidden  bool            `json:"titlebar-hidden"`
	DefaultURL      string          `json:"default-url"`
	Mac             MacWindowConfig `json:"mac"`

	OnMove           func(x, y float64)                  `json:"-"`
	OnResize         func(width float64, height float64) `json:"-"`
	OnFocus          func()                              `json:"-"`
	OnBlur           func()                              `json:"-"`
	OnFullScreen     func()                              `json:"-"`
	OnExitFullScreen func()                              `json:"-"`
	OnMinimize       func()                              `json:"-"`
	OnDeminimize     func()                              `json:"-"`
	OnClose          func() bool                         `json:"-"`
}

// MacWindowConfig is a struct that describes window fields specific to MacOS.
type MacWindowConfig struct {
	BackgroundVibrancy Vibrancy `json:"background-vibrancy"`
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

// NewWindowWithLogs returns a decorated version of the given window that logs
// all the operations.
// Use the default logger.
func NewWindowWithLogs(w Window) Window {
	return &windowWithLogs{
		base: w,
	}
}

type windowWithLogs struct {
	base Window
}

func (w *windowWithLogs) ID() uuid.UUID {
	id := w.base.ID()
}

func (w *windowWithLogs) Load(url string, v ...interface{}) error {
	fmtURL := fmt.Sprintf(url, v...)

	Logf("loading %s in window %s", fmtURL, w.base.ID())

	err := w.base.Load(url, v...)
	if err != nil {
		Errorf("loading %s in window %s failed: %s", fmtURL, w.base.ID(), err)
	}
	return err
}

func (w *windowWithLogs) Component() Component          {}
func (w *windowWithLogs) Contains(c Component) bool     {}
func (w *windowWithLogs) Render(c Component) error      {}
func (w *windowWithLogs) LastFocus() time.Time          {}
func (w *windowWithLogs) Reload() error                 {}
func (w *windowWithLogs) CanPrevious() bool             {}
func (w *windowWithLogs) Previous() error               {}
func (w *windowWithLogs) CanNext() bool                 {}
func (w *windowWithLogs) Next() error                   {}
func (w *windowWithLogs) Position() (x, y float64)      {}
func (w *windowWithLogs) Move(x, y float64)             {}
func (w *windowWithLogs) Center()                       {}
func (w *windowWithLogs) Size() (width, height float64) {}
func (w *windowWithLogs) Resize(width, height float64)  {}
func (w *windowWithLogs) Focus()                        {}
func (w *windowWithLogs) ToggleFullScreen()             {}
func (w *windowWithLogs) ToggleMinimize()               {}
func (w *windowWithLogs) Close()                        {}

func makeWindowConfigWithLogs(c WindowConfig) WindowConfig {
	return c
}

// NewConcurrentWindow returns a decorated version of the given window that is
// safe for concurrent operations.
func NewConcurrentWindow(w Window) Window {
	return &concurrentWindow{
		base: w,
	}
}

type concurrentWindow struct {
	mutex sync.Mutex
	base  Window
}

func (w *concurrentWindow) ID() uuid.UUID                           {}
func (w *concurrentWindow) Load(url string, v ...interface{}) error {}
func (w *concurrentWindow) Component() Component                    {}
func (w *concurrentWindow) Contains(c Component) bool               {}
func (w *concurrentWindow) Render(c Component) error                {}
func (w *concurrentWindow) LastFocus() time.Time                    {}
func (w *concurrentWindow) Reload() error                           {}
func (w *concurrentWindow) CanPrevious() bool                       {}
func (w *concurrentWindow) Previous() error                         {}
func (w *concurrentWindow) CanNext() bool                           {}
func (w *concurrentWindow) Next() error                             {}
func (w *concurrentWindow) Position() (x, y float64)                {}
func (w *concurrentWindow) Move(x, y float64)                       {}
func (w *concurrentWindow) Center()                                 {}
func (w *concurrentWindow) Size() (width, height float64)           {}
func (w *concurrentWindow) Resize(width, height float64)            {}
func (w *concurrentWindow) Focus()                                  {}
func (w *concurrentWindow) ToggleFullScreen()                       {}
func (w *concurrentWindow) ToggleMinimize()                         {}
func (w *concurrentWindow) Close()                                  {}
