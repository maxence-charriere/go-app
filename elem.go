package app

// Elem is the interface that describes an app element.
type Elem interface {
	// ID returns the element identifier.
	ID() string

	// Contains reports whether the component is mounted in the element.
	Contains(Compo) bool

	// WhenPage calls the given func when the element is a view.
	WhenView(func(View))

	// WhenWindow calls the given func when the element is a window.
	WhenWindow(func(Window))

	// WhenStatusMenu calls the given func when the element is a menu.
	WhenMenu(func(Menu))

	// WhenDockTile calls the given func when the element is a dock tile.
	WhenDockTile(func(DockTile))

	// WhenStatusMenu calls the given func when the element is a status menu.
	WhenStatusMenu(func(StatusMenu))

	// WhenErr call the given func when the element is in an error state.
	WhenErr(func(err error))

	// Err returns the error that prevent the element to work.
	Err() error
}

// View is the interface that describe an element that hosts components.
type View interface {
	Elem

	// Load loads the page specified by the URL.
	// URL can be formated as fmt package functions.
	// Calls with an URL which contains a component name will load the named
	// component.
	// e.g. hello will load the component named hello.
	// It returns an error if the component is not imported.
	Load(url string, v ...interface{})

	// Reload reloads the current page.
	Reload()

	// CanPrevious reports whether load the previous page is possible.
	CanPrevious() bool

	// Previous loads the previous page.
	Previous()

	// CanNext indicates if loading next page is possible.
	CanNext() bool

	// Next loads the next page.
	Next()

	// Compo returns the loaded component.
	Compo() Compo

	// Render renders the component.
	Render(Compo)
}

// Closer is the interface that describes an element that can be closed.
type Closer interface {
	// Close closes the element and free its allocated resources.
	Close()
}

// NotificationConfig is a struct that describes a notification.
type NotificationConfig struct {
	Title     string
	Subtitle  string
	Text      string
	ImageName string
	Sound     bool

	OnReply func(reply string) `json:"-"`
}
