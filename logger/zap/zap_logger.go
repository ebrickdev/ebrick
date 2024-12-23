package zap

import (
	"fmt"

	"github.com/trinitytechnology/ebrick/logger/core"
	"go.uber.org/zap"
)

// ZapLogger is a concrete implementation of the Logger interface using zap.Logger.
type ZapLogger struct {
	*zap.Logger
}

// NewZapLogger initializes a new zap-based Logger.
// It accepts the environment as a parameter to configure the logger appropriately.
func NewZapLogger(env string, fields ...core.Field) (core.Logger, error) {
	var zapConfig zap.Config

	switch env {
	case "production":
		zapConfig = zap.NewProductionConfig()
	default:
		zapConfig = zap.NewDevelopmentConfig()
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize zap logger: %w", err)
	}
	// set default fields
	logger = logger.With(convertToZapFields(fields)...)
	logger.Info("Logger initialized", zap.String("environment", env))
	return &ZapLogger{logger}, nil
}

// convertToZapFields converts custom Field instances to zap.Field with type-specific handling.
func convertToZapFields(fields []core.Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		switch v := field.Value.(type) {
		case string:
			zapFields = append(zapFields, zap.String(field.Key, v))
		case int:
			zapFields = append(zapFields, zap.Int(field.Key, v))
		case bool:
			zapFields = append(zapFields, zap.Bool(field.Key, v))
		case error:
			zapFields = append(zapFields, zap.Error(v))
		default:
			zapFields = append(zapFields, zap.Any(field.Key, v))
		}
	}
	zap.Fields()
	return zapFields
}

// Implementing the Logger interface methods for ZapLogger.
func (z *ZapLogger) Debug(msg string, fields ...core.Field) {
	z.Logger.Debug(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Info(msg string, fields ...core.Field) {
	z.Logger.Info(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...core.Field) {
	z.Logger.Warn(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...core.Field) {
	z.Logger.Error(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) DPanic(msg string, fields ...core.Field) {
	z.Logger.DPanic(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Panic(msg string, fields ...core.Field) {
	z.Logger.Panic(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Fatal(msg string, fields ...core.Field) {
	z.Logger.Fatal(msg, convertToZapFields(fields)...)
}

func (z *ZapLogger) Sync() error {
	return z.Logger.Sync()
}
