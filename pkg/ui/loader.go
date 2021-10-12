package ui

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

const (
	defaultLoaderErrorIcon = `<svg style="width:%vpx;height:%vpx" viewBox="0 0 24 24">
    	<path fill="currentColor" d="M22 14H21C21 10.13 17.87 7 14 7H13V5.73C13.6 5.39 14 4.74 14 4C14 2.9 13.11 2 12 2S10 2.9 10 4C10 4.74 10.4 5.39 11 5.73V7H10C6.13 7 3 10.13 3 14H2C1.45 14 1 14.45 1 15V18C1 18.55 1.45 19 2 19H3V20C3 21.11 3.9 22 5 22H19C20.11 22 21 21.11 21 20V19H22C22.55 19 23 18.55 23 18V15C23 14.45 22.55 14 22 14M9.86 16.68L8.68 17.86L7.5 16.68L6.32 17.86L5.14 16.68L6.32 15.5L5.14 14.32L6.32 13.14L7.5 14.32L8.68 13.14L9.86 14.32L8.68 15.5L9.86 16.68M18.86 16.68L17.68 17.86L16.5 16.68L15.32 17.86L14.14 16.68L15.32 15.5L14.14 14.32L15.32 13.14L16.5 14.32L17.68 13.14L18.86 14.32L17.68 15.5L18.86 16.68Z" />
	</svg>`
)

type ILoader interface {
	app.UI

	// Sets the ID.
	ID(v string) ILoader

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) ILoader

	// Sets the style. Multiple styles can be defined by successive calls.
	Style(k, v string) ILoader

	// Reports whether the loader is active.
	Loading(v bool) ILoader

	// Sets the size of the rotating circle in px. Default is 60px.
	Size(px int) ILoader

	// Sets the color of the rotating head. Default is white.
	Color(v string) ILoader

	// Sets the time it take to fully rotate. Default is 500ms.
	Speed(v time.Duration) ILoader

	// Sets the space between the loader and the label in px. Default is 18px.
	Spacing(px int) ILoader

	// Sets the label. Default is "Loading...".
	Label(v string) ILoader

	// Sets the error that occured during loading.
	Err(err error) ILoader

	// Sets the error icon.
	ErrIcon(v string) ILoader
}

func Loader() ILoader {
	return &loader{
		Isize:    60,
		Icolor:   "white",
		Ispeed:   time.Millisecond * 500,
		Ispacing: 18,
		Ilabel:   "Loading...",
		IerrIcon: defaultLoaderErrorIcon,
	}

}

type loader struct {
	app.Compo

	Iid      string
	Iclass   string
	Istyles  []style
	Iloading bool
	Isize    int
	Icolor   string
	Ispacing int
	Ilabel   string
	Ispeed   time.Duration
	Ierr     error
	IerrIcon string
}

func (l *loader) ID(v string) ILoader {
	l.Iid = v
	return l
}

func (l *loader) Class(v string) ILoader {
	l.Iclass = app.AppendClass(l.Iclass, v)
	return l
}

func (l *loader) Style(k, v string) ILoader {
	if v == "" {
		return l
	}
	l.Istyles = append(l.Istyles, style{
		key:   k,
		value: v,
	})
	return l
}

func (l *loader) Loading(v bool) ILoader {
	l.Iloading = v
	return l
}

func (l *loader) Size(px int) ILoader {
	l.Isize = px
	return l
}

func (l *loader) Color(v string) ILoader {
	l.Icolor = v
	return l
}

func (l *loader) Speed(v time.Duration) ILoader {
	l.Ispeed = v
	return l
}

func (l *loader) Spacing(px int) ILoader {
	l.Ispacing = px
	return l
}

func (l *loader) Label(v string) ILoader {
	l.Ilabel = v
	return l
}

func (l *loader) Err(err error) ILoader {
	l.Ierr = err
	return l
}

func (l *loader) ErrIcon(v string) ILoader {
	l.IerrIcon = v
	return l
}

func (l *loader) Render() app.UI {
	body := app.Aside().
		ID(l.Iid).
		Class(l.Iclass).
		Body(
			Stack().
				Style("width", "100%").
				Style("height", "100%").
				Center().
				Middle().
				Content(
					app.If(l.Ierr == nil,
						app.Div().
							Style("width", pxToString(l.Isize-4)).
							Style("height", pxToString(l.Isize-4)).
							Style("min-width", pxToString(l.Isize-4)).
							Style("min-height", pxToString(l.Isize-4)).
							Style("border", "2px solid currentColor").
							Style("border-top", "2px solid "+l.Icolor).
							Style("border-radius", "50%").
							Style("animation", fmt.Sprintf("goapp-spin-frames %vms linear infinite", l.Ispeed.Milliseconds())),
					).Else(
						Icon().
							Size(l.Isize).
							Src(l.IerrIcon),
					),
					app.Div().
						Style("margin-left", pxToString(l.Ispacing)).
						Body(
							app.If(l.Ierr == nil,
								app.Div().Text(l.Ilabel),
							).Else(
								app.Div().
									Style("white-space", "pre-wrap").
									Text(l.Ierr),
							),
						),
				),
		)

	for _, s := range l.Istyles {
		body.Style(s.key, s.value)
	}

	if l.Ierr == nil && !l.Iloading {
		body.Style("display", "none")
	}

	return body
}
