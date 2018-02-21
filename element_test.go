package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/tests"
)

func TestElemDB(t *testing.T) {
	tests.TestElemDB(t, func() app.ElementDB {
		return app.NewElemDB()
	})
}

func TestConcurrentElemDB(t *testing.T) {
	tests.TestElemDB(t, func() app.ElementDB {
		return app.NewConcurrentElemDB(app.NewElemDB())
	})
}
