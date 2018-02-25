package main

import (
	"github.com/murlokswarm/app"
)

// /!\ Import the component. Required to use a component.
func init() {
	app.Import(&DragDrop{})
}

// DragDrop is the component describing the Drag and Drop example.
type DragDrop struct {
	History   []string
	DragHover bool
}

// Render returns the HTML describing the Drag and Drop component.
func (d *DragDrop) Render() string {
	return `
<div class="DragAndDrop">
	<div class="DragSelect">
		<div class="DragItem" 
			 data-drag="Planet" 
			 draggable="true" 
			 ondragstart="OnDragStart">
			Planet
		</div>
		<div class="DragItem" 
			 data-drag="Spaceship" 
			 draggable="true" 
			 ondragstart="OnDragStart">
			Spaceship
		</div>
		<div class="DragItem" 
			 data-drag="Human" 
			 draggable="true" 
			 ondragstart="OnDragStart">
			Human
		</div>
	</div>

	<div class="Drop" 
		 ondragleave="OnDragLeave"
		 ondrop="OnDrop" 
		 ondragover="OnDragEnter">
		<h1>Drop something</h1>
		<img class="Blackhole {{if .DragHover}}FastRotate{{else}}Rotate{{end}}" 
			 ondragleave="js:event.preventDefault()"
			 src="blackhole.png">
	</div>

	<ul class="History">
		{{range .History}}
			<li>{{.}}</li>
		{{end}}
	</ul>
</div>
	`
}

// OnDragStart is called when draging a draggable node.
// Even if it does nothing, a component event handler must be declared in
// the ondragstart of a draggable item.
// It enables an internal mechanism that allows the node property data-drag
// to be mapped when the item is dropped.
func (d *DragDrop) OnDragStart() {
}

// OnDragEnter is called when a draggable item enters in the drop zone.
func (d *DragDrop) OnDragEnter() {
	d.DragHover = true
	app.Render(d)

}

// OnDragLeave is called when a draggable item leaves the drop zone.
func (d *DragDrop) OnDragLeave() {
	d.DragHover = false
	app.Render(d)
}

// OnDrop is the handler called when a draggable item is dropped in the drop
// zone.
func (d *DragDrop) OnDrop(e app.DragAndDropEvent) {
	// Handling external drop file.
	if len(e.Files) != 0 {
		d.History = append(e.Files, d.History...)
	}

	// Handlig drop from HTML.
	if len(e.Data) != 0 {
		d.History = append([]string{e.Data}, d.History...)
	}

	d.DragHover = false
	app.Render(d)
}
