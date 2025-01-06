package server

import "google.golang.org/grpc"

var (
	DefaultAddress = ":0"
	DefaultName    = "go.ebrick.server"
	DefaultVersion = "latest"
)

// Server defines the interface for a server with methods to initialize, start, and stop the server,
// as well as retrieve its options.
type Server interface {
	// Retrieve the options
	Options() Options
	// Start the server
	Start() error
	// Stop the server
	Stop() error
}

// GRPCServer defines the interface for a gRPC server.
type GRPCServer interface {
	Server
	// RegisterServices allows registration of gRPC services
	RegisterServices(registerFunc func(s *grpc.Server))
}
