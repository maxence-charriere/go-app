package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/tests"
)

func TestEventRegistry(t *testing.T) {
	tests.TestEventRegistry(t, func() app.EventRegistry {
		return app.NewEventRegistry(func(f func()) {
			f()
		})
	})
}

func TestEventRegistryWithLogs(t *testing.T) {
	tests.TestEventRegistry(t, func() app.EventRegistry {
		r := app.NewEventRegistry(func(f func()) {
			f()
		})
		return app.EventRegistryWithLogs(r)
	})
}

func TestConcurrentEventRegistry(t *testing.T) {
	tests.TestEventRegistry(t, func() app.EventRegistry {
		r := app.NewEventRegistry(func(f func()) {
			f()
		})
		return app.ConcurrentEventRegistry(r)
	})
}
