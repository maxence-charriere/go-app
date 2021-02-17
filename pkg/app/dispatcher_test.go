package app

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDispatcherMultipleMount(t *testing.T) {
	d := NewClientTestingDispatcher(Div())
	defer d.Close()
	d.Mount(A())
	d.Mount(Text("hello"))
	d.Mount(&hello{})
	d.Mount(&hello{})
	d.Consume()
}

func TestDispatcherAsyncWaitClient(t *testing.T) {
	d := NewClientTestingDispatcher(&hello{})
	defer d.Close()
	testDispatcherAsyncWait(t, d)
}

func TestDispatcherAsyncWaitServer(t *testing.T) {
	d := NewServerTestingDispatcher(&hello{})
	defer d.Close()
	testDispatcherAsyncWait(t, d)
}

func testDispatcherAsyncWait(t *testing.T, d Dispatcher) {
	var mu sync.Mutex
	var counts int

	inc := func() {
		mu.Lock()
		counts++
		mu.Unlock()
	}

	d.Async(inc)
	d.Async(inc)
	d.Async(inc)
	d.Async(inc)
	d.Async(inc)

	d.Wait()
	require.Equal(t, 5, counts)
}
