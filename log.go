package app

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// A Logger represents an active logging object that generates lines of output.
type Logger interface {
	// Log logs a message using the default formats for its operands.
	// Spaces are always added between operands and a newline is appended.
	Log(v ...interface{})

	// Logf logs a message according to a format specifier.
	Logf(format string, v ...interface{})

	// Log logs an error using the default formats for its operands.
	// Spaces are always added between operands and a newline is appended.
	Error(v ...interface{})

	// Logf logs an error according to a format specifier.
	Errorf(format string, v ...interface{})
}

const (
	defaultColor string = "\033[00m"
	accentColor  string = "\033[94m"
	errColor     string = "\033[91m"
)

// Console is a logger that writes messages on the standard output.
// It satisfies the Logger interface.
type Console struct {
	Debug bool
}

// Log satisfies the Logger interface.
func (c *Console) Log(v ...interface{}) {
	if !c.Debug {
		return
	}
	printLogPrefix("Log  ", accentColor)
	fmt.Println(v...)
}

// Logf satisfies the Logger interface.
func (c *Console) Logf(format string, v ...interface{}) {
	if !c.Debug {
		return
	}
	printLogPrefix("Log  ", accentColor)
	fmt.Printf(format, v...)
	fmt.Println()
}

// Error satisfies the Logger interface.
func (c *Console) Error(v ...interface{}) {
	printLogPrefix("Error", errColor)
	fmt.Println(v...)
}

// Errorf satisfies the Logger interface.
func (c *Console) Errorf(format string, v ...interface{}) {
	printLogPrefix("Error", errColor)
	fmt.Printf(format, v...)
	fmt.Println()
}

func printLogPrefix(level, color string) {
	file, line := caller()
	now := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("%s%s%s %s %s:%v |> ", color, strings.ToUpper(level), defaultColor, now, file, line)
}

func caller() (file string, line int) {
	_, file, line, _ = runtime.Caller(3)
	file = filepath.Base(file)
	return
}

// MultiLogger is a logger that logs messages using multiple loggers.
// It satisfies the Logger interface.
type MultiLogger struct {
	loggers []Logger
}

// NewMultiLogger creates a logger that uses loggers.
func NewMultiLogger(loggers ...Logger) *MultiLogger {
	return &MultiLogger{
		loggers: loggers,
	}
}

// Log satisfies the Logger interface.
func (l *MultiLogger) Log(v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Log(v...)
	}
}

// Logf satisfies the Logger interface.
func (l *MultiLogger) Logf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Logf(format, v...)
	}
}

// Error satisfies the Logger interface.
func (l *MultiLogger) Error(v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Error(v...)
	}
}

// Errorf satisfies the Logger interface.
func (l *MultiLogger) Errorf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Errorf(format, v...)
	}
}
