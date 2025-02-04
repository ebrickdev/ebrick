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

// Application defines the interface for the application.
type Application interface {
	RegisterModules(ctx context.Context, modules ...module.Module) error
	Start(ctx context.Context) error
	Options() *Options
	WebServer() web.Server
	GrpcServer() server.GRPCServer
}

// application is the implementation of the Application interface.
type application struct {
	mm      *module.ModuleManager
	web     web.Server
	grpc    server.GRPCServer
	options *Options
}

// GrpcServer returns the gRPC server instance.
func (a *application) GrpcServer() server.GRPCServer {
	return a.grpc
}

// WebServer returns the web server instance.
func (a *application) WebServer() web.Server {
	return a.web
}

// Options returns the application options.
func (a *application) Options() *Options {
	return a.options
}

// RegisterModules registers the provided modules with the application.
func (a *application) RegisterModules(ctx context.Context, modules ...module.Module) error {
	return a.mm.RegisterModules(ctx, modules...)
}

// Start starts the core modules, web server, and gRPC server concurrently.
func (a *application) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var combinedErr error

	// Start core modules
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.mm.StartAllModules(ctx); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	// Start web server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.web.Start(); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.grpc.Start(); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	wg.Wait()

	return combinedErr
}

// NewApplication creates and returns an instance of Application.
func NewApplication(opts ...Option) Application {
	appCfg := config.GetAppConfig()
	options := newOptions(opts...)

	webServer := newWebServer(appCfg)

	grpcServer, err := newGrpcServer()
	if err != nil {
		options.Logger.Fatal("failed to create grpc server", logger.Error(err))
	}

	moduleManager := module.NewModuleManager(
		module.WithLogger(options.Logger),
		module.WithCache(options.Cache),
		module.WithEventBus(options.EventBus),
	)

	app := &application{
		mm:      moduleManager,
		options: options,
		web:     webServer,
		grpc:    grpcServer,
	}

	// Register module routes
	app.registerModuleRoutes()
	return app
}

// newWebServer creates a web server based on the configuration.
func newWebServer(appCfg *config.Config) web.Server {
	var webMode string
	if appCfg.Env == "development" {
		webMode = "debug"
	} else {
		webMode = "release"
	}

	return web.NewGinEngine(
		web.WithAddress(fmt.Sprintf(":%s", appCfg.Server.Port)),
		web.WithMode(webMode),
	)
}

// newGrpcServer creates a gRPC server based on the configuration.
func newGrpcServer() (server.GRPCServer, error) {
	grpcConfig, err := server.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc server config: %w", err)
	}
	return server.NewGRPCServer(server.WithAddress(grpcConfig.Grpc.Address)), nil
}

// registerModuleRoutes iterates over the modules and registers their routes.
func (a *application) registerModuleRoutes() {
	// Assume a.web implements RouterGroup.
	routerGroup := a.web

	// Iterate over all modules registered in your module manager.
	for _, mod := range a.mm.GetModules() {
		if routable, ok := mod.(web.Routable); ok {
			routable.RegisterRoutes(routerGroup)
		}
	}
}
