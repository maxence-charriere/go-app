package core_test

import (
	"testing"

	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/tests"
)

func TestHistory(t *testing.T) {
	tests.TestHistory(t, func() core.History {
		return core.NewHistory()
	})
}

func TestConcurrentHistory(t *testing.T) {
	tests.TestHistory(t, func() core.History {
		return core.ConcurrentHistory(core.NewHistory())
	})
}
