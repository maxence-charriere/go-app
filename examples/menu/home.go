package main

import "github.com/murlokswarm/app"

type Home struct{}

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

func (h *Home) OnContextMenu() {
	ctxMenu := app.NewContextMenu()
	ctxMenu.Mount(&EditMenu{})
}

func (h *Home) OnButtonClick() {
	ctxMenu := app.NewContextMenu()
	ctxMenu.Mount(&CustomMenu{})
}

// /!\ Register the component. Required to use the component into a context.
func init() {
	app.RegisterComponent(&Home{})
}
