package win

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
)

// Driver is the app.Driver implementation for Windows.
type Driver struct {
	core.Driver `json:"-"`

	// The URL of the component to load in the default window. It overrides the
	// DefaultWindow.URL value.
	URL string `json:"-"`

	// The default window configuration.
	DefaultWindow app.WindowConfig `json:"-"`

	// The app name.
	//
	// It is used only for goapp packaging.
	Name string `json:",omitempty"`

	// The app id.
	//
	// It is used only for goapp packaging.
	ID string `json:",omitempty"`

	// The app description.
	//
	// It is used only for goapp packaging.
	Description string `json:",omitempty"`

	// The app publisher.
	//
	// It is used only for goapp packaging.
	Publisher string `json:",omitempty"`

	// The URL scheme that launches the app.
	//
	// It is used only for goapp packaging.
	URLScheme string `json:",omitempty"`

	// The app icon path relative to the resources directory. It must be a
	// ".png". Provide a big one! Other required icon sizes will be auto
	// generated.
	//
	// It is used only for goapp packaging.
	Icon string `json:",omitempty"`

	// The file types that can be opened by the app.
	//
	// It is used only for goapp packaging.
	SupportedFiles []FileType `json:",omitempty"`

	ui      chan func()
	factory *app.Factory
	events  *app.EventRegistry
	elems   *core.ElemDB
	winRPC  *bridge.PlatformRPC
	goRPC   *bridge.GoRPC
	stop    func()
}

// Target satisfies the app.Driver interface.
func (d *Driver) Target() string {
	return "windows"
}

// FileType describes a file type that can be opened by the app.
type FileType struct {
	// The  name.
	// Must be non empty, a single word and lowercased.
	Name string

	// The help that appears when the user hovers on the icon for a file of this
	// type.
	Help string

	// The path of the icon that is used to identify the file type on the
	// desktop and in the Set Default Programs on the Control Panel.
	// If no icon is specified, the application icon is used.
	//
	// Must be relative to the resource directory and be a .png file.
	Icon string

	// The associated extensions.
	// Must contain at least 1 element.
	Extensions []FileExtension
}

// FileExtension describes a file extension.
type FileExtension struct {
	// The extension.
	// Eg. ".png"
	// Must be non empty.
	Ext string

	// The mime type.
	Mime string
}
