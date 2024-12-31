package server

import (
	"context"
)

type Options struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	// The router for requests
	Router Router

	Name      string
	Id        string
	Version   string
	Advertise string
	Address   string
}

type Option func(*Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func WithRouter(router Router) Option {
	return func(o *Options) {
		o.Router = router
	}
}

func WithName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func WithId(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func WithAdvertise(advertise string) Option {
	return func(o *Options) {
		o.Advertise = advertise
	}
}

func WithAddress(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}
