package middleware

import (
	"mpc/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		timestamp := time.Now()
		latency := timestamp.Sub(start)

		// Get status and error if exists
		status := c.Writer.Status()
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// Log request details
		logger.Info("Request",
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", raw),
			logger.String("ip", c.ClientIP()),
			logger.String("user-agent", c.Request.UserAgent()),
			logger.Int("status", status),
			logger.Int("latency", int(latency.Milliseconds())),
			logger.String("error", errorMessage),
		)
	}
}
