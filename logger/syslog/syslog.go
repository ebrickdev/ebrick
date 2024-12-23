package syslog

import (
	"fmt"
	"log/syslog"

	"github.com/trinitytechnology/ebrick/logger/core"
)

// SyslogLogger implements the core.Logger interface using Syslog.
type SyslogLogger struct {
	logger *syslog.Writer
}

// NewSysLogger initializes a new SyslogLogger.
func NewSysLogger() (core.Logger, error) {
	writer, err := syslog.New(syslog.LOG_INFO|syslog.LOG_USER, "myApp")
	if err != nil {
		return nil, err
	}

	// Example of adding default fields
	writer.Info("Syslog logger initialized")

	return &SyslogLogger{logger: writer}, nil
}

// Implementing core.Logger methods

func (s *SyslogLogger) Debug(msg string, fields ...core.Field) {
	// Syslog doesn't have a Debug level; map to Info
	s.logger.Info(formatMessage(msg, fields))
}

func (s *SyslogLogger) Info(msg string, fields ...core.Field) {
	s.logger.Info(formatMessage(msg, fields))
}

func (s *SyslogLogger) Warn(msg string, fields ...core.Field) {
	s.logger.Warning(formatMessage(msg, fields))
}

func (s *SyslogLogger) Error(msg string, fields ...core.Field) {
	s.logger.Err(formatMessage(msg, fields))
}

func (s *SyslogLogger) Fatal(msg string, fields ...core.Field) {
	s.logger.Crit(formatMessage(msg, fields))
}

func (s *SyslogLogger) DPanic(msg string, fields ...core.Field) {
	s.logger.Warning(formatMessage(msg, fields))
}

func (s *SyslogLogger) Panic(msg string, fields ...core.Field) {
	s.logger.Err(formatMessage(msg, fields))
}

func (s *SyslogLogger) Sync() error {
	return s.logger.Close()
}

// Helper function to format messages with fields
func formatMessage(msg string, fields []core.Field) string {
	formatted := msg
	for _, field := range fields {
		formatted += " " + field.Key + "=" + fmt.Sprintf("%v", field.Value)
	}
	return formatted
}
