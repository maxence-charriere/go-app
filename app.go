package app

import (
	"net/url"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
)

var (
	// OnLaunch is a handler which (if set) is called when the app is
	// initialized and ready.
	// The main window should be created here.
	OnLaunch func()

	// OnFocus is a handler which (if set) is called when the app became
	// focused.
	OnFocus func()

	// OnBlur is a handler which (if set) is called when the app lost the
	// focus.
	OnBlur func()

	// OnReopen is a handler which (if set) is called when the app is reopened.
	// Eg. when the dock icon is clicked.
	OnReopen func()

	// OnFilesOpen is a handler which (if set) is called when files are targeted
	// to be opened with the app.
	OnFilesOpen func(filenames []string)

	// OnURLOpen is a handler which (if set) is called when the app is opened by
	// an URL.
	OnURLOpen func(URL url.URL)

	// OnTerminate is a handler which (if set) is called when the app is
	// requested to terminates. Return false cancels the termination request.
	OnTerminate func() bool

	// OnFinalize is a handler which (if set) is called when the app is about
	// to be terminated.
	// It should be used to perform any final cleanup before the application
	// terminates.
	OnFinalize func()
)

// Run runs the app.
func Run() {
	driver.Run()
}

// Render renders a component. Update the rendering of c.
// c must be mounted into a context.
func Render(c Componer) {
	syncs, err := markup.Synchronize(c)
	if err != nil {
		log.Error(err)
		return
	}

	ctx := Context(c)
	for _, s := range syncs {
		ctx.Render(s)
	}
}

// Resources returns the location of the resources directory.
// resources directory should contain files required by the UI.
// Its path should be used only for read only operations, otherwise it could
// mess up with the app signature.
func Resources() string {
	return driver.Resources()
}

// Storage returns the location of the app storage directory.
// Content generated (e.g. sqlite db) or downloaded (e.g. images, music)
// should be saved in this directory.
func Storage() string {
	return driver.Storage()
}

// MenuBar returns the menu bar context.
// ok will be false if there is no menubar available.
func MenuBar() (menu Contexter, ok bool) {
	return driver.MenuBar()
}

// Dock returns the dock context.
// ok will be false if there is no dock available.
func Dock() (d Docker, ok bool) {
	return driver.Dock()
}
