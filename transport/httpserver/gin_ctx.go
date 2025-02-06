package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginContext struct {
	ctx *gin.Context
}

// NewGinContext wraps Gin's context
func NewGinContext(ctx *gin.Context) Context {
	return &ginContext{ctx: ctx}
}

func (g *ginContext) JSON(statusCode int, obj interface{}) {
	g.ctx.JSON(statusCode, obj)
}

func (g *ginContext) String(statusCode int, format string, values ...interface{}) {
	g.ctx.String(statusCode, format, values...)
}

func (g *ginContext) Param(key string) string {
	return g.ctx.Param(key)
}

func (g *ginContext) Query(key string) string {
	return g.ctx.Query(key)
}

func (g *ginContext) BindJSON(obj interface{}) error {
	return g.ctx.ShouldBindJSON(obj)
}

// Next implements Context.
func (g *ginContext) Next() {
	g.ctx.Next()
}

func (g *ginContext) Request() *http.Request {
	return g.ctx.Request
}

func (g *ginContext) ClientIP() string {
	return g.ctx.ClientIP()
}

func (g *ginContext) SetHeader(key, value string) {
	g.ctx.Header(key, value)
}
