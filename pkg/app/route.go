package app

import (
	"reflect"
	"regexp"
	"sync"
)

var (
	routes = makeRouter()
)

// Route associates the type of the given component to the given path.
//
// When a page is requested and matches the route, a new instance of the given
// component is created before being displayed.
func Route(path string, c Composer) {
	routes.route(path, c)
}

// RouteWithRegexp associates the type of the given component to the given
// regular expression pattern.
//
// Patterns use the Go standard regexp format.
//
// When a page is requested and matches the pattern, a new instance of the given
// component is created before being displayed.
func RouteWithRegexp(pattern string, c Composer) {
	routes.routeWithRegexp(pattern, c)
}

type router struct {
	mu               sync.RWMutex
	routes           map[string]reflect.Type
	routesWithRegexp []regexpRoute
}

func makeRouter() router {
	return router{
		routes: make(map[string]reflect.Type),
	}
}

func (r *router) route(path string, c Composer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routes[path] = reflect.TypeOf(c)
}

func (r *router) routeWithRegexp(pattern string, c Composer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routesWithRegexp = append(r.routesWithRegexp, regexpRoute{
		regexp:    regexp.MustCompile(pattern),
		compoType: reflect.TypeOf(c),
	})
}

func (r *router) createComponent(path string) (Composer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	compoType, isRouted := r.routes[path]
	if !isRouted {
		for _, rwr := range r.routesWithRegexp {
			if rwr.regexp.MatchString(path) {
				compoType = rwr.compoType
				isRouted = true
				break
			}
		}
	}
	if !isRouted {
		return nil, false
	}

	compo := reflect.New(compoType.Elem()).Interface().(Composer)
	return compo, true
}

func (r *router) len() int {
	return len(r.routes) + len(r.routesWithRegexp)
}

type regexpRoute struct {
	regexp    *regexp.Regexp
	compoType reflect.Type
}
