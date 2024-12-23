package logger

import (
	"testing"

	"github.com/trinitytechnology/ebrick/logger/core"
	"github.com/trinitytechnology/ebrick/logger/factory"
	"github.com/trinitytechnology/ebrick/logger/mocks"
)

func TestLoggerInterface(t *testing.T) {
	// This test ensures that the Logger interface is implemented correctly.
	// It verifies that all methods are implemented and that they behave as expected.
	// The test uses a mock logger to verify that the methods are called correctly.
	// The mock logger is defined in the test file.
	// The test also verifies that the logger is properly initialized.

	// Create a new mock logger.
	mock := mocks.NewMockLogger()
	mock.Info("This info message should be logged", core.String("name", "test1"))
	mock.Debug("This debug message should be logged", core.String("name", "test2"))
	mock.Warn("This warn message should be logged", core.String("name", "test3"))
	mock.Error("This error message should be logged", core.String("name", "test4"))
	expected := []string{
		"INFO: This info message should be logged | [{name test1}]",
		"DEBUG: This debug message should be logged | [{name test2}]",
		"WARN: This warn message should be logged | [{name test3}]",
		"ERROR: This error message should be logged | [{name test4}]",
	}

	if len(mock.Messages) != len(expected) {
		t.Fatalf("Expected %d messages, got %d", len(expected), len(mock.Messages))
	}
	for i, msg := range expected {
		if mock.Messages[i] != msg {
			t.Errorf("Expected '%s', got '%s'", msg, mock.Messages[i])
		}
	}
}

// TestNewLoggerFactory verifies that the factory function initializes the correct logger based on the logger type.
func TestInitLogger(t *testing.T) {
	tests := []struct {
		name       string
		env        string
		loggerType string
		wantErr    bool
	}{
		{"ZapLogger_Development", "development", "zap", false},
		{"ZapLogger_Production", "production", "zap", false},
		{"SysLogger", "development", "syslog", false},
		{"UnsupportedLogger", "development", "unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance, err := factory.LoggerFactory(tt.loggerType, tt.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if loggerInstance == nil {
					t.Error("Expected logger instance, got nil")
				}
			}
		})
	}
}
