package main

// import (
// 	"fmt"
// 	"net/url"

// 	"github.com/murlokswarm/app"
// 	"github.com/murlokswarm/app/markup"
// )

// // Sf is the component displaying San Francisco.
// type Sf markup.ZeroCompo

// // Render returns the HTML describing the Sf component content.
// // It contains a link to show how to navigate to an other component (Paris).
// func (s *Sf) Render() string {
// 	return `
// <div class="Content Sf">
// 	<h1>San Francisco</h1>
// 	<p>
// 		San Francisco, in northern California, is a hilly city on the tip of a
// 		peninsula surrounded by the Pacific Ocean and San Francisco Bay.
// 		It's known for its year-round fog, iconic Golden Gate Bridge, cable cars
// 		and colorful Victorian houses.
// 		The Financial District's Transamerica Pyramid is its most distinctive
// 		skyscraper.
// 		In the bay sits Alcatraz Island, site of the notorious former prison.
// 	</p>
// 	<a href="Paris">Go to Paris</a>
// </div>
// 	`
// }

// // OnHref is defined to satisfy the Hrefer interface. It is called when a link
// // with href="Sf" is clicked.
// func (s *Sf) OnHref(URL *url.URL) {
// 	fmt.Println("mounted from a link click:", URL)
// }

// // /!\ Import the component. Required to use a component.
// func init() {
// 	app.Import(&Sf{})
// }
