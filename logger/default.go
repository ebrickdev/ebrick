package logger

import (
	"log"
	"os"
)

// DefaultLogger is a custom logger that mimics LogrusProvider's methods.
type defaultLogger struct {
	defaultFields map[string]any
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
}

// NewDefaultLogger creates a new DefaultLogger.
func NewDefaultLogger(mode string) *defaultLogger {
	infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	logger := &defaultLogger{
		defaultFields: make(map[string]any),
		infoLogger:    infoLogger,
		errorLogger:   errorLogger,
		debugLogger:   nil,
	}

	// Enable debug logs only in non-production mode
	if mode != "production" {
		logger.debugLogger = debugLogger
		logger.Info("Logger initialized in debug mode", map[string]any{"mode": mode})
	} else {
		logger.debugLogger = log.New(nil, "", 0) // Disable debug logs
	}

	return logger
}

// WithContext creates a new DefaultLogger with contextual fields.
func (l *defaultLogger) WithContext(fields map[string]any) *defaultLogger {
	newLogger := *l
	newLogger.defaultFields = mergeFields(l.defaultFields, fields)
	return &newLogger
}

// Helper method to merge fields.
func mergeFields(defaultFields, fields map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range defaultFields {
		merged[k] = v
	}
	for k, v := range fields {
		merged[k] = v
	}
	return merged
}

// Debug logs a debug message.
func (l *defaultLogger) Debug(msg string, fields map[string]any) {
	if l.debugLogger != nil {
		allFields := mergeFields(l.defaultFields, fields)
		l.debugLogger.Printf("%s | %v", msg, allFields)
	}
}

// Info logs an informational message.
func (l *defaultLogger) Info(msg string, fields map[string]any) {
	allFields := mergeFields(l.defaultFields, fields)
	l.infoLogger.Printf("%s | %v", msg, allFields)
}

// Warn logs a warning message.
func (l *defaultLogger) Warn(msg string, fields map[string]any) {
	allFields := mergeFields(l.defaultFields, fields)
	l.infoLogger.Printf("WARN: %s | %v", msg, allFields)
}

// Error logs an error message.
func (l *defaultLogger) Error(msg string, fields map[string]any) {
	allFields := mergeFields(l.defaultFields, fields)
	l.errorLogger.Printf("%s | %v", msg, allFields)
}

// DPanic logs a debug panic message and panics.
func (l *defaultLogger) DPanic(msg string, fields map[string]any) {
	l.Error(msg, fields)
	panic(msg)
}

// Panic logs a panic message and panics.
func (l *defaultLogger) Panic(msg string, fields map[string]any) {
	l.Error(msg, fields)
	panic(msg)
}

// Fatal logs a fatal message and exits.
func (l *defaultLogger) Fatal(msg string, fields map[string]any) {
	allFields := mergeFields(l.defaultFields, fields)
	l.errorLogger.Printf("FATAL: %s | %v", msg, allFields)
	os.Exit(1)
}

// Sync is a no-op for DefaultLogger (similar to logrus behavior).
func (l *defaultLogger) Sync() error {
	// Standard log doesn't need synchronization.
	return nil
}
