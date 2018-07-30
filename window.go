package app

import (
	"fmt"
)

// Window is the interface that describes a window.
type Window interface {
	Navigator
	Closer

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

	// FullScreen takes the window into fullscreen mode.
	FullScreen()

	// ExitFullScreen takes the window out of fullscreen mode.
	ExitFullScreen()

	// Minimize takes the window into minimized mode
	Minimize()

	// Deminimize takes the window out of minimized mode
	Deminimize()
}

// WindowConfig is a struct that describes a window.
type WindowConfig struct {
	// The title.
	Title string

	// The default position on x axis.
	X float64

	// The default position on y axis.
	Y float64

	// The default width.
	Width float64

	// The minimum width.
	MinWidth float64

	// The maximum width.
	MaxWidth float64

	// The default height.
	Height float64

	// The minimum height.
	MinHeight float64

	// The maximum height.
	MaxHeight float64

	// The background color (#rrggbb).
	BackgroundColor string

	// Reports whether the window is resizable.
	FixedSize bool

	// Reports whether the close button is hidden.
	CloseHidden bool

	// Reports whether the minimize button is hidden.
	MinimizeHidden bool

	// Reports whether the title bar is hidden.
	TitlebarHidden bool

	// The URL of the component to load when the window is created.
	DefaultURL string

	// The MacOS window specific configuration.
	Mac MacWindowConfig

	// The function that is called when the window is moved.
	OnMove func(x, y float64) `json:"-"`

	// The function that is called when the window is resized.
	OnResize func(width float64, height float64) `json:"-"`

	// The function that is called when the window get focus.
	OnFocus func() `json:"-"`

	// The function that is called when the window lose focus.
	OnBlur func() `json:"-"`

	// The function that is called when the window goes full screen.
	OnFullScreen func() `json:"-"`

	// The function that is called when the window exit full screen.
	OnExitFullScreen func() `json:"-"`

	// The function that is called when the window is minimized.
	OnMinimize func() `json:"-"`

	// The function that is called when the window is deminimized.
	OnDeminimize func() `json:"-"`

	// The function that is called when the window is closed.
	// Returning bool prevents the window to be closed.
	OnClose func() bool `json:"-"`
}

// MacWindowConfig is a struct that describes window fields specific to MacOS.
type MacWindowConfig struct {
	BackgroundVibrancy Vibrancy
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

type windowWithLogs struct {
	Window
}

func (w *windowWithLogs) Load(url string, v ...interface{}) error {
	parsedURL := fmt.Sprintf(url, v...)

	WhenDebug(func() {
		Debug("window %s is loading %s",
			w.ID(),
			parsedURL,
		)
	})

	err := w.Window.Load(url, v...)
	if err != nil {
		Log("window %s failed to load %s: %s",
			w.ID(),
			parsedURL,
			err,
		)
	}
	return err
}

func (w *windowWithLogs) Render(c Compo) error {
	WhenDebug(func() {
		Debug("window %s is rendering %T",
			w.ID(),
			c,
		)
	})

	err := w.Window.Render(c)
	if err != nil {
		Log("window %s failed to render %T: %s",
			w.ID(),
			c,
			err,
		)
	}
	return err
}

func (w *windowWithLogs) Reload() error {
	WhenDebug(func() {
		Debug("window %s is reloading", w.ID())
	})

	err := w.Window.Reload()
	if err != nil {
		Log("window %s failed to reload: %s",
			w.ID(),
			err,
		)
	}
	return err
}

func (w *windowWithLogs) Previous() error {
	WhenDebug(func() {
		Debug("window %s is loading previous", w.ID())
	})

	err := w.Window.Previous()
	if err != nil {
		Log("window %s failed to load previous: %s",
			w.ID(),
			err,
		)
	}
	return err
}

func (w *windowWithLogs) Next() error {
	WhenDebug(func() {
		Debug("window %s is loading next", w.ID())
	})

	err := w.Window.Next()
	if err != nil {
		Log("window %s failed to load next: %s",
			w.ID(),
			err,
		)
	}
	return err
}

func (w *windowWithLogs) Close() error {
	WhenDebug(func() {
		Debug("window %s is closing", w.ID())
	})

	err := w.Window.Close()
	if err != nil {
		Log("window %s failed to close: %s",
			w.ID(),
			err,
		)
	}
	return err
}

func (w *windowWithLogs) Move(x, y float64) {
	WhenDebug(func() {
		Debug("window %s is moving to x:%.2f y:%.2f",
			w.ID(),
			x,
			y,
		)
	})

	w.Window.Move(x, y)
}

func (w *windowWithLogs) Center() {
	WhenDebug(func() {
		Debug("window %s is moving to center", w.ID())
	})

	w.Window.Center()
}

func (w *windowWithLogs) Resize(width, height float64) {
	WhenDebug(func() {
		Debug("window %s is resizing to width:%.2f height:%.2f",
			w.ID(),
			width,
			height,
		)
	})

	w.Window.Resize(width, height)
}

func (w *windowWithLogs) Focus() {
	WhenDebug(func() {
		Debug("window %s is getting focus", w.ID())
	})

	w.Window.Focus()
}

func (w *windowWithLogs) FullScreen() {
	WhenDebug(func() {
		Debug("window %s is entering full screen", w.ID())
	})

	w.Window.FullScreen()
}

func (w *windowWithLogs) ExitFullScreen() {
	WhenDebug(func() {
		Debug("window %s is exiting full screen", w.ID())
	})

	w.Window.ExitFullScreen()
}

func (w *windowWithLogs) Minimize() {
	WhenDebug(func() {
		Debug("window %s is minimizing", w.ID())
	})

	w.Window.Minimize()
}

func (w *windowWithLogs) Deminimze() {
	WhenDebug(func() {
		Debug("window %s is deminimizing", w.ID())
	})

	w.Window.Deminimize()
}
