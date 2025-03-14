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
	SetTitle(format string, v ...any)

	// Returns the page language.
	Lang() string

	// Set the page language.
	SetLang(v string)

	// Returns the page description.
	Description() string

	// Sets the page description.
	SetDescription(format string, v ...any)

	// Returns the page author.
	Author() string

	// Sets the page author.
	SetAuthor(format string, v ...any)

	// Returns the page keywords.
	Keywords() string

	// Sets the page keywords.
	SetKeywords(v ...string)

	// Returns the page resources to preload.
	Preloads() []Preload

	// Sets resources to preload.
	SetPreloads(v ...Preload)

	// Set the page loading label.
	SetLoadingLabel(format string, v ...any)

	// Returns the image used by social networks when linking the page.
	Image() string

	// Set the image used by social networks when linking the page.
	SetImage(v string)

	// Returns the page URL.
	URL() *url.URL

	// Replace the the current page URL by the given one in the browser history.
	//
	// Does not work when pre-rendering.
	ReplaceURL(v *url.URL)

	// Returns the page width and height in px.
	Size() (w int, h int)

	// Set the Twitter card.
	SetTwitterCard(v TwitterCard)

	// Set the page's canonical link.
	SetCanonicalLink(format string, v ...any)
}

type requestPage struct {
	url        *url.URL
	resolveURL func(string) string

	title          string
	lang           string
	description    string
	author         string
	keywords       string
	preloads       []Preload
	loadingLabel   string
	image          string
	canonicalLink  string
	width          int
	height         int
	twitterCardMap map[string]string
}

func makeRequestPage(origin *url.URL, resolveURL func(string) string) requestPage {
	return requestPage{
		url:        origin,
		resolveURL: resolveURL,
	}
}

func (p *requestPage) Title() string {
	return p.title
}

func (p *requestPage) SetTitle(format string, v ...any) {
	p.title = FormatString(format, v...)
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

func (p *requestPage) SetDescription(format string, v ...any) {
	p.description = FormatString(format, v...)
}

func (p *requestPage) Author() string {
	return p.author
}

func (p *requestPage) SetAuthor(format string, v ...any) {
	p.author = FormatString(format, v...)
}

func (p *requestPage) Keywords() string {
	return p.keywords
}

func (p *requestPage) SetKeywords(v ...string) {
	p.keywords = strings.Join(v, ", ")
}

func (p *requestPage) Preloads() []Preload {
	return p.preloads
}

func (p *requestPage) SetPreloads(v ...Preload) {
	for i, r := range v {
		v[i].Href = p.resolveURL(r.Href)
	}
	p.preloads = v
}

func (p *requestPage) SetLoadingLabel(format string, v ...any) {
	p.loadingLabel = FormatString(format, v...)
}

func (p *requestPage) Image() string {
	return p.image
}

func (p *requestPage) SetImage(v string) {
	if v != "" {
		p.image = p.resolveURL(v)
	}
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

func (p *requestPage) SetTwitterCard(v TwitterCard) {
	v.Image = p.resolveURL(v.Image)
	p.twitterCardMap = v.toMap()
}

func (p *requestPage) SetCanonicalLink(format string, v ...any) {
	if canonicalLink := FormatString(format, v...); canonicalLink != "" {
		p.canonicalLink = p.resolveURL(canonicalLink)
	}
}

type browserPage struct {
	resolveURL func(string) string
}

func makeBrowserPage(resolveURL func(string) string) browserPage {
	return browserPage{resolveURL: resolveURL}
}

func (p browserPage) Title() string {
	return Window().
		Get("document").
		Get("title").
		String()
}

func (p browserPage) SetTitle(format string, v ...any) {
	title := FormatString(format, v...)
	Window().Get("document").Set("title", title)
	p.metaByProperty("og:title").setAttr("content", title)
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

func (p browserPage) SetDescription(format string, v ...any) {
	description := FormatString(format, v...)
	p.metaByName("description").setAttr("content", description)
	p.metaByProperty("og:description").setAttr("content", description)
}

func (p browserPage) Author() string {
	return p.metaByName("author").getAttr("content")
}

func (p browserPage) SetAuthor(format string, v ...any) {
	p.metaByName("author").setAttr("content", FormatString(format, v...))
}

func (p browserPage) Keywords() string {
	return p.metaByName("keywords").getAttr("content")
}

func (p browserPage) SetKeywords(v ...string) {
	p.metaByName("keywords").setAttr("content", strings.Join(v, ", "))
}

func (p browserPage) SetLoadingLabel(format string, v ...any) {
}

func (p browserPage) Preloads() []Preload {
	return nil
}

func (p browserPage) SetPreloads(v ...Preload) {
}

func (p browserPage) Image() string {
	return p.metaByProperty("og:image").getAttr("content")
}

func (p browserPage) SetImage(v string) {
	if v != "" {
		p.metaByProperty("og:image").setAttr("content", p.resolveURL(v))
	}
}

func (p browserPage) URL() *url.URL {
	return Window().URL()
}

func (p browserPage) ReplaceURL(v *url.URL) {
	Window().replaceHistory(v)
	p.metaByProperty("og:url").setAttr("content", v.String())
}

func (p browserPage) Size() (width int, height int) {
	return Window().Size()
}

func (p browserPage) SetTwitterCard(v TwitterCard) {
	v.Image = p.resolveURL(v.Image)
	head := Window().Get("document").Get("head")

	for k, v := range v.toMap() {
		if v == "" {
			continue
		}
		meta, _ := Window().createElement("meta", "")
		meta.setAttr("name", k)
		meta.setAttr("content", v)
		head.appendChild(meta)
	}
}

func (p browserPage) SetCanonicalLink(format string, v ...any) {
}

func (p browserPage) metaByName(v string) Value {
	meta := Window().
		Get("document").
		Call("querySelector", "meta[name='"+v+"']")

	if meta.IsNull() {
		meta, _ = Window().createElement("meta", "")
		meta.setAttr("name", v)

		Window().Get("document").
			Get("head").
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
			Get("head").
			appendChild(meta)
	}

	return meta
}

type Preload struct {
	Type          string
	As            string
	Href          string
	FetchPriority string
}
