package module

import "context"

type Module interface {
	Id() string
	Name() string
	Version() string
	Description() string
	Initialize(ctx context.Context, options *Options) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Dependencies() []string
}

type ModuleConfig struct {
	Id     string
	Name   string
	Enable bool
}
