package ebrick

import (
	"context"
	"errors"
	"sync"

	"github.com/ebrickdev/ebrick/grpc"
	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/module"
	"github.com/ebrickdev/ebrick/web"
)

// Application defines the interface for the application.
type Application interface {
	RegisterModules(ctx context.Context, modules ...module.Module) error
	GrpcServer() grpc.GRPCServer
	WebServer() web.WebServer
	Start(ctx context.Context) error
	Options() *Options
}

// application is the implementation of the Application interface.
type application struct {
	mm         *module.ModuleManager
	webServer  web.WebServer
	grpcServer grpc.GRPCServer
	options    *Options
}

// GrpcServer returns the gRPC server instance.
func (app *application) GrpcServer() grpc.GRPCServer {
	return app.grpcServer
}

// WebServer returns the web server instance.
func (app *application) WebServer() web.WebServer {
	return app.webServer
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
	log := app.options.Logger

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
				log.Info("Registering routes for module", logger.String("module", mod.Name()))
				routable.RegisterRoutes(app.options.WebServer)
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
		if err := app.options.WebServer.Start(); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.options.GRPCServer.Start(); err != nil {
			mu.Lock()
			combinedErr = errors.Join(combinedErr, err)
			mu.Unlock()
		}
	}()

	wg.Wait()

	return combinedErr
}
