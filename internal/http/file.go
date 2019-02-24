package http

import "net/http"

// FileHandler returns a handler that serves files located in the web directory.
func FileHandler(webDir string) http.Handler {
	return http.FileServer(http.Dir(webDir))
}
