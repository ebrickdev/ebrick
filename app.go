package ebrick

import (
	"context"
	"errors"
	"sync"

	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/module"
	"github.com/ebrickdev/ebrick/transport/grpc"
	"github.com/ebrickdev/ebrick/transport/http"
)

// Application defines the interface for the application.
type Application interface {
	RegisterModules(ctx context.Context, modules ...module.Module) error
	GrpcServer() grpc.GRPCServer
	HTTPServer() http.HTTPServer
	Start(ctx context.Context) error
	Options() *Options
}

// application is the implementation of the Application interface.
type application struct {
	mm         *module.ModuleManager
	httpServer http.HTTPServer
	grpcServer grpc.GRPCServer
	options    *Options
}

// GrpcServer returns the gRPC server instance.
func (app *application) GrpcServer() grpc.GRPCServer {
	return app.grpcServer
}

// HTTPServer returns the web server instance.
func (app *application) HTTPServer() http.HTTPServer {
	return app.httpServer
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

	// Step 1: Register routes for web and gRPC services before starting any modules or servers.
	if err := app.registerRoutesAndServices(log); err != nil {
		return err
	}

	// Step 2: Start all modules
	if err := app.mm.StartAllModules(ctx); err != nil {
		return err
	}

	// Step 3: Start web server and gRPC server concurrently
	var wg sync.WaitGroup
	var combinedErrMu sync.Mutex
	var combinedErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.httpServer.Start(); err != nil {
			combineError(&combinedErrMu, &combinedErr, err)
		}
	}()

	if app.options.GRPCServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := app.grpcServer.Start(); err != nil {
				combineError(&combinedErrMu, &combinedErr, err)
			}
		}()
	}

	wg.Wait()
	return combinedErr
}

// registerRoutesAndServices registers all routes for modules implementing http.Routable or grpc.ServiceRegistrar.
func (app *application) registerRoutesAndServices(log logger.Logger) error {
	for _, mod := range app.mm.GetModules() {
		if svcReg, ok := mod.(grpc.ServiceRegistrar); ok {
			log.Info("Registering gRPC service for module", logger.String("module", mod.Name()))
			svcReg.RegisterGRPCServices(app.grpcServer)
		}
	}
	return nil
}

// combineError safely combines errors using a mutex.
func combineError(mu *sync.Mutex, combinedErr *error, err error) {
	if err != nil {
		mu.Lock()
		defer mu.Unlock()
		*combinedErr = errors.Join(*combinedErr, err)
	}
}
