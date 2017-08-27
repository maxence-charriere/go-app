package url

import (
	"net/url"
	"strings"
)

// A URL represents a parsed URL (technically, a URI reference).
// It is a wrapper of the https://golang.org/pkg/net/url package URL struct.
type URL struct {
	url.URL
	componentName string
}

// Component reports if the URL targets a component and returns the component
// name if it's the case.
func (u *URL) Component() (name string, ok bool) {
	name = u.componentName
	ok = len(u.componentName) > 0
	return
}

// Parse parses rawurl into a URL structure.
// The rawurl may be relative or absolute.
func Parse(rawurl string) (u URL, err error) {
	var rawURL *url.URL
	if rawURL, err = url.Parse(rawurl); err != nil {
		return
	}
	u.URL = *rawURL

	if len(u.Scheme) == 0 {
		u.Scheme = "component"
	}

	if u.Scheme == "component" {
		u.componentName = strings.TrimLeft(u.Path, "/")
	}
	return
}
