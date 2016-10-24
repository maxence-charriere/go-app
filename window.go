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

// Window represents a window.
type Window struct {
	Title           string
	Width           uint
	Height          uint
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
func NewWindow(w Window) (ctx Contexter) {
	return driver.NewContext(w)
}
