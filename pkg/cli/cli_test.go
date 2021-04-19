package cli

import (
	"bytes"
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
	"github.com/stretchr/testify/require"
)

func init() {
	exitOnError = false
}

func TestCliSuccess(t *testing.T) {
	w := bytes.NewBufferString("\n")
	defaultManager.out = w
	defaultManager.commands = nil
	programArgs = nil

	Register().Help("A test command")
	cmd := Load()
	require.Empty(t, cmd)

	Usage()
	t.Log(w.String())
}

func TestCliIndex(t *testing.T) {
	w := bytes.NewBufferString("\n")
	defaultManager.out = w
	defaultManager.commands = nil
	programArgs = []string{"foo", "test"}

	Register("foo", "bar").Help("A test command")
	Register("foo", "buu").Help("Another test command")

	defer func() {
		recover()
		t.Log(w.String())
	}()

	Load()
	t.Fail()
}

func TestCliCmdBadOption(t *testing.T) {
	w := bytes.NewBufferString("\n")
	defaultManager.out = w
	defaultManager.commands = nil
	programArgs = []string{"-duration", "[x_x]"}

	opts := struct {
		Duration time.Duration
	}{}

	Register().Options(&opts)

	defer func() {
		recover()
		t.Log(w.String())
	}()

	Load()
	t.Fail()
}

func TestUsagePanic(t *testing.T) {
	currentUsage = nil
	require.Panics(t, func() {
		Usage()
	})
}

func TestError(t *testing.T) {
	w := bytes.NewBufferString("\n")
	defaultManager.out = w
	Error(errors.New("error error critical error"))
	t.Log(w.String())
}

func TestContextWithSignals(t *testing.T) {
	ctx, cancel := ContextWithSignals(context.TODO(), os.Interrupt)
	defer cancel()

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	<-ctx.Done()
	require.Error(t, ctx.Err())
}
