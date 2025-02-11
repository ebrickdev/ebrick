package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type httpServer struct {
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
			m(c)
		})
	}

	// Construct and return the ginEngine.
	return &httpServer{
		engine: engine,
		server: &http.Server{
			Addr:    options.Address,
			Handler: engine,
		},
	}
}
func (s *httpServer) Engine() *Engine {
	return s.engine
}

// Start launches the web server and blocks until an OS signal is received for shutdown.
func (s *httpServer) Start() error {
	// Start the server in a separate goroutine.
	go func() {
		log.Printf("WEB: Starting http on %s", s.server.Addr)
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
func (s *httpServer) Stop() error {
	// Use a shutdown timeout; options.WriteTimeout is used here,
	// but you might consider a dedicated shutdown timeout option.
	ctx, cancel := context.WithTimeout(context.Background(), s.server.WriteTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
