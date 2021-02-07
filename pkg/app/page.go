package app

import "net/url"

// PageInfo contains the page info that is modifiable when a page is pre
// rendered.
type PageInfo struct {
	// The page authors.
	Author string

	// The page description.
	Description string

	// The page keywords.
	Keywords []string

	// The text displayed while loading a page.
	LoadingLabel string

	// The page title.
	Title string

	url *url.URL
}

// URL return the page URL.
func (i *PageInfo) URL() *url.URL {
	return i.url
}
