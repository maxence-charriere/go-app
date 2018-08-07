package app

import (
	"net/url"
)

// Page is the interface that describes a webpage.
type Page interface {
	Navigator
	Closer

	// URL returns the URL used to navigate on the page.
	URL() *url.URL

	// Referer returns URL of the page that loaded the current page.
	Referer() *url.URL
}

// PageConfig is a struct that describes a webpage.
type PageConfig struct {
	// The URL of the component to load when the page is created.
	URL string
}
