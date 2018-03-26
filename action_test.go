package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/tests"
)

func TestActions(t *testing.T) {
	app.HandleAction("test", func(e app.EventDispatcher, a app.Action) {})

	app.PostAction("test", 42)
	app.PostActionBatch(
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

func TestActionRegistryWithLogs(t *testing.T) {
	tests.TestActionRegistry(t, func() app.ActionRegistry {
		dispatcher := app.NewEventRegistry(func(f func()) {
			f()
		})

		r := app.NewActionRegistry(dispatcher)
		return app.ActionRegistryWithLogs(r)
	})
}
