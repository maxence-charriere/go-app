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

// Windower represents a context with window specific interactions.
type Windower interface {
	Contexter

	Position() (x float64, y float64)

	Move(x float64, y float64)

	Size() (width float64, height float64)

	Resize(width float64, height float64)
}

// Window represents a window.
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
