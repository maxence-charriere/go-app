package app

import "testing"

func TestRenderPanic(t *testing.T) {
	defer func() { recover() }()

	hello := &Hello{}
	render(hello)

	t.Error("should panic")
}
