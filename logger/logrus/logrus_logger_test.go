package logrus

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	core "github.com/trinitytechnology/ebrick/logger/core"
)

func TestLogrusLoggerLevels(t *testing.T) {
	// Initialize a new logger instance for testing with default settings
	lrLogger := logrus.New()

	// Set the logger level to Debug to capture all log levels
	lrLogger.SetLevel(logrus.DebugLevel)

	// Create a test hook attached to the logger instance
	hook := test.NewLocal(lrLogger)

	// Initialize your custom LogrusLogger with the configured lrLogger
	loggerInstance := &logrusLogger{
		Logger: lrLogger,
	}

	// Define test cases covering various log levels
	testCases := []struct {
		name      string
		level     logrus.Level
		message   string
		fields    core.Field
		expectMsg string
	}{
		{
			name:      "InfoLevel",
			level:     logrus.InfoLevel,
			message:   "Information message",
			fields:    core.String("user", "bob"),
			expectMsg: "Information message",
		},
		{
			name:      "ErrorLevel",
			level:     logrus.ErrorLevel,
			message:   "Error occurred",
			fields:    core.String("error", "sample error"),
			expectMsg: "Error occurred",
		},
		{
			name:      "DebugLevel",
			level:     logrus.DebugLevel,
			message:   "Debugging message",
			fields:    core.String("debug", "true"),
			expectMsg: "Debugging message",
		},
		{
			name:      "WarnLevel",
			level:     logrus.WarnLevel,
			message:   "Warning message",
			fields:    core.String("warning", "true"),
			expectMsg: "Warning message",
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		// Capture the current test case to avoid variable sharing
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Clear previous entries to isolate test cases
			hook.Entries = nil

			// Log the message with the specified level and fields
			switch tc.level {
			case logrus.InfoLevel:
				loggerInstance.Info(tc.message, tc.fields)
			case logrus.ErrorLevel:
				loggerInstance.Error(tc.message, tc.fields)
			case logrus.DebugLevel:
				loggerInstance.Debug(tc.message, tc.fields)
			case logrus.WarnLevel:
				loggerInstance.Warn(tc.message, tc.fields)
				// Handle other levels as needed.
			default:
				t.Fatalf("Unhandled log level: %v", tc.level)
			}

			// Ensure exactly one log entry was recorded
			if len(hook.Entries) != 1 {
				t.Fatalf("Expected 1 log entry for test case '%s', got %d", tc.name, len(hook.Entries))
			}

			lastEntry := hook.LastEntry()

			// Assert the log level
			if lastEntry.Level != tc.level {
				t.Errorf("Expected level %v, got %v", tc.level, lastEntry.Level)
			}

			// Assert the log message
			if lastEntry.Message != tc.expectMsg {
				t.Errorf("Expected message '%s', got '%s'", tc.expectMsg, lastEntry.Message)
			}

			// Assert the fields
			// Assuming core.Field has Key and Value
			if lastEntry.Data[tc.fields.Key] != tc.fields.Value {
				t.Errorf("For field '%s', expected value '%v', got '%v'", tc.fields.Key, tc.fields.Value, lastEntry.Data[tc.fields.Key])
			}
		})
	}
}
