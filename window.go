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

	// Reports whether the window is the focus.
	IsFocus() bool

	// FullScreen takes the window into full screen mode.
	FullScreen()

	// ExitFullScreen takes the window out of full screen mode.
	ExitFullScreen()

	// Reports whether the window is in full screen mode.
	IsFullScreen() bool

	// Minimize takes the window into minimized mode.
	Minimize()

	// Deminimize takes the window out of minimized mode.
	Deminimize()

	// Reports whether the window is minimized.
	IsMinimized() bool
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

	// Enables frosted effect.
	FrostedBackground bool

	// Reports whether the window is resizable.
	FixedSize bool

	// Reports whether the close button is hidden.
	CloseHidden bool

	// Reports whether the minimize button is hidden.
	MinimizeHidden bool

	// The URL of the component to load when the window is created.
	URL string
}

const (
	// WindowMoved is the event emitted when a window is moved.
	WindowMoved Event = "window-moved"

	// WindowResized is the event emitted when a window is resized.
	WindowResized Event = "window-resized"

	// WindowFocused is the event emitted when a window gets focus.
	WindowFocused Event = "window-focused"

	// WindowBlurred is the event emitted when a window loses focus.
	WindowBlurred Event = "window-blurred"

	// WindowEnteredFullScreen is the event emitted when a window goes full screen.
	WindowEnteredFullScreen Event = "window-entered-fullscreen"

	// WindowExitedFullScreen is the event emitted when a window exits full screen.
	WindowExitedFullScreen Event = "window-exited-fullscreen"

	// WindowMinimized is the event emitted when a window is minimized.
	WindowMinimized Event = "window-minimized"

	// WindowDeminimized is the event emitted when a window is deminimized.
	WindowDeminimized Event = "window-deminimized"

	// WindowClosed is the event emitted when a window is closed.
	WindowClosed Event = "window-closed"
)
