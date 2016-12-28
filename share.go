package app

import (
	"net/url"
)

// Sharer describes a sharing service.
type Sharer interface {
	Text(v string)

	URL(v *url.URL)
}

// Share returns the sharing service.
func Share() Sharer {
	return driver.Share()
}
