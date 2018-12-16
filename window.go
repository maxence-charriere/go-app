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

	// The function that is called when the window is moved.
	OnMove func(w Window) `json:"-"`

	// The function that is called when the window is resized.
	OnResize func(w Window) `json:"-"`

	// The function that is called when the window get focus.
	OnFocus func(w Window) `json:"-"`

	// The function that is called when the window lose focus.
	OnBlur func(w Window) `json:"-"`

	// The function that is called when the window goes full screen.
	OnFullScreen func(w Window) `json:"-"`

	// The function that is called when the window exit full screen.
	OnExitFullScreen func(w Window) `json:"-"`

	// The function that is called when the window is minimized.
	OnMinimize func(w Window) `json:"-"`

	// The function that is called when the window is deminimized.
	OnDeminimize func(w Window) `json:"-"`

	// The function that is called when the window is closed.
	OnClose func(w Window) `json:"-"`
}
