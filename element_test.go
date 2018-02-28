package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/tests"
)

func TestElemDB(t *testing.T) {
	tests.TestElemDB(t, func() app.ElemDB {
		return app.NewElemDB()
	})
}

func TestConcurrentElemDB(t *testing.T) {
	tests.TestElemDB(t, func() app.ElemDB {
		return app.NewConcurrentElemDB(app.NewElemDB())
	})
}
