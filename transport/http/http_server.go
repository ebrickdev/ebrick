package http

import (
	"github.com/ebrickdev/ebrick/transport"
)

// Server defines the abstraction for the web server.
type HTTPServer interface {
	Engine() *Engine
	transport.Server
}

// Routable is an optional interface that modules can implement to register HTTP routes.
type Routable interface {
	RegisterRoutes(router RouterGroup)
}
