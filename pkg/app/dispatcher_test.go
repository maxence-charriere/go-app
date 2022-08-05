package app

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDispatcherMultipleMount(t *testing.T) {
	d := NewClientTester(Div())
	defer d.Close()
	d.Mount(A())
	d.Mount(Text("hello"))
	d.Mount(&hello{})
	d.Mount(&hello{})
	d.Consume()
}

func TestDispatcherAsyncWaitClient(t *testing.T) {
	d := NewClientTester(&hello{})
	defer d.Close()
	testDispatcherAsyncWait(t, d)
}

func TestDispatcherAsyncWaitServer(t *testing.T) {
	d := NewServerTester(&hello{})
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

func TestDispatcherLocalStorage(t *testing.T) {
	d := NewClientTester(&hello{})
	defer d.Close()
	testBrowserStorage(t, d.getLocalStorage())
}

func TestDispatcherSessionStorage(t *testing.T) {
	d := NewClientTester(&hello{})
	defer d.Close()
	testBrowserStorage(t, d.getSessionStorage())
}
