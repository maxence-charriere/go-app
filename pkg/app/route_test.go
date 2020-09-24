package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type routeCompo struct {
	Compo
	id string
}

func TestRoutes(t *testing.T) {
	routes := router{}
	routes.route("/a", &routeCompo{id: "a"})
	routes.route("/abc", &routeCompo{id: "abc"})
	routes.routeWithRegexp("^/a.*$", &routeCompo{id: "a-star"})
	routes.routeWithRegexp("^/user/.*/settings$", &routeCompo{id: "settings"})
	routes.routeWithRegexp("^/user/.*/files/.*$", &routeCompo{id: "files"})
	routes.routeWithRegexp("^/color/(red|green|blue)$", &routeCompo{id: "color"})

	tests := []struct {
		scenario   string
		path       string
		expectedID string
		notFound   bool
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
			scenario:   "path is routed",
			path:       "/a",
			expectedID: "a",
		},
		{
			scenario:   "path take priority over pattern",
			path:       "/abc",
			expectedID: "abc",
		},
		{
			scenario:   "pattern is routed",
			path:       "/ab",
			expectedID: "a-star",
		},
		{
			scenario:   "pattern with inner wildcard is routed",
			path:       "/user/42/settings",
			expectedID: "settings",
		},
		{
			scenario: "not matching pattern with inner wildcard is not routed",
			path:     "/user/42/settings/",
			notFound: true,
		},
		{
			scenario:   "pattern with end wildcard is routed",
			path:       "/user/1001/files/foo/bar/baz.png",
			expectedID: "files",
		},
		{
			scenario: "not matching pattern with end wildcard is not routed",
			path:     "/user/1001/files",
			notFound: true,
		},
		{
			scenario:   "pattern with OR condition is routed",
			path:       "/color/red",
			expectedID: "color",
		},
		{
			scenario: "not matching pattern with OR condition is not routed",
			path:     "/color/fuschia",
			notFound: true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			node, routed := routes.ui(test.path)

			if test.notFound {
				require.False(t, routed, "node is routed")
				return
			}

			id := node.(*routeCompo).id
			require.Equal(t, test.expectedID, id)
		})
	}
}

type routeSubString struct {
	Compo
	URL string `app:"0"`
	ID  string `app:"1"`
}

func TestRoutesSubstring(t *testing.T) {
	routes := router{}
	routes.routeWithRegexp(`^/url/(\w*)`, &routeSubString{})
	tests := []struct {
		scenario   string
		path       string
		expectedID string
		notFound   bool
	}{
		{
			scenario:   "matching pattern parser substring to apptag",
			path:       "/url/11223344",
			expectedID: "11223344",
		},
	}
	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			node, routed := routes.ui(test.path)

			if test.notFound {
				require.False(t, routed, "node is routed")
				return
			}

			id := node.(*routeSubString).ID
			require.Equal(t, test.expectedID, id)
		})
	}
}
