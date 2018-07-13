// +build js

package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gopherjs/gopherjs/js"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
	"github.com/murlokswarm/app/internal/core"
)

type Page struct {
	core.Elem

	id         uuid.UUID
	markup     app.Markup
	component  app.Component
	lastFocus  time.Time
	currentURL string
}

func newPage(c app.PageConfig) (app.Page, error) {
	var markup app.Markup = html.NewMarkup(driver.factory)
	markup = app.ConcurrentMarkup(markup)

	page := &Page{
		id:        uuid.New(),
		markup:    markup,
		lastFocus: time.Now(),
	}

	driver.elems.Put(page)

	js.Global.Set("golangRequest", page.onPageRequest)

	js.Global.Call("addEventListener", "unload", func() {
		driver.elems.Delete(page)
	})

	err := page.Load(page.URL().String())
	return page, err
}

// ID satisfies the app.Page interface.
func (p *Page) ID() uuid.UUID {
	return p.id
}

func (p *Page) Load(rawurl string, v ...interface{}) error {
	if p.component != nil {
		p.markup.Dismount(p.component)
	}

	rawurl = fmt.Sprintf(rawurl, v...)
	if len(p.currentURL) != 0 && p.currentURL != rawurl {
		return driver.NewPage(app.PageConfig{
			DefaultURL: rawurl,
		})
	}
	p.currentURL = rawurl

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	if len(u.Path) == 0 || u.Path == "/" {
		u.Path = driver.DefaultURL
	}
	u.Scheme = "compo"

	var compo app.Component
	if compo, err = driver.factory.New(app.ComponentNameFromURL(u)); err != nil {
		return err
	}

	if _, err = p.markup.Mount(compo); err != nil {
		return err
	}
	p.component = compo

	if navigable, ok := compo.(app.Navigable); ok {
		navigable.OnNavigate(u)
	}

	var root app.Tag
	if root, err = p.markup.Root(compo); err != nil {
		return err
	}

	var buffer bytes.Buffer
	enc := html.NewEncoder(&buffer, p.markup, false)
	if err = enc.Encode(root); err != nil {
		return err
	}

	js.Global.Get("document").Get("body").Set("innerHTML", buffer.String())
	return nil
}

func (p *Page) Component() app.Component {
	return p.component
}

func (p *Page) Contains(c app.Component) bool {
	return p.markup.Contains(c)
}

func (p *Page) Render(c app.Component) error {
	syncs, err := p.markup.Update(c)
	if err != nil {
		return err
	}

	for _, sync := range syncs {
		if sync.Replace {
			err = p.render(sync)
		} else {
			err = p.renderAttributes(sync)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Page) render(sync app.TagSync) error {
	var buffer bytes.Buffer
	enc := html.NewEncoder(&buffer, p.markup, false)
	if err := enc.Encode(sync.Tag); err != nil {
		return err
	}

	payload := &struct {
		*js.Object
		ID        string `js:"id"`
		Component string `js:"component"`
	}{
		Object: js.Global.Get("Object").New(),
	}
	payload.ID = sync.Tag.ID.String()
	payload.Component = buffer.String()

	js.Global.Call("render", payload)
	return nil
}

func (p *Page) renderAttributes(sync app.TagSync) error {
	attrs := make(app.AttributeMap, len(sync.Tag.Attributes))
	for name, val := range sync.Tag.Attributes {
		attrs[name] = html.AttrValueFormatter{
			Name:    name,
			Value:   val,
			CompoID: sync.Tag.CompoID,
			Factory: driver.factory,
		}.Format()
	}

	payload := &struct {
		*js.Object
		ID         string           `js:"id"`
		Attributes app.AttributeMap `js:"attributes"`
	}{
		Object: js.Global.Get("Object").New(),
	}
	payload.ID = sync.Tag.ID.String()
	payload.Attributes = attrs

	js.Global.Call("renderAttributes", payload)
	return nil
}

func (p *Page) LastFocus() time.Time {
	return p.lastFocus
}

func (p *Page) Reload() error {
	js.Global.Get("location").Call("reload")
	return nil
}

func (p *Page) CanPrevious() bool {
	return true
}

func (p *Page) Previous() error {
	js.Global.Get("history").Call("back")
	return nil
}

func (p *Page) CanNext() bool {
	return true
}

func (p *Page) Next() error {
	js.Global.Get("history").Call("forward")
	return nil
}

func (p *Page) URL() *url.URL {
	u, _ := url.Parse(js.
		Global.
		Get("location").
		Get("href").
		String(),
	)
	return u
}

func (p *Page) Referer() *url.URL {
	u, _ := url.Parse(js.
		Global.
		Get("document").
		Get("referrer").
		String(),
	)
	return u
}

func (p *Page) Close() error {
	js.Global.Call("close")
	return nil
}

func (p *Page) WhenPage(f func(app.Page)) {
	f(p)
}

func (p *Page) WhenNavigator(f func(app.Navigator)) {
	f(p)
}

func (p *Page) onPageRequest(j string) {
	var mapping app.Mapping
	if err := json.Unmarshal([]byte(j), &mapping); err != nil {
		app.Log("page request failed: %s", err)
		return
	}

	fn, err := p.markup.Map(mapping)
	if err != nil {
		app.Log("page request failed: %s", err)
		return
	}

	if fn != nil {
		fn()
		return
	}

	var compo app.Component
	if compo, err = p.markup.Component(mapping.CompoID); err != nil {
		app.Log("page request failed: %s", err)
		return
	}

	if err = p.Render(compo); err != nil {
		app.Log("page request failed: %s", err)
	}
}
