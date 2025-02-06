package ebrick

import (
	"fmt"

	"github.com/ebrickdev/ebrick/config"
	"github.com/ebrickdev/ebrick/grpc"
	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/module"
	"github.com/ebrickdev/ebrick/web"
)

// newWebServer creates a web server based on the configuration.
func NewWebServer(appCfg *config.Config) web.WebServer {
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
func NewGRPCServer(options *Options) grpc.GRPCServer {
	grpcConfig, err := grpc.GetConfig()
	if err != nil {
		options.Logger.Fatal("failed to create grpc server", logger.Error(err))
	}
	return grpc.NewGRPCServer(grpc.WithAddress(grpcConfig.Grpc.Address))
}

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
		webServer:  options.WebServer,
		grpcServer: options.GRPCServer,
		options:    options,
	}
	return app
}
