package app

import "testing"

type windowContext struct {
	*testContext
}

func newWindowContext() *windowContext {
	return &windowContext{
		testContext: newTestContext("window"),
	}
}

func (w *windowContext) Position() (x float64, y float64) {
	return
}

func (w *windowContext) Move(x float64, y float64) {}

func (w *windowContext) Size() (width float64, height float64) {
	return
}

func (w *windowContext) Resize(width float64, height float64) {}

func (w *windowContext) Close() {}

func TestNewWindow(t *testing.T) {
	w := Window{}
	t.Log(NewWindow(w))
}
