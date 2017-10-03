package app

import "log"

type logger struct {
}

// Log satisfies the app.Logger interface.
func (l *logger) Log(v ...interface{}) {
	log.Println(v...)
}

// Logf satisfies the app.Logger interface.
func (l *logger) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Error satisfies the app.Logger interface.
func (l *logger) Error(v ...interface{}) {
	log.Println(v...)
}

// Errorf satisfies the app.Logger interface.
func (l *logger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
