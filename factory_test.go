package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/tests"
)

func TestFactory(t *testing.T) {
	tests.TestFactory(t, func() app.Factory {
		return app.NewFactory()
	})
}

func TestFactoryWithLogs(t *testing.T) {
	tests.TestFactory(t, func() app.Factory {
		factory := app.NewFactory()
		factory = app.NewFactoryWithLogs(factory)
		return factory
	})
}

func TestConcurrentFactory(t *testing.T) {
	tests.TestFactory(t, func() app.Factory {
		factory := app.NewFactory()
		factory = app.NewConcurrentFactory(factory)
		return factory
	})
}
