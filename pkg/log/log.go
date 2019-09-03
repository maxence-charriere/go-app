package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	// CurrentLevel describe the current log level. Only entries with a level
	// inferior or equal to the current level are printed.
	CurrentLevel = WarnLevel

	outpout    io.Writer = os.Stderr
	timeLayout           = "2006/01/02 15:04:05"
)

// Log colors
var (
	DefaultColor = "\033[00m"
	InfoColor    = "\033[94m"
	ErrorColor   = "\033[91m"
	WarnColor    = "\033[93m"
	DebugColor   = "\033[90m"
)

// Constants that describe log level.
const (
	InfoLevel Level = iota
	ErrorLevel
	WarnLevel
	DebugLevel
)

// Level represents a log level.
type Level int

func (l Level) String() string {
	switch l {
	case InfoLevel:
		return "INFO"
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARN"
	default:
		return "DEBUG"
	}
}

func (l Level) color() string {
	switch l {
	case InfoLevel:
		return InfoColor
	case ErrorLevel:
		return ErrorColor
	case WarnLevel:
		return WarnColor
	default:
		return DebugColor
	}
}

// Info logs an info message.
func Info(v ...interface{}) Entry {
	return Log(Entry{
		Level:   InfoLevel,
		Message: sprint(v...),
	})
}

// Infof logs an info message according to the given format.
func Infof(format string, v ...interface{}) Entry {
	return Log(Entry{
		Level:   InfoLevel,
		Message: fmt.Sprintf(format, v...),
	})
}

// Error logs an error message.
func Error(v ...interface{}) Entry {
	return Log(Entry{
		Level:   ErrorLevel,
		Message: sprint(v...),
	})
}

// Errorf logs an error message according to the given format.
func Errorf(format string, v ...interface{}) Entry {
	return Log(Entry{
		Level:   ErrorLevel,
		Message: fmt.Sprintf(format, v...),
	})
}

// Warn logs an warning message.
func Warn(v ...interface{}) Entry {
	return Log(Entry{
		Level:   WarnLevel,
		Message: sprint(v...),
	})
}

// Warnf logs an warning message according to the given format.
func Warnf(format string, v ...interface{}) Entry {
	return Log(Entry{
		Level:   WarnLevel,
		Message: fmt.Sprintf(format, v...),
	})
}

// Debug logs an debug message.
func Debug(v ...interface{}) Entry {
	return Log(Entry{
		Level:   DebugLevel,
		Message: sprint(v...),
	})
}

// Debugf logs an debug message according to the given format.
func Debugf(format string, v ...interface{}) Entry {
	return Log(Entry{
		Level:   DebugLevel,
		Message: fmt.Sprintf(format, v...),
	})
}

// Log logs the given entry.
func Log(e Entry) Entry {
	e.printMessage()
	for k, v := range e.Tags {
		e.printTag(k, v)
	}
	return e
}

// Entry represents a log entry.
type Entry struct {
	Level   Level
	Message string
	Tags    map[string]string
}

// T appends and logs the tag described by the given key/value.
func (e Entry) T(key string, value interface{}) Entry {
	if e.Tags == nil {
		e.Tags = make(map[string]string)
	}

	v := fmt.Sprintf("%+v", value)
	e.Tags[key] = v
	e.printTag(key, v)
	return e
}

// Panic call panic with the entry message.
func (e Entry) Panic() {
	panic(e.Message)
}

func (e Entry) printMessage() {
	if e.Level > CurrentLevel {
		return
	}

	fmt.Fprintf(outpout, "%s%s%s %s %s‣%s %s\n",
		e.Level.color(),
		e.Level,
		DefaultColor,
		time.Now().Format(timeLayout),
		e.Level.color(),
		DefaultColor,
		e.Message,
	)
}

func (e Entry) printTag(k, v string) {
	if e.Level > CurrentLevel {
		return
	}

	fmt.Fprintf(outpout, "    %s: %s\n", k, v)
}

func sprint(v ...interface{}) string {
	return strings.TrimSuffix(fmt.Sprintln(v...), "\n")
}
