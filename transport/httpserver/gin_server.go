package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

// ginEngine implements the Server interface using Gin.
type ginEngine struct {
	engine *gin.Engine
	server *http.Server
}

// NewGinEngine creates a new Gin-based engine with the provided options.
func NewHTTPServer(opts ...Option) HTTPServer {
	options := newOptions(opts...)

	// Set Gin mode (e.g., release, debug, test)
	gin.SetMode(options.Mode)

	// Create the Gin engine.
	// Using gin.New() allows for custom middleware configuration.
	engine := gin.New()

	// Apply custom logger middleware.
	engine.Use(options.Logger)

	// Apply additional middleware from options.
	for _, m := range options.Middleware {
		// Wrap each middleware to convert Gin context to our custom context.
		engine.Use(func(c *gin.Context) {
			m(NewGinContext(c))
		})
	}

	// Construct and return the ginEngine.
	return &ginEngine{
		engine: engine,
		server: &http.Server{
			Addr:    options.Address,
			Handler: engine,
		},
	}
}

func (s *ginEngine) GET(route string, handler HandlerFunc) {
	s.engine.GET(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginEngine) POST(route string, handler HandlerFunc) {
	s.engine.POST(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginEngine) PUT(route string, handler HandlerFunc) {
	s.engine.PUT(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginEngine) PATCH(route string, handler HandlerFunc) {
	s.engine.PATCH(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginEngine) DELETE(route string, handler HandlerFunc) {
	s.engine.DELETE(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginEngine) OPTIONS(route string, handler HandlerFunc) {
	s.engine.OPTIONS(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

func (s *ginEngine) HEAD(route string, handler HandlerFunc) {
	s.engine.HEAD(route, func(c *gin.Context) {
		handler(NewGinContext(c))
	})
}

// Group creates a new RouterGroup with the specified prefix.
func (s *ginEngine) Group(prefix string) RouterGroup {
	return &ginRouterGroup{group: s.engine.Group(prefix)}
}

// Use registers one or more middleware handlers.
func (s *ginEngine) Use(middleware ...HandlerFunc) {
	for _, m := range middleware {
		s.engine.Use(func(c *gin.Context) {
			m(NewGinContext(c))
		})
	}
}

// Start launches the web server and blocks until an OS signal is received for shutdown.
func (s *ginEngine) Start() error {
	// Start the server in a separate goroutine.
	go func() {
		log.Printf("WEB: Starting HTTPServer on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Handle OS signals for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	return s.Stop()
}

// Stop gracefully shuts down the server using a configured timeout.
func (s *ginEngine) Stop() error {
	// Use a shutdown timeout; options.WriteTimeout is used here,
	// but you might consider a dedicated shutdown timeout option.
	ctx, cancel := context.WithTimeout(context.Background(), s.server.WriteTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
