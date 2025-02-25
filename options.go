package ebrick

import (
	"fmt"

	"github.com/ebrickdev/ebrick/cache"
	"github.com/ebrickdev/ebrick/config"
	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/messaging"
	"github.com/ebrickdev/ebrick/security/auth"
	"github.com/ebrickdev/ebrick/transport/grpc"
	"github.com/ebrickdev/ebrick/transport/http"
	"gorm.io/gorm"
)

// Options holds both configuration values and runtime dependencies
type Options struct {
	Name        string             // Application name
	Version     string             // Application version
	Cache       cache.Cache        // Cache instance
	Logger      logger.Logger      // Logger instance
	EventBus    messaging.EventBus // Event bus instance for inter-component communication
	http        http.HTTPServer    // HTTP server instance
	GRPCServer  grpc.GRPCServer    // gRPC server instance; optional
	DB          *gorm.DB
	AuthManager auth.AuthManager
}

// Option defines a function type to configure Options
type Option func(*Options)

// newOptions initializes Options with default instances for missing dependencies.
// It uses the application configuration from config.GetAppConfig().
func newOptions(opts ...Option) *Options {
	// Retrieve the application configuration
	cfg := config.GetAppConfig()

	// Initialize Options with a default Logger (could be overridden by WithLogger)
	opt := &Options{}

	// Apply functional options to override default values
	for _, o := range opts {
		o(opt)
	}

	// Ensure Logger is initialized
	if opt.Logger == nil {
		// Create a new default logger based on the environment
		logger.DefaultLogger = logger.New(logger.NewDefaultLogger(cfg.Env))
		opt.Logger = logger.DefaultLogger
		opt.Logger.Info("Default logger initiated")
	}

	// Initialize EventBus if not provided
	if opt.EventBus == nil {
		// Attempt to create an in-memory event bus. Log error if initialization fails.
		eventBus, err := messaging.NewMemoryEventBus()
		if err != nil {
			opt.Logger.Error("Failed to create event bus")
		}
		opt.EventBus = eventBus
	}

	// Initialize http if not provided.
	// Newhttp uses the application config to configure the server.
	if opt.http == nil {
		var webMode string
		if cfg.Env == "development" {
			webMode = "debug"
		} else {
			webMode = "release"
		}

		opt.http = http.NewHTTPServer(
			http.WithAddress(fmt.Sprintf(":%s", cfg.Server.Port)),
			http.WithMode(webMode),
		)
	}

	// Conditionally initialize GRPCServer if not provided and if enabled in config.
	if opt.GRPCServer == nil {
		grpcConfig, err := grpc.GetConfig()
		if err != nil {
			opt.Logger.Fatal("failed to create grpc server", logger.Error(err))
		}
		if grpcConfig.Grpc.Enabled {
			opt.GRPCServer = grpc.NewGRPCServer(grpc.WithAddress(grpcConfig.Grpc.Address))
		}
	}

	return opt
}

// WithVersion sets the application version.
func WithVersion(version string) Option {
	return func(o *Options) { o.Version = version }
}

// WithName sets the application name.
func WithName(name string) Option {
	return func(o *Options) { o.Name = name }
}

// WithCache sets the Cache dependency.
func WithCache(c cache.Cache) Option {
	return func(o *Options) { o.Cache = c }
}

// WithLogger sets the Logger dependency.
func WithLogger(l logger.Logger) Option {
	return func(o *Options) { o.Logger = l }
}

// WithEventBus sets the EventBus dependency.
func WithEventBus(eventBus messaging.EventBus) Option {
	return func(o *Options) { o.EventBus = eventBus }
}

// Withhttp sets the http dependency.
func Withhttp(http http.HTTPServer) Option {
	return func(o *Options) { o.http = http }
}

// WithGRPCServer sets the GRPCServer dependency.
func WithGRPCServer(grpcServer grpc.GRPCServer) Option {
	return func(o *Options) { o.GRPCServer = grpcServer }
}

// WithDB sets the DB dependency.
func WithDB(db *gorm.DB) Option {
	return func(o *Options) { o.DB = db }
}

func WithAuth(authManager auth.AuthManager) Option {
	return func(o *Options) { o.AuthManager = authManager }
}
