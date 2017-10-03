package log

import (
	"github.com/murlokswarm/app"
)

// MultiLogger is a logger that logs messages using multiple loggers.
// It satisfies the app.Logger interface.
type MultiLogger struct {
	loggers []app.Logger
}

// NewMultiLogger creates a logger that uses loggers.
func NewMultiLogger(loggers ...app.Logger) *MultiLogger {
	return &MultiLogger{
		loggers: loggers,
	}
}

// Log satisfies the app.Logger interface.
func (l *MultiLogger) Log(v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Log(v...)
	}
}

// Logf satisfies the app.Logger interface.
func (l *MultiLogger) Logf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Logf(format, v...)
	}
}

// Error satisfies the app.Logger interface.
func (l *MultiLogger) Error(v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Error(v...)
	}
}

// Errorf satisfies the app.Logger interface.
func (l *MultiLogger) Errorf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Errorf(format, v...)
	}
}
