package app

// Constants to specify the vibrancy to be used.
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

// Windower is the interface that describes a window context.
type Windower interface {
	Contexter

	// Position returns the position of the window.
	Position() (x float64, y float64)

	Move(x float64, y float64)

	// Size returns the size of the window.
	Size() (width float64, height float64)

	// Resize resizes the window.
	Resize(width float64, height float64)

	// Close close`s the window.
	//
	// Driver implementation:
	// - Close should call the native way to close a window.
	// - Native windows often have a handler that is called before a window
	//   is destroyed. This handler should be implemented and call
	//   Elements().Remove() to free resources allocated on go side.
	//   markup.Dismount() should be also called to release the components
	//   mounted.
	Close()
}

// Window is a struct that describes a window.
// It will be used by a driver to create a context on the top of a native
// window.
type Window struct {
	Title           string
	Lang            string
	X               float64
	Y               float64
	Width           float64
	Height          float64
	MinWidth        float64
	MinHeight       float64
	MaxWidth        float64
	MaxHeight       float64
	BackgroundColor string
	Vibrancy        Vibrancy
	Borderless      bool
	FixedSize       bool
	CloseHidden     bool
	MinimizeHidden  bool
	TitlebarHidden  bool

	OnMinimize       func()
	OnDeminimize     func()
	OnFullScreen     func()
	OnExitFullScreen func()
	OnMove           func(x float64, y float64)
	OnResize         func(width float64, height float64)
	OnFocus          func()
	OnBlur           func()
	OnClose          func() bool
}

// Vibrancy represents the NSVisualEffectView which will be applied to the
// background of the window.
// Only applicable on Apple apps.
// When set, BackgroundColor is ignored.
type Vibrancy uint8

// NewWindow creates a new window.
func NewWindow(w Window) Windower {
	return driver.NewElement(w).(Windower)
}
