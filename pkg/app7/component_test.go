package app

import "net/url"

type hello struct {
	Compo
}

func (f *hello) OnMount(Context) {
}

func (f *hello) OnNav(Context, *url.URL) {
}

func (f *hello) OnDismount(Context) {
}

func (f *hello) Render() UI {
	return Div().Body(
		H1().Text("hello world"),
	)
}
