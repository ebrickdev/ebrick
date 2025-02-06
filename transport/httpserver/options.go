package httpserver

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Options contains configurable settings for the Gin web server
type Options struct {
	Address         string
	Mode            string
	ShutdownTimeout time.Duration
	Middleware      []HandlerFunc
	Logger          gin.HandlerFunc // Custom logger middleware
}

type Option func(*Options)

// DefaultOptions returns the default configuration for the server
func newOptions(opts ...Option) *Options {
	options := Options{
		Address:         ":8080",
		Mode:            gin.ReleaseMode,
		ShutdownTimeout: 5 * time.Second,
		Logger:          gin.Logger(), // Use Gin's default logger
		Middleware:      []HandlerFunc{},
	}

	for _, o := range opts {
		o(&options)
	}

	return &options
}

// WithAddress sets the address for the server
func WithAddress(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

// WithMode sets the mode for the server
func WithMode(mode string) Option {
	return func(o *Options) {
		o.Mode = mode
	}
}

// WithShutdownTimeout sets the shutdown timeout for the server
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ShutdownTimeout = timeout
	}
}

// WithMiddleware adds middleware to the server
func WithMiddleware(middleware ...HandlerFunc) Option {
	return func(o *Options) {
		o.Middleware = append(o.Middleware, middleware...)
	}
}
