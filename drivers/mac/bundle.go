// +build darwin,amd64

package mac

// Bundle is the struct that describes the app bundle.
// It is used to set .app variables and capabilities.
// Only applied when the app is build with goapp mac build -bundle.
type Bundle struct {
	// The app name and menu bar/dock display name.
	AppName string

	// The UTI representing the app.
	ID string

	// The URL scheme that launches the app."
	URLScheme string

	// The version of the app (minified form eg 1.42).
	Version string

	// The build number.
	BuildNumber int

	// The app icon path relative to the resources directory as .png file.
	// Provide a big one! Other required icon sizes will be auto generated.
	Icon string

	// The development region.
	DevRegion string

	// A human readable copyright.
	Copyright string

	// The application role.
	Role Role

	// The application category.
	Category Category

	// Reports wheter the app runs in background mode.
	// Background apps does not appear in the dock and menubar.
	Background bool

	// Reports whether the app is a server (accepts incoming connections).
	Server bool

	// Reports whether the app uses the camera.
	Camera bool

	// Reports whether the app uses the microphone.
	Microphone bool

	// Reports whether the app uses the USB devices.
	USB bool

	// Reports whether the app uses printers.
	Printers bool

	// Reports whether the app uses bluetooth.
	Bluetooth bool

	// Reports whether the app has access to contacts.
	Contacts bool

	// Reports whether the app has access to device location.
	Location bool

	// Reports whether the app has access to calendars.
	Calendar bool

	// The file picker access mode.
	FilePickers FileAccess

	// The Download directory access mode.
	Downloads FileAccess

	// The Pictures directory access mode.
	Pictures FileAccess

	// The Music directory access mode.
	Music FileAccess

	// The Movies directory access mode.
	Movies FileAccess

	// The UTIs representing the file types that the app can open.
	SupportedFiles []string
}
