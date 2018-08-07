package app

import (
	"sync"
	"testing"
)

func TestActions(t *testing.T) {
	HandleAction("test", func(e EventDispatcher, a Action) {})

	PostAction("test", 42)
	PostActions(
		Action{Name: "test", Arg: 21},
		Action{Name: "test", Arg: 84},
	)
}

func TestActionRegistry(t *testing.T) {
	var wg sync.WaitGroup

	d := newEventRegistry(func(f func()) {
		f()
	})
	r := newActionRegistry(d)

	r.Handle("test", func(e EventDispatcher, a Action) {
		wg.Done()
	})

	r.Post("unknown", nil)

	wg.Add(1)
	r.Post("test", 42)

	wg.Add(3)
	r.PostBatch(
		Action{Name: "test", Arg: nil},
		Action{Name: "test", Arg: "hello"},
		Action{Name: "test", Arg: 21},
	)

	wg.Wait()
}
