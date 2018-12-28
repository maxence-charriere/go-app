package app

import (
	"sync"
	"testing"
)

func TestHandle(t *testing.T) {
	wg := sync.WaitGroup{}

	Handle("test", func(m Msg) {
		t.Logf("key: %s value: %v", m.Key(), m.Value())
		wg.Done()
	})

	wg.Add(3)

	NewMsg("test").WithValue(42).Post()

	Post(
		NewMsg("test").WithValue(21),
		NewMsg("test").WithValue(84),
	)

	wg.Wait()
}

func TestMsgRegistry(t *testing.T) {
	r := newMsgRegistry()
	wg := sync.WaitGroup{}

	r.handle("test", func(m Msg) {
		wg.Done()
	})

	r.post(NewMsg("unknown"))

	wg.Add(1)
	r.post(NewMsg("test").WithValue(42))

	wg.Add(3)
	r.post(
		NewMsg("test"),
		NewMsg("test").WithValue("hello"),
		NewMsg("test").WithValue(21),
	)

	wg.Wait()
}
