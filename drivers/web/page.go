// +build js

package web

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/gopherjs/gopherjs/js"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/pkg/errors"
)

// Page implements the app.Page interface.
type Page struct {
	core.Page

	id         string
	dom        dom.Engine
	compo      app.Compo
	currentURL string
}

func newPage(c app.PageConfig) app.Page {
	p := &Page{
		id: uuid.New().String(),
		dom: dom.Engine{
			Factory:        driver.factory,
			Resources:      driver.Resources,
			AttrTransforms: []dom.Transform{dom.JsToGoHandler},
		},
	}

	p.dom.Sync = p.render

	driver.elems.Put(p)

	js.Global.Set("golangRequest", p.onPageRequest)
	js.Global.Call("addEventListener", "unload", p.onClose)

	u := p.URL()
	u.Path = js.Global.Get("loadedComp").String()
	p.Load(u.String())
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

	rawurl := fmt.Sprintf(urlFmt, v...)

	if len(p.currentURL) != 0 && p.currentURL != rawurl {
		driver.NewPage(app.PageConfig{URL: rawurl})
		return
	}

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

	p.compo = c

	if err = p.dom.New(c); err != nil {
		return
	}

	if nav, ok := c.(app.Navigable); ok {
		nav.OnNavigate(u)
	}
}

func (p *Page) Compo() app.Compo {
	return p.compo
}

func (p *Page) Contains(c app.Compo) bool {
	return p.dom.Contains(c)
}

func (p *Page) Render(c app.Compo) {
	p.SetErr(p.dom.Render(c))
}

func (p *Page) render(changes interface{}) error {
	b, err := json.Marshal(changes)
	if err != nil {
		return errors.Wrap(err, "marshal changes failed")
	}

	c := js.Global.Get("JSON").Call("parse", string(b))
	js.Global.Call("render", c)
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

func (p *Page) onPageRequest(mappingStr string) {
	var m dom.Mapping
	if err := json.Unmarshal([]byte(mappingStr), &m); err != nil {
		app.Logf("page callback failed: %s", err)
		return
	}

	c, err := p.dom.CompoByID(m.CompoID)
	if err != nil {
		app.Logf("page callback failed: %s", err)
		return
	}

	var f func()
	if f, err = m.Map(c); err != nil {
		app.Logf("page callback failed: %s", err)
		return
	}

	if f != nil {
		f()
		return
	}

	app.Render(c)
}

func (p *Page) onClose() {
	driver.elems.Delete(p)
}
