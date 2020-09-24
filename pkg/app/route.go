package app

import (
	"reflect"
	"regexp"
	"strconv"
)

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
		subList := r.regexp.FindStringSubmatch(path)
		if len(subList) > 0 {
			v := reflect.ValueOf(r.node)
			if v.Kind() == reflect.Ptr {
				v := v.Elem()
				if v.Kind() == reflect.Struct {
					for i := 0; i < v.NumField(); i++ {
						t := v.Type().Field(i)
						tag := t.Tag.Get("app")
						if tag == "" {
							continue
						}
						idx, err := strconv.Atoi(tag)
						if err != nil {
							continue
						}
						if idx < len(subList) {
							vv := v.Field(i)
							if vv.Kind() == reflect.String {
								vv.SetString(subList[idx])
							}
						}

					}
				}
			}

			return r.node, true
		}
	}

	return nil, false
}

type regexpRoute struct {
	regexp *regexp.Regexp
	node   UI
}
