package ebrick

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ebrickdev/ebrick/config"
	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/module"
	"github.com/ebrickdev/ebrick/server"
	"github.com/ebrickdev/ebrick/web"
)

type Application interface {
	RegisterModules(ctx context.Context, modules ...module.Module) error
	Start(ctx context.Context) error
	Options() *Options
	WebServer() web.Server
	GrpcServer() server.GRPCServer
}

type application struct {
	mm      *module.ModuleManager
	http    web.Server
	gprc    server.GRPCServer
	options *Options
}

// GrpcServer implements Application.
func (a *application) GrpcServer() server.GRPCServer {
	return a.gprc
}

// Web implements Application.
func (a *application) WebServer() web.Server {
	return a.http
}

// GetOptions implements Application.
func (a *application) Options() *Options {
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
		if err := a.http.Start(); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	// Start Grpc Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.gprc.Start(); err != nil {
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

	httpServer := web.NewGinServer(
		web.WithAddress(fmt.Sprintf(":%s", appCfg.Server.Port)),
		web.WithMode(webMode),
	)

	grpcConfig, err := server.GetConfig()
	if err != nil {
		options.Logger.Fatal("failed to get grpc server config", logger.Error(err))
	}

	gprcServer := server.NewGRPCServer(server.WithAddress(grpcConfig.Grpc.Address))

	moduleManager := module.NewModuleManager(
		module.WithLogger(options.Logger),
		module.WithCache(options.Cache),
		module.WithEventBus(options.EventBus),
	)

	return &application{
		mm:      moduleManager,
		options: options,
		http:    httpServer,
		gprc:    gprcServer,
	}
}
