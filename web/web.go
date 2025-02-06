package web

import (
	"net/http"
)

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
	// BindJSON parses the request body as JSON into the given object.
	BindJSON(obj interface{}) error
	// Next executes the next handler in the middleware chain.
	Next()
	// Request returns the underlying HTTP request.
	Request() *http.Request
	// ClientIP returns the client's IP address.
	ClientIP() string
	// SetHeader sets a header in the response.
	SetHeader(key, value string)
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
}

// Server defines the abstraction for the web server.
type WebServer interface {
	RouterGroup
	// Use registers middleware handlers.
	Use(middleware ...HandlerFunc)
	// Start launches the web server.
	Start() error
	// Stop stops the web server.
	Stop() error
}

// Routable is an optional interface that modules can implement to register HTTP routes.
type Routable interface {
	RegisterRoutes(router RouterGroup)
}
