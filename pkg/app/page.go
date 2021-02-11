package app

import (
	"net/url"
	"strings"
)

// Page is the interface that describes a web page.
type Page interface {
	// Returns the page title.
	Title() string

	// Sets the page title.
	SetTitle(string)

	// Returns the page description.
	Description() string

	// Sets the page description.
	SetDescription(string)

	// Returns the page author.
	Author() string

	// Sets the page author.
	SetAuthor(string)

	// Returns the page keywords.
	Keywords() string

	// Sets the page keywords.
	SetKeywords(...string)

	// Set the page loading label.
	SetLoadingLabel(string)

	// Returns the page URL.
	URL() *url.URL

	// Replace the the current page URL by the given one in the browser history.
	//
	// Does not work when pre-rendering.
	ReplaceURL(*url.URL)

	// Returns the page width and height in px.
	Size() (w int, h int)
}

type requestPage struct {
	title        string
	description  string
	author       string
	keywords     string
	loadingLabel string
	url          *url.URL
	width        int
	height       int
}

func (p *requestPage) Title() string {
	return p.title
}

func (p *requestPage) SetTitle(v string) {
	p.title = v
}

func (p *requestPage) Description() string {
	return p.description
}

func (p *requestPage) SetDescription(v string) {
	p.description = v
}

func (p *requestPage) Author() string {
	return p.author
}

func (p *requestPage) SetAuthor(v string) {
	p.author = v
}

func (p *requestPage) Keywords() string {
	return p.keywords
}

func (p *requestPage) SetKeywords(v ...string) {
	p.keywords = strings.Join(v, ", ")
}

func (p *requestPage) SetLoadingLabel(v string) {
	p.loadingLabel = v
}

func (p *requestPage) URL() *url.URL {
	return p.url
}

func (p *requestPage) ReplaceURL(v *url.URL) {
	p.url = v
}

func (p *requestPage) Size() (width int, height int) {
	return p.width, p.height
}

type browserPage struct {
	url *url.URL
}

func (p browserPage) Title() string {
	return Window().
		Get("document").
		Get("title").String()
}

func (p browserPage) SetTitle(v string) {
	Window().
		Get("document").
		Set("title", v)
}

func (p browserPage) Description() string {
	return p.meta("description").getAttr("content")
}

func (p browserPage) SetDescription(v string) {
	p.meta("description").setAttr("content", v)
}

func (p browserPage) Author() string {
	return p.meta("author").getAttr("content")
}

func (p browserPage) SetAuthor(v string) {
	p.meta("author").setAttr("content", v)
}

func (p browserPage) Keywords() string {
	return p.meta("keywords").getAttr("content")
}

func (p browserPage) SetKeywords(v ...string) {
	p.meta("keywords").setAttr("content", strings.Join(v, ", "))
}

func (p browserPage) SetLoadingLabel(v string) {
}

func (p browserPage) URL() *url.URL {
	if p.url != nil {
		return p.url
	}
	return Window().URL()
}

func (p browserPage) ReplaceURL(v *url.URL) {
	Window().replaceHistory(v)
}

func (p browserPage) Size() (width int, height int) {
	return Window().Size()
}

func (p browserPage) meta(name string) Value {
	return Window().
		Get("document").
		Call("querySelector", "meta[name='"+name+"']")
}
