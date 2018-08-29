package logs

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGoapp(t *testing.T) {
	addr := ":9000"
	b := &bytes.Buffer{}
	wg := sync.WaitGroup{}

	s := GoappServer{
		Addr:   addr,
		Writer: b,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		wg.Add(1)
		time.Sleep(time.Millisecond * 5)

		c := NewGoappClient(addr, WithColoredPrompt)
		defer c.Close()

		c.Logger()("hello")
		c.Logger()("my name is %s", "Maxence")
		c.Logger()("%s", errors.New("bye"))

		c.WaitForStop(func() {
			wg.Done()
		})
	}()

	go func() {
		time.Sleep(time.Millisecond * 10)
		cancel()
	}()

	err := s.ListenAndLog(ctx)
	assert.NoError(t, err)

	wg.Wait()
	t.Log(b.String())
}
