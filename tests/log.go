package tests

import (
	"testing"

	"github.com/murlokswarm/app"
)

// TestLogger is a test used to ensure that all logger implementations behave
// the same.
func TestLogger(t *testing.T, logger app.Logger) {
	logger.Log("log", "world")
	logger.Logf("%s %s", "logf", "world")

	logger.Error("error", "world")
	logger.Errorf("%s %s", "errorf", "world")
}
