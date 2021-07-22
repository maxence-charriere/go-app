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

// RouteFactory associates the function to the given path.
//
// When a page is requested and matches the route, the function is called
// to create a new component to be displayed.
func RouteFactory(path string, cf ComposerFactory) {
	routes.routeFactory(path, cf)
}

// RouteWithRegexp associates the function to the given regular expression
// pattern.
//
// Patterns use the Go standard regexp format.
//
// When a page is requested and matches the pattern, the function is called
// to create a new component to be displayed.
func RouteWithRegexpFactory(pattern string, cf ComposerFactory) {
	routes.routeWithRegexpFactory(pattern, cf)
}

type ComposerFactory func() Composer

type router struct {
	mu                      sync.RWMutex
	routes                  map[string]reflect.Type
	routesWithRegexp        []regexpRoute
	routesFactory           map[string]ComposerFactory
	routesWithRegexpFactory []regexpRouteFactory
}

func makeRouter() router {
	return router{
		routes:        make(map[string]reflect.Type),
		routesFactory: make(map[string]ComposerFactory),
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

func (r *router) routeFactory(path string, cf ComposerFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routesFactory[path] = cf
}

func (r *router) routeWithRegexpFactory(pattern string, cf ComposerFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routesWithRegexpFactory = append(r.routesWithRegexpFactory, regexpRouteFactory{
		regexp:  regexp.MustCompile(pattern),
		factory: cf,
	})
}

func (r *router) createComponent(path string) (Composer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var compo Composer

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
	if isRouted {
		compo = reflect.New(compoType.Elem()).Interface().(Composer)
	} else {
		factory, isRouted := r.routesFactory[path]
		if !isRouted {
			for _, rwr := range r.routesWithRegexpFactory {
				if rwr.regexp.MatchString(path) {
					factory = rwr.factory
					isRouted = true
					break
				}
			}
		}
		if !isRouted {
			return nil, false
		}

		compo = factory()
	}

	return compo, true
}

func (r *router) len() int {
	return len(r.routes) + len(r.routesWithRegexp)
}

type regexpRoute struct {
	regexp    *regexp.Regexp
	compoType reflect.Type
}

type regexpRouteFactory struct {
	regexp  *regexp.Regexp
	factory ComposerFactory
}
