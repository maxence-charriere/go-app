// +build windows

package win

// Settings represents the app settings.
type Settings struct {
	// The app name.
	Name string

	// The app id.
	ID string

	// The app description.
	Description string

	// The app publisher.
	Publisher string

	// The URL scheme that activate the app.
	URLScheme string

	// The app icon path relative to the resources directory as .png file.
	// Provide a big one! Other required icon sizes will be auto generated.
	Icon string

	// The file types that can be opened by the app.
	SupportedFiles []FileType
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
