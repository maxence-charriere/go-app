package app

// Backend is the interface that describes an app backend.
// Backends define what is availabe for a given platform.
type Backend interface {
	// Run runs the backend util the app is closed.
	Run(f Factory, uiChan chan func()) error

	// Render the given value.
	// It is up to the backend implementation to decide whether it render a
	// value.
	Render(v interface{}) error
}

// AddBackend adds a backend to the app.
// Should be used only in backend implementations.
// Must be called once at package initialization.
func AddBackend(b Backend) {
	backends = append(backends, b)
}
