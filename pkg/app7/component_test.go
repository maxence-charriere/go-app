package app

import "net/url"

type hello struct {
	Compo

	Greeting string
}

func (h *hello) OnMount(Context) {
}

func (h *hello) OnNav(Context, *url.URL) {
}

func (h *hello) OnDismount(Context) {
}

func (h *hello) Render() UI {
	return Div().Body(
		H1().Body(
			Text("hello, "),
			Text(h.Greeting),
		),
	)
}

type foo struct {
	Compo
	Bar string
}

func (f *foo) Render() UI {
	return &bar{Value: f.Bar}
}

type bar struct {
	Compo
	Value string
}

func (b *bar) Render() UI {
	return Text(b.Value)
}
