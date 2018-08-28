package goapp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := NewLogger(true, true, func() {
		t.Log("stop")
	})

	go func() {
		l.Log("Hello")
		l.Debug("world")
		time.Sleep(time.Millisecond * 10)
		cancel()
	}()

	err := ListenAndWriteLogs(ctx)
	assert.NoError(t, err)
}
