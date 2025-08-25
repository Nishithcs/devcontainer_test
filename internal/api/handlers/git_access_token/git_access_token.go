package git_access_token

import (
	"clusterix-code/internal/api/api_context"
	"clusterix-code/internal/api/handlers"
	"clusterix-code/internal/api/requests"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/errors"
	"clusterix-code/internal/utils/pagination"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) GetUserAccessTokens(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()
	query := c.Query("q")
	with := c.QueryArray("with")
	page, limit := pagination.Paginate(c)

	response, err := h.services.GitPersonalAccessToken.GetUserAccessTokens(ctx, authUser.ID, query, with, page, limit)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	handlers.SuccessResponse(c, response)
}

func (h *Handler) GetUserAccessToken(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	fmt.Println(authUser)

	accessTokenId := c.Param("id")
	if accessTokenId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_ACCESS_TOKEN_ID",
			"Access Token ID is required",
			nil))
		return
	}
	id, err := strconv.ParseUint(accessTokenId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_REPOSITORY_ID",
			"Access Token ID must be a valid number",
			err))
		return
	}

	ctx := c.Request.Context()
	gitAccessTokenDTO, err := h.services.GitPersonalAccessToken.GetUserAccessToken(ctx, authUser.ID, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessGitToken(ctx, authUser, &gitAccessTokenDTO) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this git access token",
			nil))
		return
	}


	handlers.SuccessResponse(c, gitAccessTokenDTO)
}

func (h *Handler) CreateUserAccessToken(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	var req requests.CreateGitAccessTokenRequest
	req.UserID = authUser.ID
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()
	newAccessToken, err := h.services.GitPersonalAccessToken.CreateUserAccessToken(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, newAccessToken)
}

func (h *Handler) UpdateUserAccessToken(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	accessTokenId := c.Param("id")
	if accessTokenId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_REPOSITORY_ID",
			"Repository ID is required",
			nil))
		return
	}
	var req requests.UpdateGitAccessTokenRequest
	id, err := strconv.ParseUint(accessTokenId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_ACCESS_TOKEN_ID",
			"Access Token ID must be a valid number",
			err))
		return
	}
	req.ID = id
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()

	gitAccessTokenDTO, err := h.services.GitPersonalAccessToken.GetUserAccessToken(ctx, authUser.ID, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessGitToken(ctx, authUser, &gitAccessTokenDTO) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this git access token",
			nil))
		return
	}

	repo, err := h.services.GitPersonalAccessToken.UpdateUserAccessToken(ctx, authUser.ID, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}

func (h *Handler) DeleteUserAccessToken(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	fmt.Println(authUser.ID)

	gitAccessTokenId := c.Param("id")
	if gitAccessTokenId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_GIT_ACCESS_TOKEN_ID",
			"Git access token ID is required",
			nil))
		return
	}

	ctx := c.Request.Context()

	id, err := strconv.ParseUint(gitAccessTokenId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_ACCESS_TOKEN_ID",
			"Access Token ID must be a valid number",
			err))
		return
	}

	gitAccessTokenDTO, err := h.services.GitPersonalAccessToken.GetUserAccessToken(ctx, authUser.ID, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessGitToken(ctx, authUser, &gitAccessTokenDTO) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this git access token",
			nil))
		return
	}

	err = h.services.GitPersonalAccessToken.DeleteUserAccessToken(ctx, gitAccessTokenId)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, true)
}
