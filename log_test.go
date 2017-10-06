package app

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	file, err := os.Create("logger-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("logger-test")
	defer file.Close()

	loggers := []Logger{
		NewLogger(file, true),
		NewConsole(true),
		NewMultiLogger(
			NewLogger(file, false),
			NewConsole(false),
		),
	}

	for _, logger := range loggers {
		logger.Log("log", "world")
		logger.Logf("%s %s", "logf", "world")

		logger.Error("error", "world")
		logger.Errorf("%s %s", "errorf", "world")
	}
}
