package app

import (
	"testing"

	"github.com/maxence-charriere/go-app/v8/pkg/errors"
	"github.com/maxence-charriere/go-app/v8/pkg/logs"
)

func TestLog(t *testing.T) {
	DefaultLogger = t.Logf
	Log("hello", "world")
	Logf("hello %v", "Maxoo")
}

func TestServerLog(t *testing.T) {
	testSkipWasm(t)
	testLogger(t, serverLog)
}

func TestClientLog(t *testing.T) {
	testSkipNonWasm(t)
	testLogger(t, clientLog)
}

func testLogger(t *testing.T, l func(string, ...interface{})) {
	utests := []struct {
		scenario string
		value    interface{}
	}{
		{
			scenario: "log",
			value:    logs.New("test").Tag("type", "log"),
		},
		{
			scenario: "error",
			value:    errors.New("test").Tag("type", "error"),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			l("%v", u.value)
		})
	}
}
