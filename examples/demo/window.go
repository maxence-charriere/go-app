package main

import (
	"github.com/murlokswarm/app"
)

func init() {
	app.Handle("window-update-info", func(e app.Emitter, m app.Msg) {
		e.Emit("window-updated", m.Value())
	})
}

func newWindow(title, url string, frosted bool) {
	updateInfo := func(w app.Window) {
		app.NewMsg("window-update-info").WithValue(w.ID()).Post()
	}

	app.NewWindow(app.WindowConfig{
		Title:             title,
		URL:               url,
		Width:             1440,
		MinWidth:          1024,
		Height:            720,
		MinHeight:         500,
		FrostedBackground: frosted,

		OnMove:           updateInfo,
		OnResize:         updateInfo,
		OnFocus:          updateInfo,
		OnBlur:           updateInfo,
		OnFullScreen:     updateInfo,
		OnExitFullScreen: updateInfo,
		OnMinimize:       updateInfo,
		OnDeminimize:     updateInfo,
		OnClose: func(w app.Window) {
			app.Logf("window %q is closed", w.ID())
		},
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
func (w *Window) Subscribe() app.Subscriber {
	return app.NewSubscriber().
		Subscribe("window-updated", w.onUpdate)
}

func (w *Window) onUpdate(id string) {
	app.ElemByCompo(w).WhenWindow(func(win app.Window) {
		if win.ID() != id {
			return
		}

		w.ID = app.ElemByCompo(w).ID()
		w.X, w.Y = win.Position()
		w.Width, w.Height = win.Size()
		w.IsFocus = win.IsFocus()
		w.IsFullScreen = win.IsFullScreen()
		w.IsMinified = win.IsMinimized()

		app.Render(w)
	})
}

// OnMount initializes the available actions.
func (w *Window) OnMount() {
	app.ElemByCompo(w).WhenWindow(func(win app.Window) {
		w.onUpdate(win.ID())

		w.Actions = []windowAction{
			{
				Name:        "Move",
				Description: "Move the window to position {x: 100, y: 100}.",
				Action: func() {
					win.Move(100, 100)
					w.Actions[0].Err = win.Err()
					app.Render(w)
				},
			},
			{
				Name:        "Center",
				Description: "Move the window to the center of the screen.",
				Action: func() {
					win.Center()
					w.Actions[1].Err = win.Err()
					app.Render(w)
				},
			},
			{
				Name:        "Resize",
				Description: "Resize the window to its original size {width: 1440, height: 720}.",
				Action: func() {
					win.Resize(1440, 720)
					w.Actions[2].Err = win.Err()
					app.Render(w)
				},
			},
			{
				Name:        "FullScreen",
				Description: "Take the window in full screen mode.",
				Action: func() {
					win.FullScreen()
					w.Actions[3].Err = win.Err()
					app.Render(w)
				},
			},
			{
				Name:        "ExitFullScreen",
				Description: "Take the window out of fullscreen mode",
				Action: func() {
					win.ExitFullScreen()
					w.Actions[4].Err = win.Err()
					app.Render(w)
				},
			},
			{
				Name:        "Minimize",
				Description: "Take the window into minimized mode.",
				Action: func() {
					win.Minimize()
					w.Actions[5].Err = win.Err()
					app.Render(w)
				},
			},
			{
				Name:        "Minimize/Deminimize",
				Description: "Take and take out the window out of minimized mode.",
				Action: func() {
					win.Minimize()

					if win.Err() != nil {
						w.Actions[6].Err = win.Err()
						return
					}

					app.CallOnUIGoroutine(func() {
						win.Deminimize()
						w.Actions[6].Err = win.Err()
						app.Render(w)
					})

					app.Render(w)
				},
			},
			{
				Name:        "Frosted window",
				Description: "Create a window with frosted effect.",
				Action: func() {
					newWindow("frosted", "window", true)
					w.Actions[7].Err = win.Err()
				},
			},
			{
				Name:        "Close",
				Description: "Close the window.",
				Action: func() {
					win.Close()
					w.Actions[8].Err = win.Err()
				},
			},
		}

		app.Render(w)
	})
}

// Render returns a html string that describes the component.
func (w *Window) Render() string {
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
