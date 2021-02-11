package app

import "testing"

func TestDispatcherMultipleMount(t *testing.T) {
	d := NewClientTestingDispatcher(Div())
	defer d.Close()
	d.Mount(A())
	d.Mount(Text("hello"))
	d.Mount(&hello{})
	d.Mount(&hello{})
	d.Consume()
}
