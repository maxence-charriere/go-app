package dom

import "github.com/murlokswarm/app"

type node interface {
	app.Node

	CompoID() string
	SetParent(node)
	Flush() []Change
	Close()
}
