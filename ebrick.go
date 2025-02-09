package ebrick

import (
	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/messaging"
	"github.com/ebrickdev/ebrick/module"
	"github.com/ebrickdev/ebrick/security/auth"
	"github.com/ebrickdev/ebrick/transport/grpc"
	"github.com/ebrickdev/ebrick/transport/httpserver"
)

var (
	EventBus    messaging.EventBus
	Logger      logger.Logger
	HTTPServer  httpserver.HTTPServer
	GRPCServer  grpc.GRPCServer
	AuthManager auth.AuthManager
)

func NewApplication(opts ...Option) Application {
	options := newOptions(opts...)

	// Set the global variables here.
	EventBus = options.EventBus
	Logger = options.Logger
	HTTPServer = options.HTTPServer
	GRPCServer = options.GRPCServer
	AuthManager = options.AuthManager

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
