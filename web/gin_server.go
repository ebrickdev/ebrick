package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

// ginServer implements the Server interface using Gin
type ginServer struct {
	engine *gin.Engine
	server *http.Server
}

func NewGinServer(opts ...Option) Server {
	options := newOptions(opts...)

	// Set Gin mode (e.g., release, debug, test)
	gin.SetMode(options.Mode)

	// Create the Gin engine
	engine := gin.New() // Use `gin.New()` instead of `gin.Default()` for custom middleware

	// Apply custom logger
	engine.Use(options.Logger)

	// Apply middleware from options
	for _, m := range options.Middleware {
		engine.Use(func(c *gin.Context) {
			m(NewGinContext(c))
		})

	}
	return &ginServer{
		engine: engine,
		server: &http.Server{
			Addr:    options.Address,
			Handler: engine,
		},
	}
}

func (s *ginServer) GET(route string, handler HandlerFunc) {
	s.engine.GET(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginServer) POST(route string, handler HandlerFunc) {
	s.engine.POST(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginServer) PUT(route string, handler HandlerFunc) {
	s.engine.PUT(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginServer) PATCH(route string, handler HandlerFunc) {
	s.engine.PATCH(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginServer) DELETE(route string, handler HandlerFunc) {
	s.engine.DELETE(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginServer) OPTIONS(route string, handler HandlerFunc) {
	s.engine.OPTIONS(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginServer) HEAD(route string, handler HandlerFunc) {
	s.engine.HEAD(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginServer) Group(prefix string) RouterGroup {
	return &ginRouterGroup{group: s.engine.Group(prefix)}
}

func (s *ginServer) Use(middleware ...HandlerFunc) {
	for _, m := range middleware {
		s.engine.Use(func(c *gin.Context) {
			m(NewGinContext(c))
		})
	}
}

func (s *ginServer) Start() error {
	// Start the server in a goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Handle OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	return s.Stop()
}

func (s *ginServer) Stop() error {
	// Use the configured shutdown timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.server.WriteTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
