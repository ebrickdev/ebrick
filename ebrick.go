package ebrick

import (
	"github.com/ebrickdev/ebrick/module"
)

// NewApplication creates a new instance of Application with the provided options.
// It initializes the application configuration, web server, gRPC server, and module manager.
//
// Parameters:
//
//	opts - A variadic list of Option functions to configure the application.
//
// Returns:
//
//	An instance of Application.
func NewApplication(opts ...Option) Application {
	options := newOptions(opts...)

	moduleManager := module.NewModuleManager(
		module.WithLogger(options.Logger),
		module.WithCache(options.Cache),
		module.WithEventBus(options.EventBus),
	)

	app := &application{
		mm:         moduleManager,
		httpServer: options.HTTPServer,
		grpcServer: options.GRPCServer,
		options:    options,
	}
	return app
}
