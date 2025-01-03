package web

import "net/http"

// HandlerFunc defines a generic HTTP handler function type
type HandlerFunc func(ctx Context)

// Context abstracts the HTTP context used by the framework
type Context interface {
	JSON(statusCode int, obj interface{})
	String(statusCode int, format string, values ...interface{})
	Param(key string) string
	Query(key string) string
	BindJSON(obj interface{}) error
	Next()
	Request() *http.Request
	ClientIP() string
	SetHeader(key, value string)
}

// RouterGroup abstracts route grouping
type RouterGroup interface {
	GET(route string, handler HandlerFunc)
	POST(route string, handler HandlerFunc)
	PUT(route string, handler HandlerFunc)
	PATCH(route string, handler HandlerFunc)
	DELETE(route string, handler HandlerFunc)
	OPTIONS(route string, handler HandlerFunc)
	HEAD(route string, handler HandlerFunc)
}

// Server defines the abstraction for the web server
type Server interface {
	GET(route string, handler HandlerFunc)
	POST(route string, handler HandlerFunc)
	PUT(route string, handler HandlerFunc)
	PATCH(route string, handler HandlerFunc)
	DELETE(route string, handler HandlerFunc)
	OPTIONS(route string, handler HandlerFunc)
	HEAD(route string, handler HandlerFunc)
	Group(prefix string) RouterGroup
	Use(middleware ...HandlerFunc)
	Start() error
	Stop() error
}
