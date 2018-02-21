package main

import (
	"net/url"

	"github.com/murlokswarm/app"
)

// /!\ Import the component. Required to use a component.
func init() {
	app.Import(&City{})
}

type city struct {
	ID          string
	Name        string
	Description string
}

var cities = map[string]city{
	"paris": city{
		ID:   "paris",
		Name: "Paris",
		Description: `Paris, France's capital, is a major European city and a 
		global center for art, fashion, gastronomy and culture.
		Its 19th-century cityscape is crisscrossed by wide boulevards and the
		River Seine.
		Beyond such landmarks as the Eiffel Tower and the 12th-century, Gothic
		Notre-Dame cathedral, the city is known for its cafe culture and
		designer boutiques along the Rue du Faubourg Saint-Honoré.`,
	},
	"sf": city{
		ID:   "sf",
		Name: "San Francisco",
		Description: `San Francisco, in northern California, is a hilly city on 
		the tip of a peninsula surrounded by the Pacific Ocean and San Francisco 
		Bay.
		It's known for its year-round fog, iconic Golden Gate Bridge, cable cars
		and colorful Victorian houses.
		The Financial District's Transamerica Pyramid is its most distinctive
		skyscraper.
		In the bay sits Alcatraz Island, site of the notorious former prison.`,
	},
	"beijing": city{
		ID:   "beijing",
		Name: "北京",
		Description: `Beijing formerly romanized as Peking, is the capital of 
		the People's Republic of China, the world's second most populous city 
		proper, and most populous capital city. 
		The city, located in northern China, is governed as a direct-controlled
		municipality under the national government with 16 urban, suburban, and 
		rural districts.`,
	},
}

// City is the component displaying a city.
type City struct {
	City        city
	CanPrevious bool
	CanNext     bool
}

// OnNavigate is the function that is call when the component is navigated to.
// It satisfies the app.Navigable interfaces.
func (c *City) OnNavigate(u *url.URL) {
	id := u.Query().Get("id")
	if len(id) == 0 {
		id = "paris"
	}

	if w, err := app.WindowByComponent(c); err == nil {
		c.CanPrevious = w.CanPrevious()
		c.CanNext = w.CanNext()
	}

	c.City = cities[id]
	app.Render(c)
}

// Render returns the HTML describing the Paris component content.
// It contains a link to show how to navigate to an other component (Sf).
func (c *City) Render() string {
	return `
<div class="Content {{.City.ID}}">
	<h1>{{.City.Name}}</h1>
	<p>{{.City.Description}}</p>
	<div>
		<a href="/city?id=sf" class="button">San Francisco</a>
		<a href="/city?id=paris" class="button">Paris</a>
		<a href="/city?id=beijing" class="button">北京"</a>		
	</div>
	<div>
		<button class="button navButton" onclick="OnPrevious" {{if not .CanPrevious}}disabled{{end}} >Previous</button>
		<button class="button navButton" onclick="OnNext" {{if not .CanNext}}disabled{{end}} >Next</button>
	</div>
</div>
	`
}

// OnPrevious is the function that is called when the button labelled "Previous"
// is clicked.
func (c *City) OnPrevious() {
	w, err := app.WindowByComponent(c)
	if err != nil {
		app.Error(err)
		return
	}

	if err = w.Previous(); err != nil {
		app.Error(err)
	}
}

// OnNext is the function that is called when the button labelled "Next" is
// clicked.
func (c *City) OnNext() {
	w, err := app.WindowByComponent(c)
	if err != nil {
		app.Error(err)
		return
	}

	if err = w.Next(); err != nil {
		app.Error(err)
	}
}
