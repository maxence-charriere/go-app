package app

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Logger is the interface that describes an active logging object that
// generates lines of output.
type Logger interface {
	// Log logs a message according to a format specifier.
	Log(format string, args ...interface{})

	// Debug logs a debug message according to a format specifier.
	Debug(format string, args ...interface{})

	// WhenDebug execute the given function when debug mode is enabled.
	WhenDebug(f func())
}

var (
	// Loggers is the loggers used by the app.
	Loggers = []Logger{
		NewLogger(os.Stdout, os.Stderr, false),
	}
)

// Log logs a message according to a format specifier.
// It is a helper function that calls Log() for all the loggers set in
// app.Loggers.
func Log(format string, v ...interface{}) {
	for _, l := range Loggers {
		l.Log(format, v...)
	}
}

// Debug logs a debug message according to a format specifier.
// It is a helper function that calls Debug() for all the loggers set in
// app.Loggers.
func Debug(format string, v ...interface{}) {
	for _, l := range Loggers {
		l.Debug(format, v...)
	}
}

// WhenDebug execute the given function when debug mode is enabled.
// It is a helper function that calls WhenDebug() for all the loggers set in
// app.Loggers.
func WhenDebug(f func()) {
	for _, l := range Loggers {
		l.WhenDebug(f)
	}
}

// NewLogger creates a logger that writes on the given writers.
// Logs that contain errors are logged on werr.
func NewLogger(wout, werr io.Writer, debug bool) Logger {
	whenDebug := func(f func()) {}

	if debug {
		whenDebug = func(f func()) {
			f()
		}
	}

	return &logger{
		wout:      wout,
		werr:      wout,
		whenDebug: whenDebug,
	}
}

type logger struct {
	wout      io.Writer
	werr      io.Writer
	whenDebug func(func())
	indent    string
}

func (l *logger) Log(format string, v ...interface{}) {
	for _, i := range v {
		if _, ok := i.(error); ok {
			l.print(levelError, format, v...)
			return
		}
	}
	l.print(levelLog, format, v...)
}

func (l *logger) Debug(format string, v ...interface{}) {
	l.print(levelDebug, format, v...)
}

func (l *logger) WhenDebug(f func()) {
	l.whenDebug(f)
}

func (l *logger) print(level int, format string, v ...interface{}) {
	prefix := l.prefix(level)

	if len(l.indent) == 0 {
		l.indent = l.genIndent(len(prefix) - len(defaultColor)*4)
		fmt.Println("--", l.indent, "--")
	}

	format = prefix + format
	format = strings.Replace(format, "\n", "\n"+l.indent, -1)

	if format[len(format)-1] != '\n' {
		format += "\n"
	}

	if level == levelError {
		fmt.Fprintf(l.werr, format, v...)
		return
	}
	fmt.Fprintf(l.wout, format, v...)
}

func (l *logger) prefix(level int) string {
	logLevel := "LOG  "
	color := logColor

	switch level {
	case levelError:
		logLevel = "ERROR"
		color = errColor

	case levelDebug:
		logLevel = "DEBUG"
		color = debugColor
	}

	return fmt.Sprintf("%s%s%s %s %s|>%s ",
		color,
		logLevel,
		defaultColor,
		time.Now().Format("2006/01/02 15:04:05"),
		color,
		defaultColor,
	)
}

func (l *logger) genIndent(ilen int) string {
	indent := ""
	for i := 0; i < ilen; i++ {
		indent += " "
	}
	return indent
}

const (
	levelLog = iota
	levelError
	levelDebug

	defaultColor string = "\033[00m"
	logColor     string = "\033[94m"
	errColor     string = "\033[91m"
	debugColor   string = "\033[95m"
)
