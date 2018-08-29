package logs

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGoapp(t *testing.T) {
	addr := ":9000"
	b := &bytes.Buffer{}

	s := GoappServer{
		Addr:   addr,
		Writer: b,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientStopped := false

	go func() {
		time.Sleep(time.Millisecond * 5)

		c := NewGoappClient(addr, WithColoredPrompt)
		defer c.Close()

		c.Logger()("hello")
		c.Logger()("my name is %s", "Maxence")
		c.Logger()("no: %s", errors.New("i don't want to go"))

		c.WaitForStop(func() {
			c.Logger()("bye")
			clientStopped = true
		})
	}()

	go func() {
		time.Sleep(time.Millisecond * 100)
		cancel()
	}()

	err := s.ListenAndLog(ctx)
	assert.NoError(t, err)
	assert.True(t, clientStopped)
	t.Log(b.String())

}
