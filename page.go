package app

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// Page is the interface that describes a webpage.
type Page interface {
	ElementWithNavigation
	Closer

	// Base returns the base page without any decorators.
	Base() Page

	// URL returns the URL used to navigate on the page.
	URL() *url.URL

	// Referer returns URL of the page that loaded the current page.
	Referer() *url.URL
}

// PageConfig is a struct that describes a webpage.
type PageConfig struct {
	DefaultURL string `json:"default-url"`
}

// NewPageWithLogs returns a decorated version of the given page that logs
// all the operations.
// Uses the default logger.
func NewPageWithLogs(p Page) Page {
	return &pageWithLogs{
		base: p,
	}
}

type pageWithLogs struct {
	base Page
}

func (p *pageWithLogs) ID() uuid.UUID {
	id := p.base.ID()
	Log("page id is", id)
	return id
}

func (p *pageWithLogs) Base() Page {
	return p.base.Base()
}

func (p *pageWithLogs) Load(url string, v ...interface{}) error {
	fmtURL := fmt.Sprintf(url, v...)
	Logf("page %s: loading %s", p.base.ID(), fmtURL)

	err := p.base.Load(url, v...)
	if err != nil {
		Errorf("page %s: loading %s failed: %s", p.base.ID(), fmtURL, err)
	}
	return err
}

func (p *pageWithLogs) Component() Component {
	c := p.base.Component()
	Logf("page %s: mounted component is %T", p.base.ID(), c)
	return c
}

func (p *pageWithLogs) Contains(c Component) bool {
	ok := p.base.Contains(c)
	Logf("page %s: contains %T is %v", p.base.ID(), c, ok)
	return ok
}

func (p *pageWithLogs) Render(c Component) error {
	Logf("page %s: rendering %T", p.base.ID(), c)

	err := p.base.Render(c)
	if err != nil {
		Errorf("page %s: rendering %T failed: %s", p.base.ID(), c, err)
	}
	return err
}

func (p *pageWithLogs) LastFocus() time.Time {
	return p.base.LastFocus()
}

func (p *pageWithLogs) Reload() error {
	Logf("page %s: reloading component %T", p.base.ID(), p.base.Component())

	err := p.base.Reload()
	if err != nil {
		Errorf("page %s: reloading component failed: %s", p.base.ID(), err)
	}
	return err
}

func (p *pageWithLogs) CanPrevious() bool {
	ok := p.base.CanPrevious()
	Logf("page %s: can navigate to previous component is %v", p.base.ID(), ok)
	return ok
}

func (p *pageWithLogs) Previous() error {
	Logf("page %s: navigating to previous component", p.base.ID())

	err := p.base.Previous()
	if err != nil {
		Errorf("page %s: navigating to previous component failed: %s",
			p.base.ID(),
			err,
		)
	}
	return err
}

func (p *pageWithLogs) CanNext() bool {
	ok := p.base.CanNext()
	Logf("page %s: can navigate to next component is %v", p.base.ID(), ok)
	return ok
}

func (p *pageWithLogs) Next() error {
	Logf("page %s: navigating to next component", p.base.ID())

	err := p.base.Next()
	if err != nil {
		Errorf("page %s: navigating to next component failed: %s",
			p.base.ID(),
			err,
		)
	}
	return err
}

func (p *pageWithLogs) URL() *url.URL {
	u := p.base.URL()
	Logf("page %s: URL is %s", p.base.ID(), u)
	return u
}

func (p *pageWithLogs) Referer() *url.URL {
	u := p.base.Referer()
	Logf("page %s: referer is %s", p.base.ID(), u)
	return u
}

func (p *pageWithLogs) Close() {
	Logf("page %s: closing", p.base.ID())
	p.base.Close()
}
