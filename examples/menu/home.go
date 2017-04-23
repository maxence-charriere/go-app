package main

import "github.com/murlokswarm/app"

// Home is the component that displays the menu examples.
type Home struct{}

// Render returns return the HTML describing the home screen.
func (h *Home) Render() string {
	return `
<div class="Home">
	<div class="Example">
		<h1>Copy/paste</h1>
		<ul oncontextmenu="OnContextMenu">
			<li>Select me</li>
			<li>Right click</li>
			<li>Copy</li>
		</ul>
		<textarea placeholder="Right click/Paste or use meta + v" oncontextmenu="OnContextMenu"></textarea>
	</div>

	<div class="Example">
		<h1>Custom menu</h1>
		<button onclick="OnButtonClick">Show</button>
	</div>
</div>
	`
}

// OnContextMenu is called when there is a right click on the ul or textarea.
// It creates a context menu and mount the Edit component inside.
func (h *Home) OnContextMenu() {
	ctxMenu := app.NewContextMenu()
	ctxMenu.Mount(&EditMenu{})
}

// OnButtonClick is called when the Show buttton is clicked.
// It creates a context menu and mount the CustomMenu component inside.
func (h *Home) OnButtonClick() {
	ctxMenu := app.NewContextMenu()
	ctxMenu.Mount(&CustomMenu{})
}

// /!\ Register the component. Required to use the component into a context.
func init() {
	app.RegisterComponent(&Home{})
}
