package logger

import (
	"sync"

	"github.com/trinitytechnology/ebrick/logger/core"
)

type Logger struct {
	Core    core.Logger
	Options *Options
}

type Field = core.Field
type Level = core.Level

var (
	instance *Logger
	once     sync.Once
)

// InitLogger initializes the singleton logger instance.
func InitLogger(l core.Logger, opts ...Option) {
	once.Do(func() {
		options := newOptions(opts...)
		instance = &Logger{
			Core:    l,
			Options: options,
		}
	})
}

// GetLogger returns the singleton logger instance.
// It panics if the logger has not been initialized.
func GetLogger() *Logger {
	if instance == nil {
		panic("Logger not initialized. Call InitLogger first.")
	}
	return instance
}

// Logging methods delegate to the Core logger.
func (l *Logger) Debug(msg string, fields ...Field)  { l.Core.Debug(msg, fields...) }
func (l *Logger) Info(msg string, fields ...Field)   { l.Core.Info(msg, fields...) }
func (l *Logger) Warn(msg string, fields ...Field)   { l.Core.Warn(msg, fields...) }
func (l *Logger) Error(msg string, fields ...Field)  { l.Core.Error(msg, fields...) }
func (l *Logger) Fatal(msg string, fields ...Field)  { l.Core.Fatal(msg, fields...) }
func (l *Logger) DPanic(msg string, fields ...Field) { l.Core.DPanic(msg, fields...) }
func (l *Logger) Panic(msg string, fields ...Field)  { l.Core.Panic(msg, fields...) }

// Sync flushes the logger if required by the underlying implementation.
func (l *Logger) Sync() error {
	return l.Core.Sync()
}
