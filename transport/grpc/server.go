package grpc

import (
	"github.com/ebrickdev/ebrick/transport"
	"google.golang.org/grpc"
)

var (
	DefaultAddress = ":0"
	DefaultName    = "go.ebrick.server"
	DefaultVersion = "latest"
)

// GRPCServer defines the interface for a gRPC server.
type GRPCServer interface {
	Options() Options

	transport.Server
	// RegisterServices allows registration of gRPC services
	RegisterService(registerFunc func(s *grpc.Server))
}

type ServiceRegistrar interface {
	RegisterGRPCServices(s GRPCServer)
}
