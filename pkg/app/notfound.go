package app

var (
	// NotFound is the ui element that is displayed when a request is not
	// routed.
	NotFound UI = &notFound{}
)

type notFound struct {
	Compo
	Icon string
}

func (n *notFound) OnMount(Context) {
	links := Window().Get("document").Call("getElementsByTagName", "link")

	for i := 0; i < links.Length(); i++ {
		link := links.Index(i)
		rel := link.Call("getAttribute", "rel")

		if rel.String() == "icon" {
			favicon := link.Call("getAttribute", "href")
			n.Icon = favicon.String()
			return
		}
	}
}

func (n *notFound) Render() UI {
	return Div().
		Class("goapp-app-info").
		Body(
			Div().
				Class("goapp-notfound-title").
				Body(
					Text("4"),
					Img().
						Class("goapp-logo").
						Alt("0").
						Src(n.Icon),
					Text("4"),
				),
			P().
				Class("goapp-label").
				Text("Not Found"),
		)
}
