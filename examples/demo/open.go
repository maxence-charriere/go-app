package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

const (
	appOpen             = "appOpen"
	appOpened app.Event = "appOpened"
)

var (
	appOpenList appOpenings
)

type appOpenInfo struct {
	From string
	Time time.Time
}

type appOpenings struct {
	mutex    sync.Mutex
	openings []appOpenInfo
}

func (o *appOpenings) Add(open appOpenInfo) {
	o.mutex.Lock()
	o.openings = append(o.openings, open)
	o.mutex.Unlock()
}

func (o *appOpenings) Openings() []appOpenInfo {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if len(o.openings) == 0 {
		return nil
	}

	openings := make([]appOpenInfo, len(o.openings))
	copy(openings, o.openings)
	return openings
}

func init() {
	app.NewSubscriber().
		Subscribe(app.Running, func() {
			app.NewMsg(appOpen).
				WithValue(appOpenInfo{
					From: string(app.Running),
					Time: time.Now(),
				}).
				Post()
		}).
		Subscribe(app.Reopened, func(hasWindows bool) {
			app.NewMsg(appOpen).
				WithValue(appOpenInfo{
					From: string(app.Reopened),
					Time: time.Now(),
				}).
				Post()
		}).
		Subscribe(app.OpenFilesRequested, func(filenames []string) {
			app.NewMsg(appOpen).
				WithValue(appOpenInfo{
					From: fmt.Sprintf("%s(%v)", app.OpenFilesRequested, app.Pretty(filenames)),
					Time: time.Now(),
				}).
				Post()
		}).
		Subscribe(app.OpenURLRequested, func(u *url.URL) {
			app.NewMsg(appOpen).
				WithValue(appOpenInfo{
					From: fmt.Sprintf("%s(%s)", app.OpenURLRequested, u),
					Time: time.Now(),
				}).
				Post()
		})

	app.Handle(appOpen, func(m app.Msg) {
		open, ok := m.Value().(appOpenInfo)
		if !ok {
			app.Log(errors.Errorf("msg value for %q is not a %T: %T", m.Key(), open, m.Value()))
			return
		}

		appOpenList.Add(open)
		app.Emit(appOpened, appOpenList.Openings())
	})
}

// Open is a component that shows app opening behavior.
type Open struct {
	Openings []appOpenInfo
}

// Subscribe is the func to set up event listeners.
// It satisfies the app.EventSubscriber interface.
func (o *Open) Subscribe() *app.Subscriber {
	return app.NewSubscriber().Subscribe(appOpened, o.onAppOpen)
}

func (o *Open) onAppOpen(openings []appOpenInfo) {
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
