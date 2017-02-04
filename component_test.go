package app

import "testing"

func TestComponentID(t *testing.T) {
	c := &Hello{}
	ctx := NewWindow(Window{})
	ctx.Mount(c)
	defer ctx.Close()

	t.Log(ComponentID(c))
}

func TestComponentByID(t *testing.T) {
	c := &Hello{}
	ctx := NewWindow(Window{})
	ctx.Mount(c)
	defer ctx.Close()

	id := ComponentID(c)
	if c2 := ComponentByID(id); c != c2 {
		t.Error("c and c2 should be the same component")
	}
}
