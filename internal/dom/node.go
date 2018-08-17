package dom

import "github.com/murlokswarm/app"

type node interface {
	app.Node

	Close()
	CompoID() string
	SetParent(node)
}
