package ebrick

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ebrickdev/ebrick/config"
	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/module"
	"github.com/ebrickdev/ebrick/web"
	"github.com/ebrickdev/ebrick/web/middleware"
)

type Application interface {
	RegisterModules(ctx context.Context, modules ...module.Module) error
	Start(ctx context.Context) error
	GetOptions() *Options
}

type application struct {
	mm      *module.ModuleManager
	web     web.Server
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

// Start implements Application.
// Starts the application and all its components
func (a *application) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var combinedErr error

	// Start Core Modules
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.mm.StartAllModules(ctx); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	// Start Web Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.web.Start(); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	// Wait for all components to start
	wg.Wait()

	return combinedErr
}

func NewApplication(opts ...Option) Application {
	appCfg := config.GetAppConfig()
	options := newOptions(opts...)

	var webMode string
	if appCfg.Env == "development" {
		webMode = "debug"
	} else {
		webMode = "release"
	}

	webserver := web.NewGinServer(
		web.WithAddress(fmt.Sprintf(":%s", appCfg.Server.Port)),
		web.WithMode(webMode),
		web.WithMiddleware(middleware.LoggerMiddleware(logger.DefaultLogger)),
	)

	moduleManager := module.NewModuleManager(
		module.WithLogger(options.Logger),
		module.WithCache(options.Cache),
		module.WithEventBus(options.EventBus),
	)

	return &application{
		mm:      moduleManager,
		options: options,
		web:     webserver,
	}
}
