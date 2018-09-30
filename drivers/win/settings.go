package win

// Settings represents the app settings.
type Settings struct {
	// The app name.
	Name string

	// The app description.
	Description string

	// The app publisher.
	Publisher string

	// The URL scheme to call the app.
	Scheme string

	// The app icon path relative to the resources directory as .png file.
	// Provide a big one! Other required icon sizes will be auto generated.
	Icon string
}
