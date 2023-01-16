package ui

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

// IFlyer is the interface that describes a base with ad spaces surrounded by
// placeholder header and footer.
type IFlyer interface {
	app.UI

	// Sets the ID.
	ID(v string) IFlyer

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) IFlyer

	// Sets the header height in px. Default is 90px.
	HeaderHeight(px int) IFlyer

	// Sets the premium ad height. Default is 250px
	PremiumHeight(px int) IFlyer

	// Sets the footer height in px. Default is 0.
	FooterHeight(px int) IFlyer

	// Sets the banner space located at the top.
	Banner(v ...app.UI) IFlyer

	// Sets the premium space located under the banner.
	Premium(v ...app.UI) IFlyer

	// Sets a bonus space located under the premium space and displayed when
	// there is enough space.
	Bonus(v ...app.UI) IFlyer
}

// Flyer creates a base with ad spaces surrounded by placeholder header and
// footer.
func Flyer() IFlyer {
	return &flyer{
		IheaderHeight:  defaultHeaderHeight,
		IpremiumHeight: 250,
		hpadding:       BaseAdHPadding,
		vpadding:       BaseVPadding,
		layoutID:       "goapp-flyer-layout-" + uuid.NewString(),
	}
}

type flyer struct {
	app.Compo

	Iid            string
	Iclass         string
	IheaderHeight  int
	IpremiumHeight int
	IfooterHeight  int
	Ibanner        []app.UI
	Ipremium       []app.UI
	Ibonus         []app.UI

	hpadding      int
	vpadding      int
	bannerHeight  int
	premiumHeight int
	bonusHeight   int
	layoutID      string
}

func (f *flyer) ID(v string) IFlyer {
	f.Iid = v
	return f
}

func (f *flyer) Class(v string) IFlyer {
	f.Iclass = app.AppendClass(f.Iclass, v)
	return f
}

func (f *flyer) HeaderHeight(px int) IFlyer {
	if px > 0 {
		f.IheaderHeight = px
	}
	return f
}

func (f *flyer) PremiumHeight(px int) IFlyer {
	if px > 0 {
		f.IpremiumHeight = px
	}
	return f
}

func (f *flyer) FooterHeight(px int) IFlyer {
	if px > 0 {
		f.IfooterHeight = px
	}
	return f
}

func (f *flyer) Banner(v ...app.UI) IFlyer {
	f.Ibanner = app.FilterUIElems(v...)
	return f
}

func (f *flyer) Premium(v ...app.UI) IFlyer {
	f.Ipremium = app.FilterUIElems(v...)
	return f
}

func (f *flyer) Bonus(v ...app.UI) IFlyer {
	f.Ibonus = app.FilterUIElems(v...)
	return f
}

func (f *flyer) OnMount(ctx app.Context) {
	f.resize(ctx)
}

func (f *flyer) OnResize(ctx app.Context) {
	f.resize(ctx)
}

func (f *flyer) OnUpdate(ctx app.Context) {
	f.resize(ctx)
}

func (f *flyer) Render() app.UI {
	visible := func(v bool) string {
		if v {
			return "block"
		}
		return "none"
	}

	return app.Div().
		DataSet("goapp-ui", "flyer").
		ID(f.Iid).
		Class(f.Iclass).
		Body(
			app.Div().
				Style("width", "100%").
				Style("height", fmt.Sprintf("calc(100%s - %vpx)", "%", f.vpadding*2)).
				Style("padding", fmt.Sprintf("%vpx 0", f.vpadding)).
				Body(
					app.Div().
						Style("width", fmt.Sprintf("calc(100%s - %vpx)", "%", f.hpadding*2)).
						Style("padding", fmt.Sprintf("0 %vpx", f.hpadding)).
						Style("height", pxToString(f.IheaderHeight)),
					app.Div().
						ID(f.layoutID).
						Style("display", "flex").
						Style("flex-direction", "column").
						Style("justify-content", "center").
						Style("width", fmt.Sprintf("calc(100%s - %vpx)", "%", f.hpadding*2)).
						Style("height", fmt.Sprintf("calc(100%s - %vpx)", "%", f.IheaderHeight+f.IfooterHeight)).
						Style("padding", fmt.Sprintf("0 %vpx", f.hpadding)).
						Style("overflow", "hidden").
						Body(
							app.Div().
								Style("height", pxToString(f.bannerHeight)).
								Style("overflow", "hidden").
								Body(f.Ibanner...),
							app.Div().
								Style("display", visible(f.premiumHeight > 0)).
								Style("height", pxToString(f.premiumHeight)).
								Body(f.Ipremium...),
							app.Div().
								Style("display", visible(f.bonusHeight > 0)).
								Style("height", pxToString(f.bonusHeight)).
								Style("overflow", "hidden").
								Body(f.Ibonus...),
						),
					app.Div().
						Style("width", fmt.Sprintf("calc(100%s - %vpx)", "%", f.hpadding*2)).
						Style("padding", fmt.Sprintf("0 %vpx", f.hpadding)).
						Style("height", pxToString(f.IfooterHeight)),
				),
		)
}

func (f *flyer) resize(ctx app.Context) {
	if app.IsServer {
		return
	}

	layout := app.Window().GetElementByID(f.layoutID)
	if !layout.Truthy() {
		app.Log(errors.New("getting flyer container failed").WithTag("id", f.layoutID))
		return
	}

	remainingHeight := layout.Get("clientHeight").Int()

	var bannerHeight int
	if len(f.Ibanner) != 0 {
		bannerHeight = 600
	}
	if bannerHeight > remainingHeight {
		bannerHeight = remainingHeight
	}
	remainingHeight -= bannerHeight

	premiumHeight := 0
	if remainingHeight-f.IpremiumHeight >= 0 {
		premiumHeight = f.IpremiumHeight
	}
	remainingHeight -= premiumHeight

	bonusHeight := 0
	if premiumHeight > 0 && len(f.Ibonus) != 0 {
		switch {
		case remainingHeight-600 >= 0:
			bonusHeight = 600

		case remainingHeight-250 >= 0:
			bonusHeight = 250

		case remainingHeight-100 >= 0:
			bonusHeight = 100

		case remainingHeight-50 >= 0:
			bonusHeight = 50
		}
	}

	if bannerHeight != f.bannerHeight ||
		premiumHeight != f.premiumHeight ||
		bonusHeight != f.bonusHeight {
		f.bannerHeight = bannerHeight
		f.premiumHeight = premiumHeight
		f.bonusHeight = bonusHeight
	}
}
