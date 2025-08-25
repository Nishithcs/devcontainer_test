package api_context

import (
	"clusterix-code/internal/constants"
	"clusterix-code/internal/data/dto"
	"errors"

	"github.com/gin-gonic/gin"
)

func AuthUser(c *gin.Context) (*dto.User, error) {
	payload, exists := c.Get("authTokenPayload")
	if !exists {
		return nil, errors.New("user not found")
	}

	authPayload := payload.(*dto.TokenPayload)
	user := authPayload.User

	return &user, nil
}

func AuthUserDepartmentId(c *gin.Context) (*int32, error) {
	payload, exists := c.Get("authTokenPayload")
	if !exists {
		return nil, errors.New("employee not found")
	}

	authPayload := payload.(*dto.TokenPayload)
	employee := authPayload.Employee

	return &employee.DepartmentId, nil
}

func GetAuthUserLang(c *gin.Context) string {
	language := c.GetHeader("Content-Language")
	if language == "" {
		return constants.EN
	}
	return language
}
