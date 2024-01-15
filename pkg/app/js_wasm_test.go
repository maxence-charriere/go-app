package app

import (
	"sync"
	"testing"
)

func TestPromise(t *testing.T) {
	callback := FuncOf(func(this Value, args []Value) any {
		args[0].Invoke("hi")
		return nil
	})
	defer callback.Release()

	var wg sync.WaitGroup
	wg.Add(1)

	promise := Window().Get("Promise").New(callback)
	promise.Then(func(v Value) {
		t.Log(v.String())
		wg.Done()
	})

	wg.Wait()
}
