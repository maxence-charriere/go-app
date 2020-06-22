package app

import "regexp"

var (
	routes router
)

// Route binds the requested path to the given UI node.
func Route(path string, node UI) {
	routes.route(path, node)
}

// RouteWithRegexp binds the regular expression pattern to the given UI node.
// Patterns use the Go standard regexp format.
func RouteWithRegexp(pattern string, node UI) {
	routes.routeWithRegexp(pattern, node)
}

type router struct {
	routes           map[string]UI
	routesWithRegexp []regexpRoute
}

func (r *router) route(path string, node UI) {
	if r.routes == nil {
		r.routes = make(map[string]UI)
	}
	r.routes[path] = node
}

func (r *router) routeWithRegexp(pattern string, node UI) {
	r.routesWithRegexp = append(r.routesWithRegexp, regexpRoute{
		regexp: regexp.MustCompile(pattern),
		node:   node,
	})
}

func (r *router) ui(path string) (UI, bool) {
	if node, routed := r.routes[path]; routed {
		return node, true
	}

	for _, r := range r.routesWithRegexp {
		if r.regexp.MatchString(path) {
			return r.node, true
		}
	}

	return nil, false
}

type regexpRoute struct {
	regexp *regexp.Regexp
	node   UI
}
