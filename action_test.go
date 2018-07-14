package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestActions(t *testing.T) {
	app.Handle("test", func(e app.EventDispatcher, a app.Action) {})

	app.NewAction("test", 42)
	app.NewActions(
		app.Action{Name: "test", Arg: 21},
		app.Action{Name: "test", Arg: 84},
	)
}

func TestActionRegistry(t *testing.T) {
	tests.TestActionRegistry(t, func() app.ActionRegistry {
		dispatcher := app.NewEventRegistry(func(f func()) {
			f()
		})

		return app.NewActionRegistry(dispatcher)
	})
}
