package app

type notFound struct {
	Compo
	Icon string
}

func (n *notFound) OnMount() {
	links := Window().Get("document").Call("getElementsByTagName", "link")

	for i := 0; i < links.Length(); i++ {
		link := links.Index(i)
		rel := link.Call("getAttribute", "rel")

		if rel.String() == "icon" {
			favicon := link.Call("getAttribute", "href")
			n.Icon = favicon.String()
			n.Update()
			return
		}
	}
}

func (n *notFound) Render() UI {
	return Div().
		Class("app-wasm-layout").
		Body(
			Div().
				Class("app-notfound-title").
				Body(
					Text("4"),
					Img().
						Class("app-wasm-icon").
						Alt("0").
						Src(n.Icon),
					Text("4"),
				),
			P().
				Class("app-wasm-label").
				Text("Not Found"),
		)
}
