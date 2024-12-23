package core

import "context"

// Logger defines the interface for logging at various levels.
// Logger is an interface that defines a set of methods for logging messages at various levels of severity.
// The levels include Debug, Info, Warn, Error, DPanic, Panic, and Fatal. Each method accepts a message string
// and optional fields for additional context. The Sync method is used to flush any buffered log entries.

type Logger interface {
	// Debug logs messages at the debug level, typically used for low-level system information.
	Debug(msg string, fields ...Field)
	// Info logs informational messages that highlight the progress of the application.
	Info(msg string, fields ...Field)
	// Warn logs messages that indicate potentially harmful situations.
	Warn(msg string, fields ...Field)
	// Error logs error messages that indicate a failure in the application.
	Error(msg string, fields ...Field)
	// DPanic logs messages at the DPanic level, which is between Debug and Panic levels. It is used to log critical issues that are considered panics in development environments but are treated differently in production.
	DPanic(msg string, fields ...Field)
	// Panic logs messages at the panic level and then panics.
	Panic(msg string, fields ...Field)
	// Fatal logs messages at the fatal level and then calls os.Exit(1).
	Fatal(msg string, fields ...Field)
	// Sync flushes any buffered log entries.
	Sync() error
}

type Option func(*Options)

type Options struct {
	// Alternative options
	Context context.Context
	// fields to always be logged
	Fields []Field
	Level  Level
	Mode   string
}
