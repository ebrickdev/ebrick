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
func (app *application) GrpcServer() server.GRPCServer {
	return app.grpc
}

// WebServer returns the web server instance.
func (app *application) WebServer() web.Server {
	return app.web
}

// Options returns the application options.
func (app *application) Options() *Options {
	return app.options
}

// RegisterModules registers the provided modules with the application.
func (app *application) RegisterModules(ctx context.Context, modules ...module.Module) error {
	return app.mm.RegisterModules(ctx, modules...)
}

// Start starts the core modules, web server, and gRPC server concurrently.
func (app *application) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var combinedErr error

	// Start core modules
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Register routes for all modules
		for _, mod := range app.mm.GetModules() {
			if routable, ok := mod.(web.Routable); ok {
				logger.DefaultLogger.Info("Registering routes for module", logger.String("module", mod.Name()))
				routable.RegisterRoutes(app.web)
			}
		}

		// Start all modules
		if err := app.mm.StartAllModules(ctx); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	// Start web server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.web.Start(); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.grpc.Start(); err != nil {
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
