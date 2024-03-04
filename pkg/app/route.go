package app

import (
	"regexp"
	"sync"
)

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

func (r *router) routed(path string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, routed := r.routes[path]; routed {
		return true
	}

	for _, rwr := range r.routesWithRegexp {
		if rwr.regexp.MatchString(path) {
			return true
		}
	}

	return false
}

func (r *router) createComponent(path string) (Composer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if newComponent, routed := r.routes[path]; routed {
		return newComponent(), true
	}

	for _, rwr := range r.routesWithRegexp {
		if rwr.regexp.MatchString(path) {
			return rwr.newComponent(), true
		}
	}

	return nil, false
}

type regexpRoute struct {
	regexp       *regexp.Regexp
	newComponent func() Composer
}
