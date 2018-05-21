package app

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

	// ToggleFullScreen takes the window into or out of fullscreen mode.
	ToggleFullScreen()

	// Minimize takes the window into or out of minimized mode
	ToggleMinimize()
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
	FixedSize       bool
	CloseHidden     bool
	MinimizeHidden  bool
	TitlebarHidden  bool
	DefaultURL      string
	Mac             MacWindowConfig

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
