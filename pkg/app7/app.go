//go:generate go run gen/html.go
//go:generate go run gen/scripts.go
//go:generate go fmt

package app

import (
	"strings"
)

var (
	remoteRootDir string
)

// EventHandler represents a function that can handle HTML events. They are
// always called on the UI goroutine.
type EventHandler func(ctx Context, e Event)

// StaticResource makes a static resource path point to the right
// location whether the root directory is remote or not.
//
// Static resources are resources located in the web directory.
//
// This call is used internally to resolve paths within Cite, Data, Href, Src,
// and SrcSet. Paths already resolved are skipped.
func StaticResource(path string) string {
	if !strings.HasPrefix(path, "/web/") &&
		!strings.HasPrefix(path, "web/") ||
		remoteRootDir == "" {
		return path
	}

	path = strings.TrimPrefix(path, "/")

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return remoteRootDir + path
}

// Window returns the JavaScript "window" object.
func Window() BrowserWindow {
	return window
}
