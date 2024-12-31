package logger

// Logger defines the interface for logging at various levels.
// Logger is an interface that defines a set of methods for logging messages at various levels of severity.
// The levels include Debug, Info, Warn, Error, DPanic, Panic, and Fatal. Each method accepts a message string
// and optional fields for additional context. The Sync method is used to flush any buffered log entries.
var DefaultLogger Logger

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

type logger struct {
	provider Provider
}

// NewLogger creates a new instance of Logger.
func New(provider Provider) Logger {
	return &logger{provider}
}

// DPanic implements Logger.
func (l *logger) DPanic(msg string, fields ...Field) {
	l.provider.DPanic(msg, l.convertToProviderFields(fields))
}

// Debug implements Logger.
func (l *logger) Debug(msg string, fields ...Field) {
	l.provider.Debug(msg, l.convertToProviderFields(fields))
}

// Error implements Logger.
func (l *logger) Error(msg string, fields ...Field) {
	l.provider.Error(msg, l.convertToProviderFields(fields))
}

// Fatal implements Logger.
func (l *logger) Fatal(msg string, fields ...Field) {
	l.provider.Fatal(msg, l.convertToProviderFields(fields))
}

// Info implements Logger.
func (l *logger) Info(msg string, fields ...Field) {
	l.provider.Info(msg, l.convertToProviderFields(fields))
}

// Panic implements Logger.
func (l *logger) Panic(msg string, fields ...Field) {
	l.provider.Panic(msg, l.convertToProviderFields(fields))
}

// Sync implements Logger.
func (l *logger) Sync() error {
	return l.provider.Sync()

}

// Warn implements Logger.
func (l *logger) Warn(msg string, fields ...Field) {
	l.provider.Warn(msg, l.convertToProviderFields(fields))
}

func (l *logger) convertToProviderFields(fields []Field) map[string]any {
	if len(fields) == 0 {
		return nil
	}
	logFields := make(map[string]any, len(fields)) // Preallocate with size
	for _, field := range fields {
		logFields[field.Key] = field.Value
	}
	return logFields
}
