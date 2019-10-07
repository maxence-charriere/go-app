// +build wasm

package main

import (
	"syscall/js"

	"github.com/maxence-charriere/app/pkg/app"
)

// City is a component that shows how navigation works.
type City struct {
	City    string
	Current city
	Store   map[string]city
}

// OnMount initializes the component when it is mounted.
func (n *City) OnMount() {
	n.Store = map[string]city{
		"paris": {
			Name:  "Paris",
			Image: "paris.jpg",
			Description: `
Paris, France's capital, is a major European city and a global center for art,
fashion, gastronomy and culture. Its 19th-century cityscape is crisscrossed by
wide boulevards and the River Seine. Beyond such landmarks as the Eiffel Tower
and the 12th-century, Gothic Notre-Dame cathedral, the city is known for its
cafe culture and designer boutiques along the Rue du Faubourg Saint-Honoré.
			`,
		},
		"sf": {
			Name:  "SF",
			Image: "sf.jpg",
			Description: `
San Francisco, officially City and County of San Francisco and colloquially
known as SF, San Fran or "The City", is a city in—and the cultural, commercial,
and financial center of—Northern California.
			`,
		},
		"beijing": {
			Name:  "北京市",
			Image: "beijing.jpg",
			Description: `
Beijing, China’s sprawling capital, has history stretching back 3 millennia. Yet
it’s known as much for modern architecture as its ancient sites such as the
grand Forbidden City complex, the imperial palace during the Ming and Qing
dynasties. Nearby, the massive Tiananmen Square pedestrian plaza is the site of
Mao Zedong’s mausoleum and the National Museum of China, displaying a vast
collection of cultural relics.
			`,
		},
	}

	current, ok := n.Store[n.City]
	if !ok {
		current = n.Store["paris"]
	}
	n.Current = current

	app.Render(n)
}

// Render returns what to display.
func (n *City) Render() string {
	return `
<div class="City" style="background-image: url('{{.Current.Image}}')">
	<button class="Menu" onclick="OnMenuClick" oncontextmenu="OnMenuClick">☰</button>

	<main>
		<h1>{{.Current.Name}}</h1>
		<p>{{.Current.Description}}</p>
		<p class="Links">
			{{range $k, $v := .Store}}
				<a class="app_button" href="city?city={{$k}}">{{$v.Name}}</a>
			{{end}}
			<br>
		</p>
	</main>
</div>
	`
}

// OnMenuClick creates a context menu when the menu button is clicked.
func (n *City) OnMenuClick(s, e js.Value) {
	app.NewContextMenu(
		app.MenuItem{
			Label: "Reload",
			Keys:  "cmdorctrl+r",
			OnClick: func(s, e js.Value) {
				app.Reload()
			},
		},
		app.MenuItem{Separator: true},
		app.MenuItem{
			Label: "Go to repository",
			OnClick: func(s, e js.Value) {
				app.Navigate("https://github.com/maxence-charriere/app")
			}},
		app.MenuItem{
			Label: "Source code",
			OnClick: func(s, e js.Value) {
				app.Navigate("https://github.com/maxence-charriere/app/blob/master/demo/cmd/demo-wasm/hello.go")
			}},
		app.MenuItem{Separator: true},
		app.MenuItem{
			Label: "Hello example",
			OnClick: func(s, e js.Value) {
				app.Navigate("hello")
			}},
	)
}

type city struct {
	Name        string
	Description string
	Image       string
}
