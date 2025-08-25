package middleware

import (
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {

		payload, exists := c.Get("authTokenPayload")

		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		tokenPayload := payload.(*dto.TokenPayload)

		permissionService := services.NewPermissionService()
		if !permissionService.IsAdmin(tokenPayload.Roles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}
		c.Next()
	}
}
