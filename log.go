package app

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
