package main

import "github.com/murlokswarm/app"

// Paris is the component displaying Paris.
type Paris struct{}

// Render returns the HTML describing the Paris component content.
// It contains a link to show how to navigate to an other component (Sf).
func (p *Paris) Render() string {
	return `
<div class="Content Paris">
	<h1>Paris</h1>
	<p>
		Paris, France's capital, is a major European city and a global center 
		for art, fashion, gastronomy and culture. 
		Its 19th-century cityscape is crisscrossed by wide boulevards and the 
		River Seine. 
		Beyond such landmarks as the Eiffel Tower and the 12th-century, Gothic 
		Notre-Dame cathedral, the city is known for its cafe culture and 
		designer boutiques along the Rue du Faubourg Saint-Honoré.
	</p>
	<a href="Sf">Go to San Francisco</a>
</div>
	`
}

// /!\ Register the component. Required to use the component into a context.
func init() {
	app.RegisterComponent(&Paris{})
}
