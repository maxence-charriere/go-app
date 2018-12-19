package main

import (
	"github.com/murlokswarm/app"
)

func newWindow(title, url string, frosted bool) {
	app.NewWindow(app.WindowConfig{
		Title:             title,
		URL:               url,
		Width:             1440,
		MinWidth:          1024,
		Height:            720,
		MinHeight:         500,
		FrostedBackground: frosted,
	})
}

// Window is a component that contain window related examples.
type Window struct {
	ID           string
	X            float64
	Y            float64
	Width        float64
	Height       float64
	IsFocus      bool
	IsFullScreen bool
	IsMinified   bool
	Actions      []windowAction
}

// Subscribe is the func to set up event listeners.
// It satisfies the app.EventSubscriber interface.
func (c *Window) Subscribe() app.Subscriber {
	return app.NewSubscriber().
		Subscribe(app.WindowMoved, c.onUpdate).
		Subscribe(app.WindowResized, c.onUpdate).
		Subscribe(app.WindowFocused, c.onUpdate).
		Subscribe(app.WindowBlurred, c.onUpdate).
		Subscribe(app.WindowEnteredFullScreen, c.onUpdate).
		Subscribe(app.WindowExitedFullScreen, c.onUpdate).
		Subscribe(app.WindowMinimized, c.onUpdate).
		Subscribe(app.WindowDeminimized, c.onUpdate)
}

func (c *Window) onUpdate(e app.Window) {
	if !e.Contains(c) {
		return
	}

	app.ElemByCompo(c).WhenWindow(func(w app.Window) {
		c.ID = app.ElemByCompo(c).ID()
		c.X, c.Y = w.Position()
		c.Width, c.Height = w.Size()
		c.IsFocus = w.IsFocus()
		c.IsFullScreen = w.IsFullScreen()
		c.IsMinified = w.IsMinimized()

		app.Render(c)
	})
}

// OnMount initializes the available actions.
func (c *Window) OnMount() {
	app.ElemByCompo(c).WhenWindow(func(w app.Window) {
		c.onUpdate(w)

		c.Actions = []windowAction{
			{
				Name:        "Move",
				Description: "Move the window to position {x: 100, y: 100}.",
				Action: func() {
					w.Move(100, 100)
					c.Actions[0].Err = w.Err()
					app.Render(c)
				},
			},
			{
				Name:        "Center",
				Description: "Move the window to the center of the screen.",
				Action: func() {
					w.Center()
					c.Actions[1].Err = w.Err()
					app.Render(c)
				},
			},
			{
				Name:        "Resize",
				Description: "Resize the window to its original size {width: 1440, height: 720}.",
				Action: func() {
					w.Resize(1440, 720)
					c.Actions[2].Err = w.Err()
					app.Render(c)
				},
			},
			{
				Name:        "FullScreen",
				Description: "Take the window in full screen mode.",
				Action: func() {
					w.FullScreen()
					c.Actions[3].Err = w.Err()
					app.Render(c)
				},
			},
			{
				Name:        "ExitFullScreen",
				Description: "Take the window out of fullscreen mode",
				Action: func() {
					w.ExitFullScreen()
					c.Actions[4].Err = w.Err()
					app.Render(c)
				},
			},
			{
				Name:        "Minimize",
				Description: "Take the window into minimized mode.",
				Action: func() {
					w.Minimize()
					c.Actions[5].Err = w.Err()
					app.Render(c)
				},
			},
			{
				Name:        "Minimize/Deminimize",
				Description: "Take and take out the window out of minimized mode.",
				Action: func() {
					w.Minimize()

					if w.Err() != nil {
						c.Actions[6].Err = w.Err()
						app.Render(c)
						return
					}

					app.CallOnUIGoroutine(func() {
						w.Deminimize()
						c.Actions[6].Err = w.Err()
						app.Render(c)
					})

					app.Render(c)
				},
			},
			{
				Name:        "Frosted window",
				Description: "Create a window with frosted effect.",
				Action: func() {
					newWindow("frosted", "window", true)
				},
			},
			{
				Name:        "Close",
				Description: "Close the window.",
				Action: func() {
					w.Close()
					c.Actions[8].Err = w.Err()
					app.Render(c)
				},
			},
		}

		app.Render(c)
	})
}

// Render returns a html string that describes the component.
func (c *Window) Render() string {
	return `
<div class="Layout">
	<navpane current="window">
	<div class="Window-Tracking">
		<h1>Window</h1>
		<table>
			<tr>
				<td>ID</td>
				<td>{{.ID}}</td>
			</tr>
			<tr>
				<td>X</td>
				<td>{{.X}}</td>
			</tr>
			<tr>
				<td>Y</td>
				<td>{{.Y}}</td>
			</tr>
			<tr>
				<td>Width</td>
				<td>{{.Width}}</td>
			</tr>
			<tr>
				<td>Height</td>
				<td>{{.Height}}</td>
			</tr>
			<tr>
				<td>Focus</td>
				<td>{{.IsFocus}}</td>
			</tr>
			<tr>
				<td>FullScreen</td>
				<td>{{.IsFullScreen}}</td>
			</tr>
		</table>
	</div>
	<div class="Window-Actions">
		<h1 class="TopTitle">Actions</h1>
		<div class="Window-Action-List">
			{{range $idx, $v := .Actions}}
			<div class="Window-Action" onclick="{{to "Actions" $idx "Action"}}">
				<h2>
					{{.Name}}	
				</h2>
				<p>{{.Description}}</p>

				{{if .Err}} 
					<p class="Error">{{.Err.Error}}</p>
				{{end}}
			</div>
			{{else}}
				<h2>Not supported</h2>
			{{end}}
		</div>
	</div>
</div>
	`
}

type windowAction struct {
	Name        string
	Description string
	Action      func()
	Err         error
}
