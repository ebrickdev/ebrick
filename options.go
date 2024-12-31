package ebrick

import (
	"github.com/ebrickdev/ebrick/cache"
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

	// Init Options
	opt := &Options{}

	for _, o := range opts {
		o(opt)
	}

	// // Logger
	if opt.Logger == nil {
		opt.Logger = logger.DefaultLogger
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
