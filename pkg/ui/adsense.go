package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/logs"
)

// IAdsenseDisplay is the interface that describes a responsive Adsense display
// unit.
//
// Note that the Adsense script must be loaded in the app.Handler.RawHeaders.
type IAdsenseDisplay interface {
	app.UI

	// Sets the ID.
	ID(v string) IAdsenseDisplay

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) IAdsenseDisplay

	// Sets the AdSense slot.
	Client(v string) IAdsenseDisplay

	// Sets the AdSense slot.
	Slot(v string) IAdsenseDisplay
}

// AdsenseDisplay creates a responsive Adsense display unit.
func AdsenseDisplay() IAdsenseDisplay {
	return &adsenseDisplay{
		id: "goapp-adsense-display-" + uuid.NewString(),
	}
}

type adsenseDisplay struct {
	app.Compo

	Iid     string
	Iclass  string
	Iclient string
	Islot   string

	id          string
	currentPath string
	width       int
	height      int
}

func (d *adsenseDisplay) ID(v string) IAdsenseDisplay {
	d.Iid = v
	return d
}

func (d *adsenseDisplay) Class(v string) IAdsenseDisplay {
	d.Iclass = app.AppendClass(d.Iclass, v)
	return d
}

func (d *adsenseDisplay) Client(v string) IAdsenseDisplay {
	d.Iclient = v
	return d
}

func (d *adsenseDisplay) Slot(v string) IAdsenseDisplay {
	d.Islot = v
	return d
}

func (d *adsenseDisplay) OnMount(ctx app.Context) {
	ctx.Defer(d.load)
}

func (d *adsenseDisplay) OnNav(ctx app.Context) {
	ctx.Defer(d.load)
}

func (d *adsenseDisplay) OnResize(ctx app.Context) {
	ctx.Defer(d.load)
}

func (d *adsenseDisplay) OnUpdate(ctx app.Context) {
	ctx.Defer(d.load)
}

func (d *adsenseDisplay) Render() app.UI {
	return app.Div().
		DataSet("goapp-ui", "adsenseDisplay").
		ID(d.Iid).
		Class(d.Iclass).
		Body(
			app.Ins().ID(d.id),
		)
}

func (d *adsenseDisplay) containerID() string {
	return d.id
}

func (d *adsenseDisplay) client() string {
	return d.Iclient
}

func (d *adsenseDisplay) slot() string {
	return d.Islot
}

func (d *adsenseDisplay) load(ctx app.Context) {
	ins := app.Window().GetElementByID(d.id)
	if !ins.Truthy() {
		return
	}

	layout := ins.Get("parentElement")
	w := layout.Get("clientWidth").Int()
	h := layout.Get("clientHeight").Int()

	currentPath := ctx.Page().URL().Path

	if w != d.width || h != d.height || currentPath != d.currentPath {
		d.width = w
		d.height = h
		d.currentPath = currentPath
		ads.Push(ctx, d)
	}
}

type adUnit interface {
	app.UI

	containerID() string
	client() string
	slot() string
}

type adPimp struct {
	mutex    sync.Mutex
	units    map[adUnit]struct{}
	interval time.Duration
}

func newAdPimp() *adPimp {
	return &adPimp{
		units:    make(map[adUnit]struct{}),
		interval: time.Millisecond * 100,
	}
}

func (p *adPimp) Push(ctx app.Context, units ...adUnit) {
	ctx.Async(func() {
		time.Sleep(p.interval)

		p.mutex.Lock()
		defer p.mutex.Unlock()

		p.push(ctx, units...)
	})
}

func (p *adPimp) push(ctx app.Context, units ...adUnit) {
	if app.IsServer {
		return
	}

	for _, u := range units {
		p.addUnit(u)
	}

	for u := range p.units {
		p.pushUnit(u)
	}

	if len(p.units) == 0 {
		return
	}

	p.Push(ctx)
}

func (p *adPimp) addUnit(u adUnit) {
	containerID := u.containerID()
	if containerID == "" {
		return
	}

	ins := app.Window().GetElementByID(containerID)
	if !ins.Truthy() {
		p.removeUnit(u)
		return
	}

	for lastChild := ins.Get("lastChild"); lastChild.Truthy(); lastChild = ins.Get("lastChild") {
		ins.Call("removeChild", lastChild)
	}
	ins.Set("className", "")
	ins.Get("dataset").Set("adsbygoogleStatus", "")
	ins.Get("dataset").Set("adStatus", "")
	ins.Set("style", "")

	p.units[u] = struct{}{}
}

func (p *adPimp) pushUnit(u adUnit) {
	if !u.Mounted() {
		p.removeUnit(u)
		return
	}

	ins := app.Window().GetElementByID(u.containerID())
	if !ins.Truthy() {
		p.removeUnit(u)
		return
	}

	if status := ins.Get("dataset").Get("adsbygoogleStatus").String(); status != "" {
		p.removeUnit(u)
		return
	}

	layout := ins.Get("parentElement")
	w := layout.Get("clientWidth").Int()
	h := layout.Get("clientHeight").Int()
	if w == 0 && h == 0 {
		app.Log(logs.New("ad unit not visible").
			WithTag("width", w).
			WithTag("height", h).
			WithTag("slot", u.slot()).
			WithTag("container-id", u.containerID()),
		)
		p.removeUnit(u)
		return
	}
	if w*h == 0 {
		return
	}

	adsbygoogle := app.Window().Get("adsbygoogle")
	if !adsbygoogle.Truthy() {
		app.Log(logs.New("adsbygoogle is not loaded"))
		return
	}

	app.Log(logs.New("loading ad unit").
		WithTag("width", w).
		WithTag("height", h).
		WithTag("slot", u.slot()).
		WithTag("container-id", u.containerID()),
	)

	ins.Set("className", "adsbygoogle")
	ins.Get("dataset").Set("adClient", u.client())
	ins.Get("dataset").Set("adSlot", u.slot())
	ins.Set("style", fmt.Sprintf("display:inline-block;width:%vpx;height:%vpx;overflow:hidden", w, h))

	adsbygoogle.Call("push", map[string]interface{}{})
}

func (p *adPimp) removeUnit(u adUnit) {
	delete(p.units, u)
}

var (
	ads = newAdPimp()
)
