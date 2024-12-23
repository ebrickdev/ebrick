package logger

import (
	"context"

	"github.com/trinitytechnology/ebrick/config"
)

type Option func(*Options)

type Options struct {
	// Alternative options
	Context context.Context
	// fields to always be logged
	Fields []Field
	// Caller skip frame count for file:line info
	CallerSkipCount int
	// The logging level the logger should log at. default is `InfoLevel`
	Level Level

	Mode string

	LogType string
}

func newOptions(opts ...Option) *Options {
	mode := config.GetConfig().Env
	if mode == "" {
		mode = "development"
	}
	logType := config.GetConfig().LogType
	if mode == "" {
		mode = "default"
	}
	opt := &Options{
		Context:         nil,
		Fields:          []Field{},
		CallerSkipCount: 0,
		Level:           0,
		Mode:            mode,
		LogType:         logType,
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

// WithFields set default fields for the logger.
func WithFields(fields ...Field) Option {
	return func(args *Options) {
		args.Fields = fields
	}
}

// WithLevel set default level for the logger.
func WithLevel(level Level) Option {
	return func(args *Options) {
		args.Level = level
	}
}

// WithCallerSkipCount set frame count to skip.
func WithCallerSkipCount(c int) Option {
	return func(args *Options) {
		args.CallerSkipCount = c
	}
}
