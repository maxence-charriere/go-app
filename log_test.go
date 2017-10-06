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
		logger.Log("hello", "world")
		logger.Logf("%s %s", "hello", "world")

		logger.Error("hello", "world")
		logger.Errorf("%s %s", "hello", "world")
	}
}
