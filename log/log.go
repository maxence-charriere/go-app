package log

import "runtime"
import "path/filepath"
import "fmt"
import "strings"

const (
	defaultColor string = "\033[00m"
	accentColor  string = "\033[94m"
	errColor     string = "\033[91m"
)

// Logger is a logger that satisfies the app.Logger interface.
type Logger struct {
	Debug bool
}

// Log satisfies the app.Logger interface.
func (l *Logger) Log(v ...interface{}) {
	if !l.Debug {
		return
	}
	printPrefix("Log", accentColor)
	fmt.Println(v...)
}

// Logf satisfies the app.Logger interface.
func (l *Logger) Logf(format string, v ...interface{}) {
	if !l.Debug {
		return
	}
	printPrefix("Log", accentColor)
	fmt.Printf(format, v...)
}

// Error satisfies the app.Logger interface.
func (l *Logger) Error(v ...interface{}) {
	printPrefix("Error", errColor)
	fmt.Println(v...)
}

// Errorf satisfies the app.Logger interface.
func (l *Logger) Errorf(format string, v ...interface{}) {
	printPrefix("Error", errColor)
	fmt.Printf(format, v...)
}

func printPrefix(level, color string) {
	file, line := caller()
	fmt.Printf("%s%s%s - %s: %v - ", color, strings.ToUpper(level), defaultColor, file, line)
}

func caller() (file string, line int) {
	_, file, line, _ = runtime.Caller(3)
	file = filepath.Base(file)
	return
}
