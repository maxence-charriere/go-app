package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type menu struct {
	app.Compo

	currentPath string
}

func (m *menu) OnNav(ctx app.Context) {
	path := ctx.Page.URL().Path
	if path == "/" {
		path = "/start"
	}
	m.currentPath = path

	m.Update()
}

func (m *menu) Render() app.UI {
	return app.Nav().
		Class("menu").
		Body(
			app.Div().Body(
				app.A().
					Class("title").
					Href("/start").
					Text("go-app"),
			),
			app.Div().
				Class("content").
				Body(
					app.Section().Body(
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M13.13 22.19L11.5 18.36C13.07 17.78 14.54 17 15.9 16.09L13.13 22.19M5.64 12.5L1.81 10.87L7.91 8.1C7 9.46 6.22 10.93 5.64 12.5M21.61 2.39C21.61 2.39 16.66 .269 11 5.93C8.81 8.12 7.5 10.53 6.65 12.64C6.37 13.39 6.56 14.21 7.11 14.77L9.24 16.89C9.79 17.45 10.61 17.63 11.36 17.35C13.5 16.53 15.88 15.19 18.07 13C23.73 7.34 21.61 2.39 21.61 2.39M14.54 9.46C13.76 8.68 13.76 7.41 14.54 6.63S16.59 5.85 17.37 6.63C18.14 7.41 18.15 8.68 17.37 9.46C16.59 10.24 15.32 10.24 14.54 9.46M8.88 16.53L7.47 15.12L8.88 16.53M6.24 22L9.88 18.36C9.54 18.27 9.21 18.12 8.91 17.91L4.83 22H6.24M2 22H3.41L8.18 17.24L6.76 15.83L2 20.59V22M2 19.17L6.09 15.09C5.88 14.79 5.73 14.47 5.64 14.12L2 17.76V19.17Z" />
							</svg>
							`).
							Text("Getting started").
							Selected(m.currentPath == "/start").
							Href("/start"),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M9,2V8H11V11H5C3.89,11 3,11.89 3,13V16H1V22H7V16H5V13H11V16H9V22H15V16H13V13H19V16H17V22H23V16H21V13C21,11.89 20.11,11 19,11H13V8H15V2H9Z" />
							</svg>
							`).
							Text("Architecture").
							Selected(m.currentPath == "/architecture").
							Href("/architecture"),
						newMenuItem().
							Icon(`
							<svg height="24px" viewBox="0 0 207 78" width="24px" xmlns="http://www.w3.org/2000/svg">
								<g fill="currentColor" fill-rule="evenodd">
									<path d="m16.2 24.1c-.4 0-.5-.2-.3-.5l2.1-2.7c.2-.3.7-.5 1.1-.5h35.7c.4 0 .5.3.3.6l-1.7 2.6c-.2.3-.7.6-1 .6z"/>
									<path d="m1.1 33.3c-.4 0-.5-.2-.3-.5l2.1-2.7c.2-.3.7-.5 1.1-.5h45.6c.4 0 .6.3.5.6l-.8 2.4c-.1.4-.5.6-.9.6z"/>
									<path d="m25.3 42.5c-.4 0-.5-.3-.3-.6l1.4-2.5c.2-.3.6-.6 1-.6h20c.4 0 .6.3.6.7l-.2 2.4c0 .4-.4.7-.7.7z"/>
									<g transform="translate(55)">
										<path d="m74.1 22.3c-6.3 1.6-10.6 2.8-16.8 4.4-1.5.4-1.6.5-2.9-1-1.5-1.7-2.6-2.8-4.7-3.8-6.3-3.1-12.4-2.2-18.1 1.5-6.8 4.4-10.3 10.9-10.2 19 .1 8 5.6 14.6 13.5 15.7 6.8.9 12.5-1.5 17-6.6.9-1.1 1.7-2.3 2.7-3.7-3.6 0-8.1 0-19.3 0-2.1 0-2.6-1.3-1.9-3 1.3-3.1 3.7-8.3 5.1-10.9.3-.6 1-1.6 2.5-1.6h36.4c-.2 2.7-.2 5.4-.6 8.1-1.1 7.2-3.8 13.8-8.2 19.6-7.2 9.5-16.6 15.4-28.5 17-9.8 1.3-18.9-.6-26.9-6.6-7.4-5.6-11.6-13-12.7-22.2-1.3-10.9 1.9-20.7 8.5-29.3 7.1-9.3 16.5-15.2 28-17.3 9.4-1.7 18.4-.6 26.5 4.9 5.3 3.5 9.1 8.3 11.6 14.1.6.9.2 1.4-1 1.7z"/>
										<path d="m107.2 77.6c-9.1-.2-17.4-2.8-24.4-8.8-5.9-5.1-9.6-11.6-10.8-19.3-1.8-11.3 1.3-21.3 8.1-30.2 7.3-9.6 16.1-14.6 28-16.7 10.2-1.8 19.8-.8 28.5 5.1 7.9 5.4 12.8 12.7 14.1 22.3 1.7 13.5-2.2 24.5-11.5 33.9-6.6 6.7-14.7 10.9-24 12.8-2.7.5-5.4.6-8 .9zm23.8-40.4c-.1-1.3-.1-2.3-.3-3.3-1.8-9.9-10.9-15.5-20.4-13.3-9.3 2.1-15.3 8-17.5 17.4-1.8 7.8 2 15.7 9.2 18.9 5.5 2.4 11 2.1 16.3-.6 7.9-4.1 12.2-10.5 12.7-19.1z" fill-rule="nonzero"/>
									</g>
								</g>
							</svg>
							`).
							Text("API reference").
							Selected(m.currentPath == "/reference").
							Href("/reference"),
					),
					app.Section().Body(
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M10,5V11H21V5M16,18H21V12H16M4,18H9V5H4M10,18H15V12H10V18Z" />
							</svg>
							`).
							Text("Components").
							Selected(m.currentPath == "/components").
							Href("/components"),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M16,4.5V7H5V9H16V11.5L19.5,8M16,12.5V15H5V17H16V19.5L19.5,16" />
							</svg>
							`).
							Text("Concurrency").
							Selected(m.currentPath == "/concurrency").
							Href("/concurrency"),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M19,10H17V8H19M19,13H17V11H19M16,10H14V8H16M16,13H14V11H16M16,17H8V15H16M7,10H5V8H7M7,13H5V11H7M8,11H10V13H8M8,8H10V10H8M11,11H13V13H11M11,8H13V10H11M20,5H4C2.89,5 2,5.89 2,7V17A2,2 0 0,0 4,19H20A2,2 0 0,0 22,17V7C22,5.89 21.1,5 20,5Z" />
							</svg>
							`).
							Text("Declarative syntax").
							Selected(m.currentPath == "/syntax").
							Href("/syntax"),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M3,3H21V21H3V3M7.73,18.04C8.13,18.89 8.92,19.59 10.27,19.59C11.77,19.59 12.8,18.79 12.8,17.04V11.26H11.1V17C11.1,17.86 10.75,18.08 10.2,18.08C9.62,18.08 9.38,17.68 9.11,17.21L7.73,18.04M13.71,17.86C14.21,18.84 15.22,19.59 16.8,19.59C18.4,19.59 19.6,18.76 19.6,17.23C19.6,15.82 18.79,15.19 17.35,14.57L16.93,14.39C16.2,14.08 15.89,13.87 15.89,13.37C15.89,12.96 16.2,12.64 16.7,12.64C17.18,12.64 17.5,12.85 17.79,13.37L19.1,12.5C18.55,11.54 17.77,11.17 16.7,11.17C15.19,11.17 14.22,12.13 14.22,13.4C14.22,14.78 15.03,15.43 16.25,15.95L16.67,16.13C17.45,16.47 17.91,16.68 17.91,17.26C17.91,17.74 17.46,18.09 16.76,18.09C15.93,18.09 15.45,17.66 15.09,17.06L13.71,17.86Z" />
							</svg>
							`).
							Text("JS/Dom").
							Selected(m.currentPath == "/js").
							Href("/js"),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M11,10H5L3,8L5,6H11V3L12,2L13,3V4H19L21,6L19,8H13V10H19L21,12L19,14H13V20A2,2 0 0,1 15,22H9A2,2 0 0,1 11,20V10Z" />
							</svg>
							`).
							Text("Lifecycle").
							Selected(m.currentPath == "/lifecycle").
							Href("/lifecycle"),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M11,10H5L3,8L5,6H11V3L12,2L13,3V4H19L21,6L19,8H13V10H19L21,12L19,14H13V20A2,2 0 0,1 15,22H9A2,2 0 0,1 11,20V10Z" />
							</svg>
							`).
							Text("Routing").
							Selected(m.currentPath == "/routing").
							Href("/routing"),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
								<path fill="currentColor" d="M15,7H20.5L15,1.5V7M8,0H16L22,6V18A2,2 0 0,1 20,20H8C6.89,20 6,19.1 6,18V2A2,2 0 0,1 8,0M4,4V22H20V24H4A2,2 0 0,1 2,22V4H4Z" />
							</svg>
							`).
							Text("Static resources").
							Selected(m.currentPath == "/static-resources").
							Href("/static-resources"),
						// newMenuItem().
						// 	Icon(`
						// 	<svg style="width:24px;height:24px" viewBox="0 0 24 24">
						// 		<path fill="currentColor" d="M3,3H11V7.34L16.66,1.69L22.31,7.34L16.66,13H21V21H13V13H16.66L11,7.34V11H3V3M3,13H11V21H3V13Z" />
						// 	</svg>
						// 	`).
						// 	Text("Widgets").
						// 	Selected(m.currentPath == "/widgets").
						// 	Href("/widgets"),

					),
					app.Section().Body(
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
								<path fill="currentColor" d="M2 19.63L13.43 8.2L12.72 7.5L14.14 6.07L12 3.89C13.2 2.7 15.09 2.7 16.27 3.89L19.87 7.5L18.45 8.91H21.29L22 9.62L18.45 13.21L17.74 12.5V9.62L16.27 11.04L15.56 10.33L4.13 21.76L2 19.63Z" />
							</svg>
							`).
							Text("Built with go-app").
							Selected(m.currentPath == "/built-with").
							Href("/built-with"),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M5,20H19V18H5M19,9H15V3H9V9H5L12,16L19,9Z" />
							</svg>
							`).
							Text("Install").
							Selected(m.currentPath == "/install").
							Href("/install"),
					),
					app.Section().Body(
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
								<path fill="currentColor" d="M2,21H20V19H2M20,8H18V5H20M20,3H4V13A4,4 0 0,0 8,17H14A4,4 0 0,0 18,13V10H20A2,2 0 0,0 22,8V5C22,3.89 21.1,3 20,3Z" />
							</svg>
							`).
							Text("Buy me a coffee").
							Href(buyMeACoffeeURL).
							External(),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
								<path fill="currentColor" d="M15.41,22C15.35,22 15.28,22 15.22,22C15.1,21.95 15,21.85 14.96,21.73L12.74,15.93C12.65,15.69 12.77,15.42 13,15.32C13.71,15.06 14.28,14.5 14.58,13.83C15.22,12.4 14.58,10.73 13.15,10.09C11.72,9.45 10.05,10.09 9.41,11.5C9.11,12.21 9.09,13 9.36,13.69C9.66,14.43 10.25,15 11,15.28C11.24,15.37 11.37,15.64 11.28,15.89L9,21.69C8.96,21.81 8.87,21.91 8.75,21.96C8.63,22 8.5,22 8.39,21.96C3.24,19.97 0.67,14.18 2.66,9.03C4.65,3.88 10.44,1.31 15.59,3.3C18.06,4.26 20.05,6.15 21.13,8.57C22.22,11 22.29,13.75 21.33,16.22C20.32,18.88 18.23,21 15.58,22C15.5,22 15.47,22 15.41,22M12,3.59C7.03,3.46 2.9,7.39 2.77,12.36C2.68,16.08 4.88,19.47 8.32,20.9L10.21,16C8.38,15 7.69,12.72 8.68,10.89C9.67,9.06 11.96,8.38 13.79,9.36C15.62,10.35 16.31,12.64 15.32,14.47C14.97,15.12 14.44,15.65 13.79,16L15.68,20.93C17.86,19.95 19.57,18.16 20.44,15.93C22.28,11.31 20.04,6.08 15.42,4.23C14.33,3.8 13.17,3.58 12,3.59Z" />
							</svg>
							`).
							Text("Open Collective").
							Href(openCollectiveURL).
							External(),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M12,2A10,10 0 0,0 2,12C2,16.42 4.87,20.17 8.84,21.5C9.34,21.58 9.5,21.27 9.5,21C9.5,20.77 9.5,20.14 9.5,19.31C6.73,19.91 6.14,17.97 6.14,17.97C5.68,16.81 5.03,16.5 5.03,16.5C4.12,15.88 5.1,15.9 5.1,15.9C6.1,15.97 6.63,16.93 6.63,16.93C7.5,18.45 8.97,18 9.54,17.76C9.63,17.11 9.89,16.67 10.17,16.42C7.95,16.17 5.62,15.31 5.62,11.5C5.62,10.39 6,9.5 6.65,8.79C6.55,8.54 6.2,7.5 6.75,6.15C6.75,6.15 7.59,5.88 9.5,7.17C10.29,6.95 11.15,6.84 12,6.84C12.85,6.84 13.71,6.95 14.5,7.17C16.41,5.88 17.25,6.15 17.25,6.15C17.8,7.5 17.45,8.54 17.35,8.79C18,9.5 18.38,10.39 18.38,11.5C18.38,15.32 16.04,16.16 13.81,16.41C14.17,16.72 14.5,17.33 14.5,18.26C14.5,19.6 14.5,20.68 14.5,21C14.5,21.27 14.66,21.59 15.17,21.5C19.14,20.16 22,16.42 22,12A10,10 0 0,0 12,2Z" />
							</svg>
							`).
							Text("GitHub").
							Href(githubURL).
							External(),
						newMenuItem().
							Icon(`
							<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    							<path fill="currentColor" d="M22.46,6C21.69,6.35 20.86,6.58 20,6.69C20.88,6.16 21.56,5.32 21.88,4.31C21.05,4.81 20.13,5.16 19.16,5.36C18.37,4.5 17.26,4 16,4C13.65,4 11.73,5.92 11.73,8.29C11.73,8.63 11.77,8.96 11.84,9.27C8.28,9.09 5.11,7.38 3,4.79C2.63,5.42 2.42,6.16 2.42,6.94C2.42,8.43 3.17,9.75 4.33,10.5C3.62,10.5 2.96,10.3 2.38,10C2.38,10 2.38,10 2.38,10.03C2.38,12.11 3.86,13.85 5.82,14.24C5.46,14.34 5.08,14.39 4.69,14.39C4.42,14.39 4.15,14.36 3.89,14.31C4.43,16 6,17.26 7.89,17.29C6.43,18.45 4.58,19.13 2.56,19.13C2.22,19.13 1.88,19.11 1.54,19.07C3.44,20.29 5.7,21 8.12,21C16,21 20.33,14.46 20.33,8.79C20.33,8.6 20.33,8.42 20.32,8.23C21.16,7.63 21.88,6.87 22.46,6Z" />
							</svg>
							`).
							Text("Twitter").
							Href(twitterURL).
							External(),
					),
				),
		)
}

type menuItem struct {
	app.Compo

	Iicon     string
	Ihref     string
	Iselected string
	Itext     string
	Itarget   string
	Irel      string
}

func newMenuItem() *menuItem {
	return &menuItem{}
}

func (i *menuItem) Icon(svg string) *menuItem {
	i.Iicon = svg
	return i
}

func (i *menuItem) Href(v string) *menuItem {
	i.Ihref = v
	return i
}

func (i *menuItem) Selected(v bool) *menuItem {
	if v {
		i.Iselected = "focus"
	}
	return i
}

func (i *menuItem) Text(v string) *menuItem {
	i.Itext = v
	return i
}

func (i *menuItem) External() *menuItem {
	i.Itarget = "_blank"
	i.Irel = "noopener"
	return i
}

func (i *menuItem) Render() app.UI {
	return app.A().
		Class("item").
		Class(i.Iselected).
		Href(i.Ihref).
		Target(i.Itarget).
		Rel(i.Irel).
		Body(
			app.Stack().
				Center().
				Content(
					app.Div().
						Class("icon").
						Body(
							app.Raw(i.Iicon),
						),
					app.Div().
						Class("label").
						Text(i.Itext),
				),
		)
}

type overlayMenu struct {
	app.Compo
}

func (m *overlayMenu) Render() app.UI {
	return app.Div().
		Class("overlay-menu").
		Body(
			&menu{},
		)
}
