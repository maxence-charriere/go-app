package app

import (
	"sync"
	"syscall/js"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
		defer callback.Release()

		return Window().Get("Promise").New(callback)
	})
	defer promiseBuilder.Release()

	Window().Set("promiseBuilder", promiseBuilder)
	defer Window().Delete("promiseBuilder")
	promiseBuilderCopy := Window().Get("promiseBuilder")

	wg.Add(1)
	promiseBuilderCopy.Invoke().Then(func(arg Value) {
		t.Log("then:", arg)
		wg.Done()
	})
}

// Capture the chain returned by then so tests can observe rejection propagation.
type trackedPromise struct {
	js.Value
	chain js.Value
}

func (p *trackedPromise) Call(method string, args ...any) js.Value {
	p.chain = p.Value.Call(method, args...)
	return p.chain
}

func awaitTestPromise(t *testing.T, promise js.Value) (Value, bool) {
	t.Helper()
	type result struct {
		value    Value
		rejected bool
	}
	done := make(chan result, 1)
	fulfilled := FuncOf(func(_ Value, args []Value) any {
		done <- result{value: args[0]}
		return nil
	})
	defer fulfilled.Release()
	rejected := FuncOf(func(_ Value, args []Value) any {
		done <- result{value: args[0], rejected: true}
		return nil
	})
	defer rejected.Release()
	valueOf(promise).Call("then", fulfilled, rejected)
	select {
	case r := <-done:
		return r.value, r.rejected
	case <-time.After(5 * time.Second):
		t.Fatal("promise did not settle")
		return nil, false
	}
}

func TestPromiseSettlement(t *testing.T) {
	for _, tt := range []struct {
		name     string
		method   string
		argument any
		rejected bool
	}{
		{name: "fulfillment", method: "resolve", argument: "hello"},
		{name: "undefined fulfillment", method: "resolve", argument: js.Undefined()},
		{name: "rejection", method: "reject", argument: js.Global().Get("Error").New("failed"), rejected: true},
		{name: "undefined rejection", method: "reject", argument: js.Undefined(), rejected: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			promise := &trackedPromise{Value: js.Global().Get("Promise").Call(tt.method, tt.argument)}
			calls := 0
			value{jsValue: promise}.Then(func(v Value) {
				calls++
				require.True(t, v.Equal(valueOf(tt.argument)))
			})
			result, rejected := awaitTestPromise(t, promise.chain)
			require.Equal(t, tt.rejected, rejected)
			if tt.rejected {
				require.Zero(t, calls)
				require.True(t, result.Equal(valueOf(tt.argument)), "rejection reason changed")
			} else {
				require.Equal(t, 1, calls)
			}
		})
	}
}

func TestPromisePending(t *testing.T) {
	var resolve Value
	executor := FuncOf(func(_ Value, args []Value) any {
		resolve = args[0]
		return nil
	})
	defer executor.Release()
	promise := &trackedPromise{Value: js.Global().Get("Promise").New(syscalJSValueOf(executor))}
	calls := 0
	value{jsValue: promise}.Then(func(v Value) { calls++ })
	require.Zero(t, calls)
	resolve.Invoke(42)
	_, rejected := awaitTestPromise(t, promise.chain)
	require.False(t, rejected)
	require.Equal(t, 1, calls)
}

func TestPromiseUnhandledRejection(t *testing.T) {
	reason := Window().Get("Error").New("unhandled promise test")
	done := make(chan Value, 1)
	listener := FuncOf(func(_ Value, args []Value) any {
		event := args[0]
		if event.Get("reason").Equal(reason) {
			event.Call("preventDefault")
			done <- event.Get("promise")
		}
		return nil
	})
	Window().Call("addEventListener", "unhandledrejection", listener)
	defer func() {
		Window().Call("removeEventListener", "unhandledrejection", listener)
		listener.Release()
	}()

	promise := &trackedPromise{Value: js.Global().Get("Promise").Call("reject", syscalJSValueOf(reason))}
	value{jsValue: promise}.Then(func(Value) {
		t.Error("rejection called fulfillment handler")
	})
	select {
	case reported := <-done:
		require.True(t, reported.Equal(valueOf(promise.chain)))
	case <-time.After(5 * time.Second):
		t.Fatal("browser did not report the unhandled rejection")
	}
}
