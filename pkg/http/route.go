package http

import (
	"net/http"
)

// RouteHandler is the interface that describe a routable http.Handler.
type RouteHandler interface {
	http.Handler

	// Handle reports whether the handler can handle the given request.
	CanHandle(r *http.Request) bool
}

// Route returns a http.Handler that routes requests to the appropriate given
// handler.
func Route(handlers ...RouteHandler) http.Handler {
	return route(handlers)
}

type route []RouteHandler

func (r route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, h := range r {
		if h.CanHandle(req) {
			h.ServeHTTP(w, req)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
