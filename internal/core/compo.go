package core

import (
	"net/url"
	"strings"
)

// CompoNameFromURL returns the component name targeted by the given URL.
func CompoNameFromURL(u *url.URL) string {
	if len(u.Scheme) != 0 && u.Scheme != "compo" {
		return ""
	}

	p := u.Path
	p = strings.TrimPrefix(p, "/")

	path := strings.SplitN(p, "/", 2)
	if len(path[0]) == 0 {
		return ""
	}

	names := strings.SplitN(path[0], "?", 2)
	name := names[0]
	name = strings.ToLower(name)
	name = strings.TrimPrefix(name, "main.")

	return name
}

// CompoNameFromURLString returns the component name targeted by the given URL
// string.
func CompoNameFromURLString(rawurl string) string {
	u, _ := url.Parse(rawurl)
	return CompoNameFromURL(u)
}
