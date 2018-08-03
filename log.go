package app

import (
	"fmt"
	"io"
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

// NewLogger creates a logger that writes on the given writers.
// Logs that contain errors are logged on werr.
func NewLogger(wout, werr io.Writer, debug, colors bool) Logger {
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
		colors:    colors,
	}
}

type logger struct {
	wout      io.Writer
	werr      io.Writer
	whenDebug func(func())
	colors    bool
}

func (l *logger) Log(format string, v ...interface{}) {
	for _, i := range v {
		if _, ok := i.(error); ok {
			l.print(levelError, format, v...)
			return
		}
	}
	l.print(levelInfo, format, v...)
}

func (l *logger) Debug(format string, v ...interface{}) {
	l.print(levelDebug, format, v...)
}

func (l *logger) WhenDebug(f func()) {
	l.whenDebug(f)
}

func (l *logger) print(level int, format string, v ...interface{}) {
	format = l.prefix(level) + format

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
	logLevel := "INFO "
	color := infoColor
	endColor := defaultColor

	switch level {
	case levelError:
		logLevel = "ERROR"
		color = errColor

	case levelDebug:
		logLevel = "DEBUG"
		color = debugColor
	}

	if !l.colors {
		color = ""
		endColor = ""
	}

	return fmt.Sprintf("%s%s%s %s %s|>%s ",
		color,
		logLevel,
		endColor,
		time.Now().Format("2006/01/02 15:04:05"),
		color,
		endColor,
	)
}

const (
	levelInfo = iota
	levelError
	levelDebug

	defaultColor string = "\033[00m"
	infoColor    string = "\033[94m"
	errColor     string = "\033[91m"
	debugColor   string = "\033[95m"
)
