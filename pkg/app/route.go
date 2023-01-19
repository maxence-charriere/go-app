package app

import (
	"reflect"
	"regexp"
	"sync"
)

var (
	routes = makeRouter()
)

// Route set the type of component to be mounted when a page is navigated to the
// given path.
func Route(path string, c Composer) {
	RouteFunc(path, newZeroComponentFunc(c))
}

// RouteWithRegexp set the type of component to be mounted when a page is
// navigated to a path that matches the given pattern.
func RouteWithRegexp(pattern string, c Composer) {
	RouteWithRegexpFunc(pattern, newZeroComponentFunc(c))

}

// RouteFunc set a function that creates the component to be mounted when a page
// is navigated to the given path.
func RouteFunc(path string, newComponent func() Composer) {
	routes.route(path, newComponent)
}

// RouteWithRegexpFunc set a function that creates the component to be mounted
// when a page is navigated to a path that matches the given pattern.
func RouteWithRegexpFunc(pattern string, newComponent func() Composer) {
	routes.routeWithRegexp(pattern, newComponent)
}

func newZeroComponentFunc(c Composer) func() Composer {
	componentType := reflect.TypeOf(c)

	return func() Composer {
		return reflect.New(componentType.Elem()).Interface().(Composer)
	}
}

type router struct {
	mu               sync.RWMutex
	routes           map[string]func() Composer
	routesWithRegexp []regexpRoute
}

func makeRouter() router {
	return router{
		routes: make(map[string]func() Composer),
	}
}

func (r *router) route(path string, newComponent func() Composer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routes[path] = newComponent
}

func (r *router) routeWithRegexp(pattern string, newComponent func() Composer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routesWithRegexp = append(r.routesWithRegexp, regexpRoute{
		regexp:       regexp.MustCompile(pattern),
		newComponent: newComponent,
	})
}

func (r *router) createComponent(path string) (Composer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	newComponent, isRouted := r.routes[path]
	if !isRouted {
		for _, rwr := range r.routesWithRegexp {
			if rwr.regexp.MatchString(path) {
				newComponent = rwr.newComponent
				isRouted = true
				break
			}
		}
	}
	if !isRouted {
		return nil, false
	}

	return newComponent(), true
}

type regexpRoute struct {
	regexp       *regexp.Regexp
	newComponent func() Composer
}
