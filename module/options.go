package module

import (
	"github.com/ebrickdev/ebrick/cache"
	"github.com/ebrickdev/ebrick/event"
	"github.com/ebrickdev/ebrick/logger"
)

type Options struct {
	Cache    cache.Cache
	Logger   logger.Logger
	EventBus event.EventBus
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	opt := &Options{}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

func WithCache(c cache.Cache) Option {
	return func(o *Options) {
		o.Cache = c
	}
}

func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func WithEventBus(e event.EventBus) Option {
	return func(o *Options) {
		o.EventBus = e
	}
}
