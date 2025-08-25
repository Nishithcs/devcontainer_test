package middleware

import (
	"clusterix-code/internal/api/handlers"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/utils/errors"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ConvertMapClaimsToTokenPayload(mapClaims jwt.MapClaims) (*dto.TokenPayload, error) {
	jsonBytes, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, fmt.Errorf("error marshalling MapClaims to JSON: %v", err)
	}

	var tokenPayload dto.TokenPayload
	if err := json.Unmarshal(jsonBytes, &tokenPayload); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON to TokenPayload: %v", err)
	}

	return &tokenPayload, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handlers.ErrorResponse(c, errors.NewAuthenticationError("Authorization header is missing"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			handlers.ErrorResponse(c, errors.NewAuthenticationError("Invalid Authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		//if err != nil || !token.Valid {
		//	handlers.ErrorResponse(c, errors.NewAuthenticationError("Invalid token"))
		//	c.Abort()
		//	return
		//}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			payload, err := ConvertMapClaimsToTokenPayload(claims)
			if err != nil {
				handlers.ErrorResponse(c, errors.NewInternalError("INTERNAL_ERROR", err))
				c.Abort()
				return
			}
			c.Set("authTokenPayload", payload)
		}

		c.Next()
	}
}

func AuthMiddlewareWithQueryParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			handlers.ErrorResponse(c, errors.NewAuthenticationError("Missing authentication token"))
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewAuthenticationError("Unexpected signing method")
			}
			return []byte(os.Getenv("AUTH_JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			handlers.ErrorResponse(c, errors.NewAuthenticationError("Invalid or expired token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handlers.ErrorResponse(c, errors.NewAuthenticationError("Invalid token claims"))
			c.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
			handlers.ErrorResponse(c, errors.NewAuthenticationError("Token has expired"))
			c.Abort()
			return
		}

		if purpose, ok := claims["purpose"].(string); !ok || purpose != "ws-auth" {
			handlers.ErrorResponse(c, errors.NewAuthenticationError("Invalid token purpose"))
			c.Abort()
			return
		}

		payload := &dto.TokenPayload{}
		if userID, ok := claims["user_id"].(float64); ok {
			payload.User.ID = uint64(userID)
		}
		if orgID, ok := claims["org_id"].(float64); ok {
			payload.User.OrganizationID = uint32(orgID)
		}

		c.Set("authTokenPayload", payload)
		c.Next()
	}
}
