package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type grpcServer struct {
	server  *grpc.Server
	options Options
	address string
}

// NewGRPCServer creates a new gRPC server with the given options.
func NewGRPCServer(opts ...Option) GRPCServer {
	options := newOptions()
	for _, o := range opts {
		o(&options)
	}

	return &grpcServer{
		server: grpc.NewServer(
			grpc.MaxConcurrentStreams(options.GRPCOptions.MaxConcurrentStreams),
			grpc.UnaryInterceptor(options.GRPCOptions.Interceptor),
		), // Create a new gRPC server
		options: options,
		address: options.Address,
	}
}

// Options returns the server options.
func (s *grpcServer) Options() Options {
	return s.options
}

// Start starts the gRPC server and listens on the configured address.
func (s *grpcServer) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	log.Printf("GRPC: Starting gRPC server on %s...", s.address)
	go func() {
		if err := s.server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the gRPC server.
func (s *grpcServer) Stop() error {
	log.Println("Stopping gRPC server...")
	s.server.GracefulStop()
	return nil
}

// RegisterServices registers gRPC services with the server.
func (s *grpcServer) RegisterService(registerFunc func(s *grpc.Server)) {
	registerFunc(s.server)
}
