package factory

import (
	"testing"
)

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
			loggerInstance, err := LoggerFactory(tt.loggerType, tt.env)
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
