package app

import (
	"net/url"
	"testing"
)

const (
	// NOT_FOUND is returned from drill when no route matched the URL
	NOT_FOUND = "_NOT_FOUND_"
	// DOM_ERROR is returned from drill if the DOM had an unepected configuration
	DOM_ERROR = "_DOM_ERROR_"
)

// testComp is a very simple test component that renders as a text node containing its name
type testComp struct {
	Compo
	name string
}

func (c *testComp) Render() UI {
	return Div().Body(Text(c.name))
}

// TestRoute tests the basic routing tables
func TestRoute(t *testing.T) {

	resetRoutes()
	Route("/a", &testComp{name: "a"})
	Route("/abc", &testComp{name: "abc"})
	Route("/azz", &testComp{name: "azz"})
	Route("/xyz", &testComp{name: "xyz"})

	navExpect(t, "/abc", "abc")
	navExpect(t, "/a", "a")
	navExpect(t, "/x", NOT_FOUND)
	navExpect(t, "", NOT_FOUND)
	navExpect(t, "/abc?foo=bar", "abc") // query is not used in route lookup
}

// TestRouteRe tests some basic regex patterns and priority of Route vs RouteRe
func TestRouteRe(t *testing.T) {

	resetRoutes()
	Route("/a", &testComp{name: "a"})
	Route("/abc", &testComp{name: "abc"})
	Route("/123", &testComp{name: "123"})

	RouteRe("/x.*", &testComp{name: "x-star"})
	RouteRe("/a.*", &testComp{name: "a-star"}) // test match when not first regex route

	navExpect(t, "/a", "a") // specific routes should be found before regex
	navExpect(t, "/ab", "a-star")
	navExpect(t, "/aa", "a-star")
	navExpect(t, "/abc", "abc")
	navExpect(t, "/.*b", NOT_FOUND)
	navExpect(t, "/.*c", NOT_FOUND)

	navExpect(t, "/123", "123")
	// these two test cases confirm that regex needs full string, not just subset
	navExpect(t, "/4/123", NOT_FOUND)
	navExpect(t, "/123/123", NOT_FOUND)
}

// TestRouteReDirs1 tests regex paths with inner wildcard
func TestRouteReDirs1(t *testing.T) {

	resetRoutes()
	RouteRe("/user/.*/settings", &testComp{name: "settings"})

	navExpect(t, "/user/1001/settings", "settings")
	navExpect(t, "/user/1/settings", "settings")
	navExpect(t, "/user/1001/settings/", NOT_FOUND) // extra trailing slash
}

// TestRouteReDirs2 test regex wildcard at end and middle
func TestRouteReDirs2(t *testing.T) {

	resetRoutes()
	RouteRe("/user/.*/files/.*", &testComp{name: "files"})

	navExpect(t, "/user/1001/files/", "files")
	navExpect(t, "/user/1001/files/index.html", "files")
	navExpect(t, "/user/1001/files/foo/bar/baz.png", "files")
	navExpect(t, "/user/team/green/files/index.html", "files")
}

// TestRouteReDirs3 tests regex with "OR" condition
func TestRouteReDirs3(t *testing.T) {

	resetRoutes()
	RouteRe("/color/(red|green|blue)", &testComp{name: "color"})

	navExpect(t, "/color/red", "color")
	navExpect(t, "/color/blue", "color")
	navExpect(t, "/color/fuschia", NOT_FOUND)
}

// navExpect does the test case that simulates navigating to a new url via route tables.
// Parameters are the url input and the text string to look for in the new DOM.
func navExpect(t *testing.T, path string, result string) {

	resetDOM(t)
	u, err := url.Parse(path)
	if err != nil {
		t.Fatalf("invalid url path '%s': %v", path, err)
	}
	if err = navigate(u, false); err != nil {
		if result == NOT_FOUND {
			t.Logf("Pass: route path: '%s' expected '%s'", path, result)
			return
		}
		t.Fatalf("FAIL: nav %s: %v", path, err)
	}
	if val := drill(t, body); val != result {
		t.Errorf("FAIL: route path:'%s' expect:'%s' actual:'%s'\n",
			path, result, val)
		return
	}
	t.Logf("Pass: route path:'%s' expect:'%s'\n",
		path, result)
}

// resetRoutes empties routing tables, and is called at beginning of each test function
func resetRoutes() {
	routes = make(map[string]UI)
	routesRe = make([]regexRoute, 0)
}

// resetDOM initializes body and content - should be called before navigate
// so that it's easy to determine what changed
func resetDOM(t *testing.T) {
	body = Body()
	content = Div()
	initContent()
}

// drill traverses down DOM tree to leaf text node and return its contents
func drill(t *testing.T, node UI) string {
	if tc, ok := node.(*testComp); ok {
		node = tc.root
	}
	if tx, isText := node.(*text); isText {
		return tx.text()
	}
	if s, ok := node.(standardNode); ok {
		if len(s.children()) > 1 {
			t.Errorf("FAIL: (drill) unexpected numChildren(%d) for %v\n", len(s.children()), node)
			return DOM_ERROR
		}
		return drill(t, s.children()[0])
	}
	if _, ok := node.(*notFound); ok {
		return NOT_FOUND
	}
	t.Errorf("FAIL: (drill) unexpected node %v type %v\n", node, node.nodeType())
	return DOM_ERROR
}
