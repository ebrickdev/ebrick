package module

import (
	"github.com/gin-gonic/gin"
	"github.com/trinitytechnology/ebrick/cache"
	"github.com/trinitytechnology/ebrick/messaging"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Options struct {
	Database    *gorm.DB
	Cache       cache.Cache
	EventStream messaging.CloudEventStream
	Logger      *zap.Logger
	Router      *gin.Engine
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	opt := &Options{}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

// Option functions
func Database(db *gorm.DB) Option {
	return func(o *Options) {
		o.Database = db
	}
}

func Cache(c cache.Cache) Option {
	return func(o *Options) {
		o.Cache = c
	}
}

func EventStream(es messaging.CloudEventStream) Option {
	return func(o *Options) {
		o.EventStream = es
	}
}

func Logger(l *zap.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func Router(r *gin.Engine) Option {
	return func(o *Options) {
		o.Router = r
	}
}
