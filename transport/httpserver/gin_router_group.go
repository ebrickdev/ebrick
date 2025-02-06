package httpserver

import (
	"github.com/gin-gonic/gin"
)

// ginRouterGroup implements the RouterGroup interface
type ginRouterGroup struct {
	group *gin.RouterGroup
}

// Group implements RouterGroup.
func (g *ginRouterGroup) Group(prefix string) RouterGroup {
	return &ginRouterGroup{group: g.group.Group(prefix)}
}

func (g *ginRouterGroup) GET(route string, handler HandlerFunc) {
	g.group.GET(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (g *ginRouterGroup) POST(route string, handler HandlerFunc) {
	g.group.POST(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (g *ginRouterGroup) PUT(route string, handler HandlerFunc) {
	g.group.PUT(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (g *ginRouterGroup) PATCH(route string, handler HandlerFunc) {
	g.group.PATCH(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (g *ginRouterGroup) DELETE(route string, handler HandlerFunc) {
	g.group.DELETE(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (g *ginRouterGroup) OPTIONS(route string, handler HandlerFunc) {
	g.group.OPTIONS(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (g *ginRouterGroup) HEAD(route string, handler HandlerFunc) {
	g.group.HEAD(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}
