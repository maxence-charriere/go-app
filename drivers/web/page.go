// +build js

package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/gopherjs/gopherjs/js"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/html"
)

type Page struct {
	core.Page

	markup     app.Markup
	id         string
	compo      app.Compo
	currentURL string
}

func newPage(c app.PageConfig) app.Page {
	p := &Page{
		markup: app.ConcurrentMarkup(html.NewMarkup(driver.factory)),
		id:     uuid.New().String(),
	}

	driver.elems.Put(p)

	js.Global.Set("golangRequest", p.onPageRequest)
	js.Global.Call("addEventListener", "unload", p.onClose)

	p.Load(p.URL().String())
	return p
}

// ID satisfies the app.Page interface.
func (p *Page) ID() string {
	return p.id
}

func (p *Page) Load(urlFmt string, v ...interface{}) {
	var err error
	defer func() {
		p.SetErr(err)
	}()

	if p.compo != nil {
		p.markup.Dismount(p.compo)
		p.compo = nil
	}

	rawurl := fmt.Sprintf(urlFmt, v...)
	if len(p.currentURL) != 0 && p.currentURL != rawurl {
		driver.NewPage(app.PageConfig{URL: rawurl})
		return
	}
	p.currentURL = rawurl

	var u *url.URL
	if u, err = url.Parse(rawurl); err != nil {
		return
	}

	if len(rawurl) == 0 || u.Path == "/" {
		u, err = url.Parse(driver.URL)
		if err != nil {
			return
		}
	}

	u.Scheme = "compo"
	p.currentURL = u.String()
	n := core.CompoNameFromURL(u)

	var c app.Compo
	if c, err = driver.factory.NewCompo(n); err != nil {
		return
	}

	if _, err = p.markup.Mount(c); err != nil {
		return
	}

	p.compo = c

	if nav, ok := c.(app.Navigable); ok {
		nav.OnNavigate(u)
	}

	var root app.Tag
	if root, err = p.markup.Root(c); err != nil {
		return
	}

	var buffer bytes.Buffer
	enc := html.NewEncoder(&buffer, p.markup, false)
	if err = enc.Encode(root); err != nil {
		return
	}

	js.Global.Get("document").Get("body").Set("innerHTML", buffer.String())
}

func (p *Page) Compo() app.Compo {
	return p.compo
}

func (p *Page) Contains(c app.Compo) bool {
	return p.markup.Contains(c)
}

func (p *Page) Render(c app.Compo) {
	var err error
	defer func() {
		p.SetErr(err)
	}()

	var syncs []app.TagSync
	if syncs, err = p.markup.Update(c); err != nil {
		return
	}

	for _, sync := range syncs {
		if sync.Replace {
			err = p.render(sync)
		} else {
			err = p.renderAttributes(sync)
		}

		if err != nil {
			return
		}
	}
}

func (p *Page) render(sync app.TagSync) error {
	var buffer bytes.Buffer
	enc := html.NewEncoder(&buffer, p.markup, false)
	if err := enc.Encode(sync.Tag); err != nil {
		return err
	}

	payload := &struct {
		*js.Object
		ID    string `js:"id"`
		Compo string `js:"component"`
	}{
		Object: js.Global.Get("Object").New(),
	}
	payload.ID = sync.Tag.ID
	payload.Compo = buffer.String()

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
	payload.ID = sync.Tag.ID
	payload.Attributes = attrs

	js.Global.Call("renderAttributes", payload)
	return nil
}

func (p *Page) Reload() {
	js.Global.Get("location").Call("reload")
}

func (p *Page) CanPrevious() bool {
	return true
}

func (p *Page) Previous() {
	js.Global.Get("history").Call("back")
}

func (p *Page) CanNext() bool {
	return true
}

func (p *Page) Next() {
	js.Global.Get("history").Call("forward")
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

func (p *Page) Close() {
	js.Global.Call("close")
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

	var c app.Compo
	if c, err = p.markup.Compo(mapping.CompoID); err != nil {
		app.Log("page request failed: %s", err)
		return
	}

	p.Render(c)
	if p.Err() != nil {
		app.Log("page request failed: %s", err)
	}
}

func (p *Page) onClose() {
	driver.elems.Delete(p)
}
