package httpserver

import (
	"net/http"

	"github.com/ebrickdev/ebrick/transport"
)

type H map[string]any

// HandlerFunc defines a generic HTTP handler function type.
type HandlerFunc func(ctx Context)

// Context abstracts the HTTP context used by the framework.
type Context interface {
	// JSON sends a JSON response with the specified status code.
	JSON(statusCode int, obj interface{})
	// String sends a formatted string response with the specified status code.
	String(statusCode int, format string, values ...interface{})
	// Param retrieves a route parameter by name.
	Param(key string) string
	// Query retrieves a query parameter by name.
	Query(key string) string
	// ShouldBindJSON parses the request body as JSON into the given object.
	ShouldBindJSON(obj interface{}) error
	// ShouldBindBodyWithJSON is a shortcut for c.ShouldBindBodyWith(obj, binding.JSON).
	ShouldBindBodyWithJSON(obj any) error
	// ShouldBindBodyWithXML is a shortcut for c.ShouldBindBodyWith(obj, binding.XML).
	ShouldBindBodyWithXML(obj any) error
	// ShouldBindBodyWithYAML is a shortcut for c.ShouldBindBodyWith(obj, binding.YAML).
	ShouldBindBodyWithYAML(obj any) error
	// ShouldBindBodyWithTOML is a shortcut for c.ShouldBindBodyWith(obj, binding.TOML).
	ShouldBindBodyWithTOML(obj any) error
	// Next executes the next handler in the middleware chain.
	Next()
	// Request returns the underlying HTTP request.
	Request() *http.Request
	// ClientIP returns the client's IP address.
	ClientIP() string
	// SetHeader sets a header in the response.
	SetHeader(key, value string)
	// Get retrieves a value from the context by key.
	Get(key string) (value any, exists bool)
	// Set sets a key-value pair in the context.
	Set(key string, value any)
	// AbortWithStatus aborts the request with the specified status code.
	AbortWithStatus(code int)
	// Abort aborts the request.
	Abort()
}

// Router defines the core HTTP routing interface.
type Router interface {
	GET(route string, handler HandlerFunc)
	POST(route string, handler HandlerFunc)
	PUT(route string, handler HandlerFunc)
	PATCH(route string, handler HandlerFunc)
	DELETE(route string, handler HandlerFunc)
	OPTIONS(route string, handler HandlerFunc)
	HEAD(route string, handler HandlerFunc)
}

// RouterGroup abstracts route grouping; it embeds Router.
type RouterGroup interface {
	Router
	// Group creates a new RouterGroup with the given prefix.
	Group(prefix string) RouterGroup
	Use(middleware ...HandlerFunc)
}

// Server defines the abstraction for the web server.
type HTTPServer interface {
	RouterGroup
	transport.Server
}

// Routable is an optional interface that modules can implement to register HTTP routes.
type Routable interface {
	RegisterRoutes(router RouterGroup)
}
