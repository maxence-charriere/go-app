package actionevent

import (
	"encoding/json"

	"github.com/murlokswarm/app"
)

func init() {
	// Set the handler to be called when an action is created.
	app.HandleAction("click-action", func(e app.EventDispatcher, a app.Action) {
		// Dispatch the event to all event subcriber.
		e.Dispatch("click-event", a.Arg)
	})
}

// Clickbox is a component that displays a clickable area that produce
// click action when clicked.
type Clickbox app.ZeroCompo

// Render returns the markup that describe the component.
func (b *Clickbox) Render() string {
	return `
<div class="Layout">
	<div class="Clickbox">
		<h1>Click area</h1>
		<div class="ClickArea" onclick="OnClick"></div>	
	</div>
	<actionevent.ClickListener>
	<actionevent.ClickListener>
</div>
	`
}

// OnClick is called when a click on the click area occurs.
func (b *Clickbox) OnClick(e app.MouseEvent) {
	// Create a new action.
	app.PostAction("click-action", e)
}

// ClickListener is a component that listen for click-action and display
// click info.
type ClickListener struct {
	Logs []string
}

// Subscribe satisfie the app.Subscriber interface.
// It is where event subscription have to be setup.
func (l *ClickListener) Subscribe() *app.EventSubscriber {
	// Returns the subscriber.
	// No need to close/unsubscribe, this is internally handled.
	// No memory leak here!
	return app.NewEventSubscriber().
		Subscribe("click-event", l.OnClickEvent)

}

// Render returns the markup that describe the component.
func (l *ClickListener) Render() string {
	return `
<div class="ClickListener">
	<h1>Click Listener</h1>
	<div class="ClickOutput">
		{{range .Logs}}
			<p>{{.}}</p>
		{{end}}
	</div>
</div>d
	`
}

// OnClickEvent is the function that is called when a click-event is dispatched
// fron the click-action handler.
func (l *ClickListener) OnClickEvent(e app.MouseEvent) {
	d, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		app.Logf("on click event: %s", err)
		return
	}

	l.Logs = append([]string{string(d)}, l.Logs...)
	app.Render(l)
}
