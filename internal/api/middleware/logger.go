package middleware

import (
	"time"

	"clusterix-code/internal/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger middleware for HTTP request logging
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		requestID := c.GetString("RequestID")

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		status := c.Writer.Status()

		// Get error if exists
		var err error
		if c.Errors.Last() != nil {
			err = c.Errors.Last().Err
		}

		// Create logger fields
		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("request_id", requestID),
		}

		// Add error field if exists
		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		// Log based on status code
		switch {
		case status >= 500:
			logger.Error("Server error", nil, fields...)
		case status >= 400:
			logger.Warn("Client error", fields...)
		case status >= 300:
			logger.Info("Redirect", fields...)
		default:
			logger.Info("Success", fields...)
		}
	}
}
