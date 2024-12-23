package logrus

import (
	"github.com/sirupsen/logrus"
	"github.com/trinitytechnology/ebrick/logger/core"
)

// logrusLogger is a concrete implementation of the Logger interface using logrus.Logger.
type logrusLogger struct {
	*logrus.Logger
}

// NewLogrusLogger initializes a new logrus-based Logger.
func NewLogrusLogger(env string, fields ...core.Field) (core.Logger, error) {
	logger := logrus.New()

	switch env {
	case "production":
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.InfoLevel)
	default:
		logger.SetFormatter(&logrus.TextFormatter{})
		logger.SetLevel(logrus.DebugLevel)
	}

	// set default fields
	// Add default fields if any
	if len(fields) > 0 {
		logger = logger.WithFields(convertToLogrusFields(fields)).Logger
	}

	logger.Info("Logger initialized", logrus.Fields{"environment": env})

	return &logrusLogger{logger}, nil
}

// Implementing the Logger interface methods for logrusLogger.
func (l *logrusLogger) Debug(msg string, fields ...core.Field) {
	logrusFields := convertToLogrusFields(fields)
	l.Logger.WithFields(logrusFields).Debug(msg)
}

func (l *logrusLogger) Info(msg string, fields ...core.Field) {
	logrusFields := convertToLogrusFields(fields)
	l.Logger.WithFields(logrusFields).Info(msg)
}

func (l *logrusLogger) Warn(msg string, fields ...core.Field) {
	logrusFields := convertToLogrusFields(fields)
	l.Logger.WithFields(logrusFields).Warn(msg)
}

func (l *logrusLogger) Error(msg string, fields ...core.Field) {
	logrusFields := convertToLogrusFields(fields)
	l.Logger.WithFields(logrusFields).Error(msg)
}

func (l *logrusLogger) DPanic(msg string, fields ...core.Field) {
	logrusFields := convertToLogrusFields(fields)
	l.Logger.WithFields(logrusFields).Panic(msg)

	if l.Logger.GetLevel() <= logrus.DebugLevel {
		panic(msg)
	}
}

func (l *logrusLogger) Panic(msg string, fields ...core.Field) {
	logrusFields := convertToLogrusFields(fields)
	l.Logger.WithFields(logrusFields).Panic(msg)
}

func (l *logrusLogger) Fatal(msg string, fields ...core.Field) {
	logrusFields := convertToLogrusFields(fields)
	l.Logger.WithFields(logrusFields).Fatal(msg)
}

func (l *logrusLogger) Sync() error {
	// logrus does not require explicit synchronization.
	return nil
}

// convertToLogrusFields converts custom Field instances to logrus.Fields with type-specific handling.
func convertToLogrusFields(fields []core.Field) logrus.Fields {
	logrusFields := make(logrus.Fields)
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}
	return logrusFields
}
