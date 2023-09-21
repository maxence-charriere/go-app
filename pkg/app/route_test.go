package app

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type routeCompo struct {
	Compo
}

type routeWithRegexpCompo struct {
	Compo
}

func TestRoutes(t *testing.T) {
	utests := []struct {
		scenario     string
		createRoutes func(*router)
		path         string
		expected     Composer
		notFound     bool
	}{
		{
			scenario: "path is not routed",
			path:     "/goodbye",
			notFound: true,
		},
		{
			scenario: "empty path is not routed",
			path:     "",
			notFound: true,
		},
		{
			scenario: "path is routed",
			createRoutes: func(r *router) {
				r.route("/a", NewComponentFunc(&routeCompo{}))
			},
			expected: &routeCompo{},
			path:     "/a",
		},
		{
			scenario: "path take priority over pattern",
			path:     "/abc",
			createRoutes: func(r *router) {
				r.route("/abc", NewComponentFunc(&routeCompo{}))
				r.routeWithRegexp("^/a.*$", NewComponentFunc(&routeWithRegexpCompo{}))
			},
			expected: &routeCompo{},
		},
		{
			scenario: "pattern is routed",
			path:     "/ab",
			createRoutes: func(r *router) {
				r.route("/abc", NewComponentFunc(&routeCompo{}))
				r.routeWithRegexp("^/a.*$", NewComponentFunc(&routeWithRegexpCompo{}))
			},
			expected: &routeWithRegexpCompo{},
		},
		{
			scenario: "pattern with inner wildcard is routed",
			path:     "/user/42/settings",
			createRoutes: func(r *router) {
				r.routeWithRegexp("^/user/.*/settings$", NewComponentFunc(&routeWithRegexpCompo{}))
			},
			expected: &routeWithRegexpCompo{},
		},
		{
			scenario: "not matching pattern with inner wildcard is not routed",
			path:     "/user/42/settings/",
			createRoutes: func(r *router) {
				r.routeWithRegexp("^/user/.*/settings$", NewComponentFunc(&routeWithRegexpCompo{}))
			},
			notFound: true,
		},
		{
			scenario: "pattern with end wildcard is routed",
			path:     "/user/1001/files/foo/bar/baz.png",
			createRoutes: func(r *router) {
				r.routeWithRegexp("^/user/.*/files/.*$", NewComponentFunc(&routeWithRegexpCompo{}))
			},
			expected: &routeWithRegexpCompo{},
		},
		{
			scenario: "not matching pattern with end wildcard is not routed",
			path:     "/user/1001/files",
			createRoutes: func(r *router) {
				r.routeWithRegexp("^/user/.*/files/.*$", NewComponentFunc(&routeWithRegexpCompo{}))
			},
			notFound: true,
		},
		{
			scenario: "pattern with OR condition is routed",
			path:     "/color/red",
			createRoutes: func(r *router) {
				r.routeWithRegexp("^/color/(red|green|blue)$", NewComponentFunc(&routeWithRegexpCompo{}))
			},
			expected: &routeWithRegexpCompo{},
		},
		{
			scenario: "not matching pattern with OR condition is not routed",
			path:     "/color/fuschia",
			createRoutes: func(r *router) {
				r.routeWithRegexp("^/color/(red|green|blue)$", NewComponentFunc(&routeWithRegexpCompo{}))
			},
			notFound: true,
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			r := makeRouter()
			if u.createRoutes != nil {
				u.createRoutes(&r)
			}

			compo, isRouted := r.createComponent(u.path)
			if u.notFound {
				require.Nil(t, compo)
				require.False(t, isRouted)
				return
			}
			require.True(t, isRouted)
			require.NotNil(t, compo)
			require.Equal(t, reflect.TypeOf(u.expected), reflect.TypeOf(compo))
		})
	}
}
