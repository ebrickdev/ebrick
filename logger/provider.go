package logger

type Provider interface {
	// Debug logs messages at the debug level, typically used for low-level system information.
	Debug(msg string, fields map[string]any)
	// Info logs informational messages that highlight the progress of the application.
	Info(msg string, fields map[string]any)
	// Warn logs messages that indicate potentially harmful situations.
	Warn(msg string, fields map[string]any)
	// Error logs error messages that indicate a failure in the application.
	Error(msg string, fields map[string]any)
	// DPanic logs messages at the DPanic level, which is between Debug and Panic levels. It is used to log critical issues that are considered panics in development environments but are treated differently in production.
	DPanic(msg string, fields map[string]any)
	// Panic logs messages at the panic level and then panics.
	Panic(msg string, fields map[string]any)
	// Fatal logs messages at the fatal level and then calls os.Exit(1).
	Fatal(msg string, fields map[string]any)
	// Sync flushes any buffered log entries.
	Sync() error
}
