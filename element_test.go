package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/tests"
)

func TestElemDB(t *testing.T) {
	tests.TestElementDB(t, func() app.ElementDB {
		return app.NewElementDB()
	})
}

func TestConcurrentElemDB(t *testing.T) {
	tests.TestElementDB(t, func() app.ElementDB {
		return app.NewConcurrentElemDB(app.NewElementDB())
	})
}
