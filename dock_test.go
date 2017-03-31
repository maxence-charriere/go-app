package app

type dockContext struct {
	*testContext
}

func newDockContext() *dockContext {
	return &dockContext{
		testContext: newTestContext("dock"),
	}
}

func (d *dockContext) SetIcon(path string) {}

func (d *dockContext) SetBadge(v interface{}) {}
