package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestEventRegistry(t *testing.T) {
	tests.TestEventRegistry(t, func() app.EventRegistry {
		return app.NewEventRegistry(func(f func()) {
			f()
		})
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

func TestEventSubscriber(t *testing.T) {
	s := app.NewEventSubscriber()
	defer s.Close()

	s.Subscribe("test-event-subscriber", func() {})
}
