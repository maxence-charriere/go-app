package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
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

func TestLoadErrorOutput(t *testing.T) {
	for _, tt := range []struct {
		name       string
		args       []string
		env        string
		noCommand  bool
		nilOptions bool
		message    string
	}{
		{name: "invalid flag value", args: []string{"-count", "invalid"}, message: "invalid value"},
		{name: "unknown flag", args: []string{"-unknown"}, message: "flag provided but not defined"},
		{name: "unknown command", args: []string{"unknown"}, message: "unknown command"},
		{name: "no root command", noCommand: true, message: "no root command"},
		{name: "nil options", nilOptions: true, message: "receiver"},
		{name: "invalid environment", env: "private-config-value", message: "GOAPP_TEST_COUNT"},
		{name: "flag cannot override invalid environment", args: []string{"-count", "42"}, env: "private-config-value", message: "GOAPP_TEST_COUNT"},
		{name: "help", args: []string{"-h"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			oldManager, oldUsage, oldArgs, oldExit := defaultManager, currentUsage, programArgs, exitOnError
			t.Cleanup(func() {
				defaultManager, currentUsage, programArgs, exitOnError = oldManager, oldUsage, oldArgs, oldExit
			})
			var out bytes.Buffer
			defaultManager = commandManager{out: &out}
			programArgs = tt.args
			exitOnError = false
			type options struct {
				Count int `env:"GOAPP_TEST_COUNT"`
			}
			opts := &options{}
			if tt.nilOptions {
				opts = nil
			}
			if !tt.noCommand {
				Register().Options(opts)
			}
			env := "21"
			if tt.env != "" {
				env = tt.env
			}
			t.Setenv("GOAPP_TEST_COUNT", env)
			require.Panics(t, func() { Load() })
			if tt.message == "" {
				require.NotContains(t, out.String(), "Error:")
				require.Contains(t, out.String(), "Usage:")
				return
			}
			require.Contains(t, out.String(), "Error:")
			require.Contains(t, out.String(), tt.message)
			require.NotContains(t, out.String(), "private-config-value")
		})
	}
}

func TestError(t *testing.T) {
	w := bytes.NewBufferString("\n")
	defaultManager.out = w
	Error(errors.New("error error critical error"))
	t.Log(w.String())
}

func TestContextWithSignals(t *testing.T) {
	t.Run("explicit cancellation", func(t *testing.T) {
		parent, parentCancel := context.WithCancel(context.Background())
		defer parentCancel()
		ctx, cancel := ContextWithSignals(parent, os.Interrupt)
		defer cancel()
		require.NoError(t, ctx.Err())
		cancel()
		cancel() // Cleanup is safe to repeat.
		waitForSignalContext(t, ctx)
		require.NoError(t, parent.Err())
	})

	for _, alreadyCanceled := range []bool{false, true} {
		name := "parent canceled later"
		if alreadyCanceled {
			name = "parent already canceled"
		}
		t.Run(name, func(t *testing.T) {
			parent, parentCancel := context.WithCancel(context.Background())
			defer parentCancel()
			if alreadyCanceled {
				parentCancel()
			}
			ctx, cancel := ContextWithSignals(parent, os.Interrupt)
			defer cancel()
			parentCancel()
			waitForSignalContext(t, ctx)
		})
	}

	for _, scenario := range []string{
		"signals",
		"all signals",
		"stop",
		"stop after signal",
	} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			executable, err := os.Executable()
			require.NoError(t, err)
			cmd := exec.CommandContext(ctx, executable, "-test.run=^TestContextWithSignalsSubprocess$")
			cmd.Env = append(os.Environ(), "GOAPP_TEST_SIGNAL_CONTEXT="+scenario)
			output, err := cmd.CombinedOutput()
			require.NoError(t, ctx.Err(), "subprocess timed out: %s", output)
			if scenario == "stop" || scenario == "stop after signal" {
				var exitError *exec.ExitError
				require.True(t, errors.As(err, &exitError), "expected signal termination, got %v: %s", err, output)
				status, ok := exitError.Sys().(syscall.WaitStatus)
				require.True(t, ok)
				require.True(t, status.Signaled(), "subprocess did not terminate by signal: %s", output)
				require.Equal(t, syscall.SIGINT, status.Signal())
				return
			}
			require.NoError(t, err, "%s", output)
		})
	}
}

func TestContextWithSignalsSubprocess(t *testing.T) {
	scenario := os.Getenv("GOAPP_TEST_SIGNAL_CONTEXT")
	if scenario == "" {
		return
	}
	// Ensure the subprocess starts with the default interrupt behavior.
	signal.Reset(os.Interrupt)
	signals := []os.Signal{os.Interrupt}
	if scenario == "all signals" {
		signals = nil
	}
	ctx, cancel := ContextWithSignals(context.Background(), signals...)
	defer cancel()

	if scenario == "stop" || scenario == "stop after signal" {
		if scenario == "stop after signal" {
			require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGINT))
			waitForSignalContext(t, ctx)
		}
		cancel()
		require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGINT))
		time.Sleep(time.Second)
		t.Fatal("interrupt did not terminate the process after cleanup")
	}

	// Observe delivery so the second signal is not coalesced with the first.
	delivered := make(chan os.Signal, 1)
	signal.Notify(delivered, os.Interrupt)
	defer signal.Stop(delivered)
	for range 2 {
		require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGINT))
		select {
		case <-delivered:
		case <-time.After(5 * time.Second):
			t.Fatal("signal was not delivered")
		}
		waitForSignalContext(t, ctx)
	}
	// Wait for in-flight signal delivery before the subprocess exits.
	signal.Stop(delivered)
}

func TestContextWithSignalsGoroutineCleanup(t *testing.T) {
	// Initialize os/signal's shared delivery goroutine before measuring.
	delivered := make(chan os.Signal, 1)
	signal.Notify(delivered, os.Interrupt)
	defer signal.Stop(delivered)

	for _, scenario := range []string{
		"explicit cancellation",
		"parent canceled later",
		"parent already canceled",
	} {
		t.Run(scenario, func(t *testing.T) {
			baseline := runtime.NumGoroutine()
			for range 100 {
				parent, parentCancel := context.WithCancel(context.Background())
				if scenario == "parent already canceled" {
					parentCancel()
				}
				ctx, cancel := ContextWithSignals(parent, os.Interrupt)
				t.Cleanup(cancel)
				if scenario == "explicit cancellation" {
					cancel()
				}
				parentCancel()
				waitForSignalContext(t, ctx)
			}
			// Allow small runtime fluctuations, but detect retained workers.
			require.Eventually(t, func() bool {
				return runtime.NumGoroutine() <= baseline+2
			}, 5*time.Second, 10*time.Millisecond)
		})
	}
}

func waitForSignalContext(t *testing.T, ctx context.Context) {
	t.Helper()
	select {
	case <-ctx.Done():
		require.Equal(t, context.Canceled, ctx.Err())
	case <-time.After(5 * time.Second):
		t.Fatal("context was not canceled")
	}
}
