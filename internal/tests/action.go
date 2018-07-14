package tests

import (
	"sync"
	"testing"

	"github.com/murlokswarm/app"
)

// TestActionRegistry is a test suite that ensure all the action registries
// behave the same.
func TestActionRegistry(t *testing.T, newRegistry func() app.ActionRegistry) {
	var wg sync.WaitGroup
	r := newRegistry()

	r.Handle("test", func(e app.EventDispatcher, a app.Action) {
		wg.Done()
	})

	r.Post("unknown", nil)

	wg.Add(1)
	r.Post("test", 42)

	wg.Add(3)
	r.PostBatch(
		app.Action{Name: "test", Arg: nil},
		app.Action{Name: "test", Arg: "hello"},
		app.Action{Name: "test", Arg: 21},
	)

	wg.Wait()
}
