package ui

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/errors"
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

	// Sets the duration to wait before ads are loaded. Default is 150ms.
	Refresh(d time.Duration) IAdsenseDisplay
}

// AdsenseDisplay creates a responsive Adsense display unit.
func AdsenseDisplay() IAdsenseDisplay {
	return &adsenseDisplay{
		Irefresh: time.Millisecond * 150,
		id:       "goapp-adsense-display-" + uuid.NewString(),
	}
}

type adsenseDisplay struct {
	app.Compo

	Iid      string
	Iclass   string
	Iclient  string
	Islot    string
	Irefresh time.Duration

	id     string
	width  int
	height int
}

func (d *adsenseDisplay) ID(v string) IAdsenseDisplay {
	d.Iid = v
	return d
}

func (d *adsenseDisplay) Class(v string) IAdsenseDisplay {
	if v == "" {
		return d
	}
	if d.Iclass != "" {
		d.Iclass += " "
	}
	d.Iclass += v
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

func (d *adsenseDisplay) Refresh(du time.Duration) IAdsenseDisplay {
	d.Irefresh = du
	return d
}

func (d *adsenseDisplay) OnMount(ctx app.Context) {
	ctx.Defer(d.resize)
}

func (d *adsenseDisplay) OnNav(ctx app.Context) {
	ctx.Defer(d.resize)
}

func (d *adsenseDisplay) OnResize(ctx app.Context) {
	ctx.Dispatch(d.resize)
}

func (d *adsenseDisplay) OnUpdate(ctx app.Context) {
	ctx.Dispatch(d.resize)
}

func (d *adsenseDisplay) Render() app.UI {
	return app.Div().
		DataSet("goapp-ui", "adsenseDisplay").
		ID(d.Iid).
		Class(d.Iclass).
		Body(
			app.Ins().
				ID(d.id).
				Class("adsbygoogle").
				Style("display", "block").
				Style("width", "100%").
				Style("height", "100%").
				Style("overflow", "hidden").
				DataSet("ad-client", d.Iclient).
				DataSet("ad-slot", d.Islot),
		)
}

func (d *adsenseDisplay) resize(ctx app.Context) {
	if app.IsServer {
		return
	}

	ins := app.Window().GetElementByID(d.id)
	if !ins.Truthy() {
		app.Log(errors.New("getting adsense display ins failed").Tag("id", d.id))
		return
	}
	w := ins.Get("clientWidth").Int()
	h := ins.Get("clientHeight").Int()

	if w != d.width || h != d.height {
		ins.Set("innerHTML", "")
		ins.Get("dataset").Set("adsbygoogleStatus", nil)
		ins.Get("dataset").Set("adStatus", nil)
		d.width = w
		d.height = h
		fmt.Println("ad:", w, h)
		refreshAdsenseUnits(ctx, d.Irefresh)
	}
}

var (
	adRefresh *time.Timer
)

func refreshAdsenseUnits(ctx app.Context, refresh time.Duration) {
	if adRefresh != nil {
		adRefresh.Reset(refresh)
		return
	}

	adRefresh = time.AfterFunc(refresh, func() {
		adsbygoogle := app.Window().Get("adsbygoogle")
		if !adsbygoogle.Truthy() {
			app.Logf("%s", errors.New("getting adsbygoogle failed"))
			return
		}
		adsbygoogle.Call("push", map[string]interface{}{})
		fmt.Println("adsense PUSH")
	})
}
