package server

import (
	"google.golang.org/grpc"
)

type Options struct {
	Address     string
	Name        string
	Version     string
	GRPCOptions GRPCOptions
}

// GRPCOptions holds gRPC-specific configuration options.
type GRPCOptions struct {
	MaxConcurrentStreams uint32
	Interceptor          grpc.UnaryServerInterceptor
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{
		Address: DefaultAddress,
		Name:    DefaultName,
		Version: DefaultVersion,
		GRPCOptions: GRPCOptions{
			MaxConcurrentStreams: 1000, // Default max streams
			Interceptor:          nil,  // No default interceptor
		},
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
func WithAddress(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

func WithName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func WithMaxConcurrentStreams(maxStreams uint32) Option {
	return func(o *Options) {
		o.GRPCOptions.MaxConcurrentStreams = maxStreams
	}
}

func WithInterceptor(interceptor grpc.UnaryServerInterceptor) Option {
	return func(o *Options) {
		o.GRPCOptions.Interceptor = interceptor
	}
}
