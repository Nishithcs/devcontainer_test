package auth

import (
	"clusterix-code/internal/api/api_context"
	"clusterix-code/internal/api/handlers"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

type ShortTokenResponse struct {
	Token string `json:"token"`
}

func (h *Handler) GenerateShortAuthToken(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewAuthenticationError("user not authenticated"))
		return
	}

	expiration := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": authUser.ID,
		"org_id":  authUser.OrganizationID,
		"exp":     expiration,
		"iat":     time.Now().Unix(),
		"purpose": "ws-auth",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("AUTH_JWT_SECRET")
	if secret == "" {
		handlers.ErrorResponse(c, errors.NewInternalError("JWT_SECRET_MISSING", fmt.Errorf("JWT secret is not set")))
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		handlers.ErrorResponse(c, errors.NewInternalError("TOKEN_SIGNING_FAILED", fmt.Errorf("failed to sign token")))
		return
	}

	handlers.SuccessResponse(c, ShortTokenResponse{Token: tokenString})
}
