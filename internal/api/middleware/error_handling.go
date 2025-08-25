package middleware

import (
	"clusterix-code/internal/api/handlers"
	"clusterix-code/internal/api/handlers/metrics"
	"clusterix-code/internal/utils/errors"
	"clusterix-code/internal/utils/logger"
	"fmt"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware handles errors and standardizes responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()
			handleError(c, err.Err)
		}
	}
}

func handleError(c *gin.Context, err error) {
	// Convert to AppError if it isn't already
	var appErr *errors.AppError
	if e, ok := err.(*errors.AppError); ok {
		appErr = e
	} else {
		appErr = errors.NewInternalError("INTERNAL_ERROR", err)
	}

	// Record metric
	recordErrorMetric(appErr)

	// Log error
	logError(appErr)

	// Send response
	c.JSON(appErr.HTTPCode, handlers.Response{
		Success: false,
		Error:   appErr,
	})
}

func recordErrorMetric(err *errors.AppError) {
	metrics.RecordError(string(err.Type), err.Code)
}

func logError(err *errors.AppError) {
	logger.Error(fmt.Sprintf("[%s] %s", err.Type, err.Message),
		err.Raw,
		zap.String("error_code", err.Code),
		zap.String("error_type", string(err.Type)),
		zap.Int("http_code", err.HTTPCode),
	)
}
