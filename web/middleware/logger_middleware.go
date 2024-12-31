package middleware

import (
	"time"

	"github.com/ebrickdev/ebrick/logger"
	"github.com/ebrickdev/ebrick/web"
)

// LoggerMiddleware logs request details for each HTTP request
func LoggerMiddleware(log logger.Logger) web.HandlerFunc {
	return func(ctx web.Context) {
		start := time.Now()

		ctx.Next()

		// Log details
		duration := time.Since(start)
		req := ctx.Request()
		log.Info("Request completed",
			logger.String("method", req.Method), // Adjusting for the web.Context abstraction
			logger.String("path", req.URL.Path), // Use abstraction methods for retrieving path
			logger.String("ip", ctx.ClientIP()),
			logger.String("user_agent", req.UserAgent()),
			logger.Any("latency", duration), // Time duration the request took
		)
	}
}
