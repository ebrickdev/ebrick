package ebrick

import (
	"github.com/ebrickdev/ebrick/cache"
	"github.com/ebrickdev/ebrick/config"
	"github.com/ebrickdev/ebrick/event"
	"github.com/ebrickdev/ebrick/logger"
)

type Options struct {
	Name     string
	Version  string
	Cache    cache.Cache
	Logger   logger.Logger
	EventBus event.EventBus
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	cfg := config.GetAppConfig()
	// Init Options
	opt := &Options{
		Logger: logger.DefaultLogger,
	}

	for _, o := range opts {
		o(opt)
	}

	if opt.Logger == nil {
		logger.DefaultLogger = logger.New(logger.NewDefaultLogger(cfg.Env))
		opt.Logger = logger.DefaultLogger
		logger.DefaultLogger.Info("Default logger initiated")
	}

	// // Init Event Bus
	// if opt.EventBus == nil {
	// 	// EventBus
	// 	eventBus, err := inmemory.NewEventBus()
	// 	if err != nil {
	// 		opt.Logger.Error("Failed to create event bus")
	// 	}
	// 	opt.EventBus = eventBus
	// }

	return opt
}

func WithVersion(version string) Option {
	return func(o *Options) { o.Version = version }
}

func WithName(name string) Option {
	return func(o *Options) { o.Name = name }
}

func WithCache(c cache.Cache) Option {
	return func(o *Options) { o.Cache = c }
}

func WithLogger(l logger.Logger) Option {
	return func(o *Options) { o.Logger = l }
}

func WithEventBus(eventBus event.EventBus) Option {
	return func(o *Options) { o.EventBus = eventBus }
}
