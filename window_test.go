package app

import "testing"

type WindowCtx struct {
	*ZeroContext
}

func newWindowCtx() *WindowCtx {
	return &WindowCtx{
		ZeroContext: NewZeroContext("window"),
	}
}

func (w *WindowCtx) Position() (x float64, y float64) {
	return
}

func (w *WindowCtx) Move(x float64, y float64) {}

func (w *WindowCtx) Size() (width float64, height float64) {
	return
}

func (w *WindowCtx) Resize(width float64, height float64) {}

func TestNewWindow(t *testing.T) {
	w := Window{}
	t.Log(NewWindow(w))
}
