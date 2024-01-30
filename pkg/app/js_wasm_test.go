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
	defer wg.Wait()
	wg.Add(1)

	promise := Window().Get("Promise").New(callback)
	promise.Then(func(v Value) {
		t.Log(v.String())
		wg.Done()
	})
}

func TestPromiseBuilder(t *testing.T) {
	var wg sync.WaitGroup
	defer wg.Wait()

	promiseBuilder := FuncOf(func(this Value, args []Value) any {
		callback := FuncOf(func(this Value, args []Value) any {
			wg.Add(1)
			args[0].Invoke("hi")
			wg.Done()
			return nil
		})

		return Window().Get("Promise").New(callback)
	})

	Window().Set("promiseBuilder", promiseBuilder)
	promiseBuilderCopy := Window().Get("promiseBuilder")

	wg.Add(1)
	promiseBuilderCopy.Invoke().Then(func(arg Value) {
		t.Log("then:", arg)
		wg.Done()
	})
}
