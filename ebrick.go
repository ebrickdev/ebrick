package ebrick

import (
	"context"

	"github.com/ebrickdev/ebrick/module"
)

type Application interface {
	RegisterModules(ctx context.Context, modules ...module.Module) error
	Start(ctx context.Context) error
	GetOptions() *Options
}

type application struct {
	mm      *module.ModuleManager
	options *Options
}

// GetOptions implements Application.
func (a *application) GetOptions() *Options {
	return a.options
}

// RegisterModules implements Application.
func (a *application) RegisterModules(ctx context.Context, modules ...module.Module) error {
	return a.mm.RegisterModules(ctx, modules...)
}

func (a *application) Start(ctx context.Context) error {
	// Start Cores Services

	// Start all modules
	err := a.mm.StartAllModules(ctx)
	return err
}

func NewApplication(opts ...Option) Application {

	options := newOptions(opts...)

	moduleManager := module.NewModuleManager(module.WithLogger(options.Logger), module.WithCache(options.Cache), module.WithEventBus(options.EventBus))

	return &application{
		mm:      moduleManager,
		options: options,
	}
}
