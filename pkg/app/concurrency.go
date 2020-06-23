package app

var (
	dispatch Dispatcher = Dispatch
	uiChan              = make(chan func(), 512)
)

// Dispatcher is a function that executes the given function on the goroutine
// dedicated to UI.
type Dispatcher func(func())

// Dispatch executes the given function on the UI goroutine.
func Dispatch(f func()) {
	uiChan <- f
}
