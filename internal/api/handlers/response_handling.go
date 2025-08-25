package handlers

import (
	internalErrors "clusterix-code/internal/utils/errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response
type Response struct {
	Success bool                `json:"success"`
	Data    interface{}         `json:"data,omitempty"`
	Error   interface{}         `json:"error,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

// SuccessResponse sends a standardized success response
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, err error) {
	var appErr *internalErrors.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPCode, Response{
			Success: false,
			Error:   appErr.Message,
			Errors:  appErr.Errors,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error:   err.Error(),
	})
}
