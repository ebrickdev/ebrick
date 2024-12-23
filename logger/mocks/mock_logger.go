package mocks

import (
	"fmt"
	"sync"

	"github.com/trinitytechnology/ebrick/logger/core"
)

// MockLogger is a mock implementation of the Logger interface for testing purposes.
type MockLogger struct {
	mu       sync.Mutex
	Messages []string
}

// NewMockLogger creates a new instance of MockLogger.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		Messages: []string{},
	}
}

func (m *MockLogger) Debug(msg string, fields ...core.Field) {
	m.record("DEBUG", msg, fields)
}

func (m *MockLogger) Info(msg string, fields ...core.Field) {
	m.record("INFO", msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...core.Field) {
	m.record("WARN", msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...core.Field) {
	m.record("ERROR", msg, fields)
}

func (m *MockLogger) DPanic(msg string, fields ...core.Field) {
	m.record("DPANIC", msg, fields)
}

func (m *MockLogger) Panic(msg string, fields ...core.Field) {
	m.record("PANIC", msg, fields)
}

func (m *MockLogger) Fatal(msg string, fields ...core.Field) {
	m.record("FATAL", msg, fields)
}

func (m *MockLogger) Sync() error {
	return nil
}

func (m *MockLogger) record(level, msg string, fields []core.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Messages = append(m.Messages, fmt.Sprintf("%s: %s | %v", level, msg, fields))
}
