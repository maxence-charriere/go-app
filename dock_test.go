package app

type DockCtx struct {
	*ZeroContext
}

func newDockCtx() *DockCtx {
	return &DockCtx{
		ZeroContext: NewZeroContext("dock"),
	}
}

func (d *DockCtx) SetIcon(path string) {}

func (d *DockCtx) SetBadge(v interface{}) {}
