package dom

import "github.com/murlokswarm/app"

type compo struct {
	ID       string
	rootID   string
	parentID string
	compo    app.Compo
	events   *app.EventSubscriber
}
