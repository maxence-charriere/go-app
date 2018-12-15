package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

type appOpen struct {
	From string
	Time time.Time
}

type appOpenings struct {
	mutex    sync.Mutex
	openings []appOpen
}

func (o *appOpenings) Add(open appOpen) {
	o.mutex.Lock()
	o.openings = append(o.openings, open)
	o.mutex.Unlock()
}

func (o *appOpenings) Openings() []appOpen {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if len(o.openings) == 0 {
		return nil
	}

	openings := make([]appOpen, len(o.openings))
	copy(openings, o.openings)
	return openings
}

var (
	appOpenList appOpenings
)

func init() {
	app.Handle("app-open", func(e app.Emitter, m app.Msg) {
		open, ok := m.Value().(appOpen)
		if !ok {
			app.Log(errors.Errorf("msg value for %q is not a %T: %T", m.Key(), open, m.Value()))
			return
		}

		appOpenList.Add(open)
		fmt.Println(app.Pretty(appOpenList.Openings()))
		e.Emit("app-opened", appOpenList.Openings())
	})
}

// Open is a component that shows app opening behavior.
type Open struct {
	Openings []appOpen
}

// Subscribe is the func to set up event listeners.
// It satisfies the app.EventSubscriber interface.
func (o *Open) Subscribe() app.Subscriber {
	return app.NewSubscriber().Subscribe("app-opened", o.onAppOpen)
}

func (o *Open) onAppOpen(openings []appOpen) {
	o.Openings = openings
	app.Render(o)
}

// OnMount initializes the component openings.
func (o *Open) OnMount() {
	o.Openings = appOpenList.Openings()
	app.Render(o)
}

// Render returns a html string that describes the component.
func (o *Open) Render() string {
	return `
<div class="Layout">
	<navpane current="open">
	<div class="Open">
		<h1>Open</h1>
		<div class="Open-List">
			<table>
				<tr>
					<th>Time</th>
					<th>From</th>
				</tr>

				{{range .Openings}}
				<tr>
					<td>{{time .Time "15:04:05"}}</td>
					<td>{{.From}}</td>
				</tr>
				{{end}}
			</table>
		</div>
	</div>
</div>
	`
}
