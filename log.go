package app

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// A Logger represents an active logging object that generates lines of output.
type Logger interface {
	// Log logs a message using the default formats for its operands.
	// Spaces are always added between operands and a newline is appended.
	Log(v ...interface{})

	// Logf logs a message according to a format specifier.
	Logf(format string, v ...interface{})

	// Error logs an error using the default formats for its operands.
	// Spaces are always added between operands and a newline is appended.
	Error(v ...interface{})

	// Errorf logs an error according to a format specifier.
	Errorf(format string, v ...interface{})
}

const (
	defaultColor string = "\033[00m"
	accentColor  string = "\033[94m"
	errColor     string = "\033[91m"
)

// NewLogger creates a logger that writes on the given writer.
// Logs are written only if debug is enabled.
func NewLogger(w io.Writer, debug bool) Logger {
	return &logger{
		writer: w,
		debug:  debug,
	}
}

type logger struct {
	writer io.Writer
	debug  bool
}

func (l *logger) Log(v ...interface{}) {
	if !l.debug {
		return
	}
	printLogPrefix(l.writer, "Log  ", accentColor)
	fmt.Fprintln(l.writer, v...)
}

func (l *logger) Logf(format string, v ...interface{}) {
	if !l.debug {
		return
	}
	printLogPrefix(l.writer, "Log  ", accentColor)
	fmt.Fprintf(l.writer, format, v...)
	fmt.Fprintln(l.writer)
}

func (l *logger) Error(v ...interface{}) {
	printLogPrefix(l.writer, "Error", errColor)
	fmt.Fprintln(l.writer, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	printLogPrefix(l.writer, "Error", errColor)
	fmt.Fprintf(l.writer, format, v...)
	fmt.Fprintln(l.writer)
}

func printLogPrefix(w io.Writer, level, color string) {
	now := time.Now().Format("2006/01/02 15:04:05")
	fmt.Fprintf(w,
		"%s%s%s %s %s|>%s ",
		color,
		strings.ToUpper(level),
		defaultColor,
		now,
		color,
		defaultColor,
	)
}

// NewConsole creates a logger that writes messages on standard outputs.
// Logs are written on stdout, only if debug is enabled.
// Errors are written on stderr.
// It is safe for concurrent access.
func NewConsole(debug bool) Logger {
	logger := newConsole(debug)
	return NewConcurrentLogger(logger)
}

type console struct {
	std Logger
	err Logger
}

func newConsole(debug bool) *console {
	return &console{
		std: NewLogger(os.Stdout, debug),
		err: NewLogger(os.Stderr, debug),
	}
}

func (c *console) Log(v ...interface{}) {
	c.std.Log(v...)
}

func (c *console) Logf(format string, v ...interface{}) {
	c.std.Logf(format, v...)
}

func (c *console) Error(v ...interface{}) {
	c.err.Error(v...)
}

func (c *console) Errorf(format string, v ...interface{}) {
	c.err.Errorf(format, v...)
}

// NewMultiLogger creates a logger that aggregate multiple loggers.
func NewMultiLogger(loggers ...Logger) Logger {
	return &multiLogger{
		loggers: loggers,
	}
}

type multiLogger struct {
	loggers []Logger
}

func (l *multiLogger) Log(v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Log(v...)
	}
}

func (l *multiLogger) Logf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Logf(format, v...)
	}
}

func (l *multiLogger) Error(v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Error(v...)
	}
}

func (l *multiLogger) Errorf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Errorf(format, v...)
	}
}

// NewConcurrentLogger decorates the given logger to ensure concurrent access
// safety.
func NewConcurrentLogger(l Logger) Logger {
	return &concurrentLogger{
		logger: l,
	}
}

type concurrentLogger struct {
	mutex  sync.Mutex
	logger Logger
}

func (l *concurrentLogger) Log(v ...interface{}) {
	l.mutex.Lock()
	l.logger.Log(v...)
	l.mutex.Unlock()
}

func (l *concurrentLogger) Logf(format string, v ...interface{}) {
	l.mutex.Lock()
	l.logger.Logf(format, v...)
	l.mutex.Unlock()
}

func (l *concurrentLogger) Error(v ...interface{}) {
	l.mutex.Lock()
	l.logger.Error(v...)
	l.mutex.Unlock()
}

func (l *concurrentLogger) Errorf(format string, v ...interface{}) {
	l.mutex.Lock()
	l.logger.Errorf(format, v...)
	l.mutex.Unlock()
}
