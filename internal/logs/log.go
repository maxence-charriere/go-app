package logs

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Logger represents a function that formats using the default formats for its
// operands and logs the resulting string.
type Logger func(format string, a ...interface{})

// ToWriter returns a logger that writes on the given writer.
// It is safe for concurrent operations.
func ToWriter(w io.Writer) Logger {
	var mu sync.Mutex

	return func(format string, a ...interface{}) {
		mu.Lock()
		fmt.Fprintf(w, format, a...)
		fmt.Fprintln(w)
		mu.Unlock()
	}
}

// WithPrompt add a prompt the the logger.
func WithPrompt(l Logger) Logger {
	return func(format string, a ...interface{}) {
		format = prompt(false, a...) + format
		l(format, a...)
	}
}

// WithColoredPrompt add a prompt the the logger.
func WithColoredPrompt(l Logger) Logger {
	return func(format string, a ...interface{}) {
		format = prompt(true, a...) + format
		l(format, a...)
	}
}

func prompt(colors bool, a ...interface{}) string {
	defaultColor := ""
	logColor := ""
	errColor := ""

	if colors {
		defaultColor = "\033[00m"
		logColor = "\033[94m"
		errColor = "\033[91m"
	}

	level := "LOG"
	color := logColor

	for _, arg := range a {
		if _, ok := arg.(error); ok {
			level = "ERR"
			color = errColor
			break
		}
	}

	return fmt.Sprintf("%s%s%s %s %s‣%s ",
		color,
		level,
		defaultColor,
		time.Now().Format("2006/01/02 15:04:05"),
		color,
		defaultColor,
	)
}
