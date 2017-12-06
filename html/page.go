package html

// PageConfig is the struct that describes a page.
type PageConfig struct {
	// The page title.
	Title string

	// The app.js script.
	AppJS string

	// The default component rendering.
	DefaultComponent string

	// The CSS filenames to include.
	CSS []string

	// The javascript filenames to include.
	Javasripts []string
}

func Page(config PageConfig) (page string, err error) {
	tmpl := ``

	return
}
