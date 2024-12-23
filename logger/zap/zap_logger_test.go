package zap

import (
	"fmt"
	"testing"

	"github.com/trinitytechnology/ebrick/logger/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLoggerLevels(t *testing.T) {
	// Create a zap observer to capture log entries.
	zcore, recorded := observer.New(zapcore.DebugLevel)
	zLogger := zap.New(zcore)
	loggerInstance := &ZapLogger{
		Logger: zLogger,
	}

	// Define test cases.
	testCases := []struct {
		name      string
		level     zapcore.Level
		message   string
		fields    []core.Field
		expectMsg string
	}{
		{
			name:      "InfoLevel",
			level:     zap.InfoLevel,
			message:   "Information message",
			fields:    []core.Field{core.String("user", "bob")},
			expectMsg: "Information message",
		},
		{
			name:      "ErrorLevel",
			level:     zap.ErrorLevel,
			message:   "Error occurred",
			fields:    []core.Field{core.Error("error", fmt.Errorf("sample error"))},
			expectMsg: "Error occurred",
		},
		{
			name:      "DebugLevel",
			level:     zap.DebugLevel,
			message:   "Debugging message",
			fields:    []core.Field{core.String("debug", "true")},
			expectMsg: "Debugging message",
		},
		{
			name:      "WarnLevel",
			level:     zap.WarnLevel,
			message:   "Warning message",
			fields:    []core.Field{core.String("warning", "true")},
			expectMsg: "Warning message",
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.level {
			case zap.InfoLevel:
				loggerInstance.Info(tc.message, tc.fields...)
			case zap.ErrorLevel:
				loggerInstance.Error(tc.message, tc.fields...)
			case zap.DebugLevel:
				loggerInstance.Debug(tc.message, tc.fields...)
			case zap.WarnLevel:
				loggerInstance.Warn(tc.message, tc.fields...)
				// Handle other levels as needed.
			}

			// Retrieve the last log entry.
			if recorded.Len() < 1 {
				t.Fatalf("No log entries recorded for test case '%s'", tc.name)
			}

			entry := recorded.All()[recorded.Len()-1]

			// Assert the log level.
			if entry.Level != tc.level {
				t.Errorf("Expected level %v, got %v", tc.level, entry.Level)
			}

			// Assert the log message.
			if entry.Message != tc.expectMsg {
				t.Errorf("Expected message '%s', got '%s'", tc.expectMsg, entry.Message)
			}

			// Assert the fields.
			for _, field := range tc.fields {
				found := false
				for _, f := range entry.Context {
					if f.Key == field.Key {
						switch v := field.Value.(type) {
						case string:
							if f.String != v {
								t.Errorf("For field '%s', expected value '%v', got '%v'", field.Key, v, f.String)
							}
						case error:
							if f.Interface != v {
								t.Errorf("For field '%s', expected value '%v', got '%v'", field.Key, v, f.Interface)
							}
						// Add more cases as needed for different types.
						default:
							t.Errorf("Unhandled field type for key '%s'", field.Key)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected field '%s' not found in log entry", field.Key)
				}
			}
		})
	}
}
