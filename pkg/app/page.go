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

	// Returns the page language.
	Lang() string

	// Set the page language.
	SetLang(string)

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

	// Returns the image used by social networks when linking the page.
	Image() string

	// Set the image used by social networks when linking the page.
	SetImage(string)

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
	lang         string
	description  string
	author       string
	keywords     string
	loadingLabel string
	image        string
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

func (p *requestPage) Lang() string {
	return p.lang
}

func (p *requestPage) SetLang(v string) {
	p.lang = v
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

func (p *requestPage) Image() string {
	return p.image
}

func (p *requestPage) SetImage(v string) {
	p.image = v
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
	url        *url.URL
	dispatcher Dispatcher
}

func (p browserPage) Title() string {
	return Window().
		Get("document").
		Get("title").
		String()
}

func (p browserPage) SetTitle(v string) {
	Window().Get("document").Set("title", v)
	p.metaByProperty("og:title").setAttr("content", v)
}

func (p browserPage) Lang() string {
	return Window().
		Get("document").
		Get("documentElement").
		Get("lang").
		String()
}

func (p browserPage) SetLang(v string) {
	Window().
		Get("document").
		Get("documentElement").
		Set("lang", v)
}

func (p browserPage) Description() string {
	return p.metaByName("description").getAttr("content")
}

func (p browserPage) SetDescription(v string) {
	p.metaByName("description").setAttr("content", v)
	p.metaByProperty("og:description").setAttr("content", v)
}

func (p browserPage) Author() string {
	return p.metaByName("author").getAttr("content")
}

func (p browserPage) SetAuthor(v string) {
	p.metaByName("author").setAttr("content", v)
}

func (p browserPage) Keywords() string {
	return p.metaByName("keywords").getAttr("content")
}

func (p browserPage) SetKeywords(v ...string) {
	p.metaByName("keywords").setAttr("content", strings.Join(v, ", "))
}

func (p browserPage) SetLoadingLabel(v string) {
}

func (p browserPage) Image() string {
	return p.metaByProperty("og:image").getAttr("content")
}

func (p browserPage) SetImage(v string) {
	p.metaByProperty("og:image").setAttr("content", p.dispatcher.resolveStaticResource(v))
}

func (p browserPage) URL() *url.URL {
	if p.url != nil {
		return p.url
	}
	return Window().URL()
}

func (p browserPage) ReplaceURL(v *url.URL) {
	Window().replaceHistory(v)
	p.metaByProperty("og:url").setAttr("content", v.String())
}

func (p browserPage) Size() (width int, height int) {
	return Window().Size()
}

func (p browserPage) metaByName(v string) Value {
	meta := Window().
		Get("document").
		Call("querySelector", "meta[name='"+v+"']")

	if meta.IsNull() {
		meta, _ = Window().createElement("meta", "")
		meta.setAttr("name", v)

		Window().Get("document").
			Call("getElementsByTagName", "head").
			Index(0).
			appendChild(meta)
	}

	return meta
}

func (p browserPage) metaByProperty(v string) Value {
	meta := Window().
		Get("document").
		Call("querySelector", "meta[property='"+v+"']")

	if meta.IsNull() {
		meta, _ = Window().createElement("meta", "")
		meta.setAttr("property", v)

		Window().Get("document").
			Call("getElementsByTagName", "head").
			Index(0).
			appendChild(meta)
	}

	return meta
}
