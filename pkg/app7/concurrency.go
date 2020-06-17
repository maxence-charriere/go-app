package app

var (
	dispatch dispatcher = Dispatch
	uiChan              = make(chan func(), 512)
)

// dispatcher is a function that executes the given function on the goroutine
// dedicated to UI.
type dispatcher func(func())

// Dispatch executes the given function on the UI goroutine.
func Dispatch(f func()) {
	uiChan <- f
}
