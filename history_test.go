package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/tests"
)

func TestHistory(t *testing.T) {
	tests.TestHistory(t, func() app.History {
		return app.NewHistory()
	})
}

func TestConcurrentHistory(t *testing.T) {
	tests.TestHistory(t, func() app.History {
		return app.ConcurrentHistory(app.NewHistory())
	})
}
