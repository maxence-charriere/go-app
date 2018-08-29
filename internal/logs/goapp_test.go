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

type concurrentBuffer struct {
	buffer bytes.Buffer
	mutex  sync.Mutex
}

func (b *concurrentBuffer) Write(d []byte) (int, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.Write(d)
}

func (b *concurrentBuffer) Read(d []byte) (int, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.Read(d)
}

func (b *concurrentBuffer) String() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.String()
}

func TestGoapp(t *testing.T) {
	addr := ":9000"
	b := &concurrentBuffer{}

	s := GoappServer{
		Addr:   addr,
		Writer: b,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopc := make(chan struct{})
	defer close(stopc)

	go func() {
		time.Sleep(time.Millisecond * 5)

		c := NewGoappClient(addr, WithColoredPrompt)
		defer c.Close()

		c.Logger()("hello")
		c.Logger()("my name is %s", "Maxence")
		c.Logger()("%s", errors.New("bye"))

		c.WaitForStop(func() {
			stopc <- struct{}{}
		})
	}()

	go func() {
		time.Sleep(time.Millisecond * 100)
		cancel()
	}()

	err := s.ListenAndLog(ctx)
	assert.NoError(t, err)

	<-stopc
	t.Log(b.String())

}
