package app

var (
	driver Driver
)

// Driver is the interface that describes the implementation to handle platform
// specific rendering.
type Driver interface {
	// Run runs the application.
	//
	// Driver implementation:
	// - Should start the app loop.
	Run()

	// Resources returns the location of the resources directory.
	Resources() string

	// Storage returns the location of the app storage directory.
	Storage() string

	// JavascriptBridge is the javascript function to call when a driver want to
	// pass data to the native platform.
	JavascriptBridge() string
}
