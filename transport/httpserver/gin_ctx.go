package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ginContext wraps Gin's context to implement the Context interface.
type ginContext struct {
	ctx *gin.Context
}

// ShouldBindBodyWithJSON implements Context.
func (g *ginContext) ShouldBindBodyWithJSON(obj any) error {
	return g.ctx.ShouldBindBodyWithJSON(obj)
}

// ShouldBindBodyWithTOML implements Context.
func (g *ginContext) ShouldBindBodyWithTOML(obj any) error {
	return g.ctx.ShouldBindBodyWithTOML(obj)
}

// ShouldBindBodyWithXML implements Context.
func (g *ginContext) ShouldBindBodyWithXML(obj any) error {
	return g.ctx.ShouldBindBodyWithXML(obj)
}

// ShouldBindBodyWithYAML implements Context.
func (g *ginContext) ShouldBindBodyWithYAML(obj any) error {
	return g.ctx.ShouldBindBodyWithYAML(obj)
}

// NewGinContext wraps Gin's context and returns a new ginContext.
func NewGinContext(ctx *gin.Context) Context {
	return &ginContext{ctx: ctx}
}

// JSON sends a JSON response with the specified status code.
func (g *ginContext) JSON(statusCode int, obj interface{}) {
	g.ctx.JSON(statusCode, obj)
}

// String sends a string response with the specified status code and format.
func (g *ginContext) String(statusCode int, format string, values ...interface{}) {
	g.ctx.String(statusCode, format, values...)
}

// Param returns the value of the URL parameter with the given key.
func (g *ginContext) Param(key string) string {
	return g.ctx.Param(key)
}

// Query returns the value of the query parameter with the given key.
func (g *ginContext) Query(key string) string {
	return g.ctx.Query(key)
}

// ShouldBindJSON binds the JSON payload to the given object.
func (g *ginContext) ShouldBindJSON(obj interface{}) error {
	return g.ctx.ShouldBindJSON(obj)
}

// Next calls the next handler in the middleware chain.
func (g *ginContext) Next() {
	g.ctx.Next()
}

func (g *ginContext) Get(key string) (value any, exists bool) {
	return g.ctx.Get(key)
}

// Set sets a key-value pair in the context.
func (g *ginContext) Set(key string, value any) {
	g.ctx.Set(key, value)
}

// AbortWithStatus aborts the request with the specified status code.
func (g *ginContext) AbortWithStatus(code int) {
	g.ctx.AbortWithStatus(code)
}

// Abort aborts the request.
func (g *ginContext) Abort() {
	g.ctx.Abort()
}

// Request returns the HTTP request.
func (g *ginContext) Request() *http.Request {
	return g.ctx.Request
}

// ClientIP returns the client's IP address.
func (g *ginContext) ClientIP() string {
	return g.ctx.ClientIP()
}

// SetHeader sets a header key-value pair in the response.
func (g *ginContext) SetHeader(key, value string) {
	g.ctx.Header(key, value)

}
