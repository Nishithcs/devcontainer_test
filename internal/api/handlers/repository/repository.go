package repository

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

func (h *Handler) GetRepositories(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()
	query := c.Query("q")
	status := c.Query("status")
	with := c.QueryArray("with")
	page, limit := pagination.Paginate(c)

	response, err := h.services.Repository.GetRepositories(ctx, authUser.OrganizationID, query, status, with, page, limit)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	handlers.SuccessResponse(c, response)
}

func (h *Handler) GetRepository(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	fmt.Println(authUser)

	repositoryId := c.Param("id")
	if repositoryId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_REPOSITORY_ID",
			"Repository ID is required",
			nil))
		return
	}
	id, err := strconv.ParseUint(repositoryId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_REPOSITORY_ID",
			"Repository ID must be a valid number",
			err))
		return
	}

	ctx := c.Request.Context()
	repositoryDTO, err := h.services.Repository.GetRepository(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessRepository(ctx, authUser, &repositoryDTO) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this repository",
			nil))
		return
	}

	handlers.SuccessResponse(c, repositoryDTO)
}

func (h *Handler) CreateRepository(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	var req requests.CreateRepositoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	req.CreatedByID = authUser.ID
	req.OrganizationID = authUser.OrganizationID
	req.AddedByAdmin = true
	req.Status = "confirmed"

	ctx := c.Request.Context()
	repo, err := h.services.Repository.CreateRepository(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}

func (h *Handler) UpdateRepository(c *gin.Context) {
	repositoryId := c.Param("id")
	if repositoryId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_REPOSITORY_ID",
			"Repository ID is required",
			nil))
		return
	}
	var req requests.UpdateRepositoryRequest
	id, err := strconv.ParseUint(repositoryId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_REPOSITORY_ID",
			"Repository ID must be a valid number",
			err))
		return
	}
	req.ID = id
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()

	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	repositoryDTO, err := h.services.Repository.GetRepository(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessRepository(ctx, authUser, &repositoryDTO) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this repository",
			nil))
		return
	}

	repo, err := h.services.Repository.UpdateRepository(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}

func (h *Handler) DeleteRepository(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	fmt.Println(authUser.ID)

	repositoryId := c.Param("id")
	if repositoryId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_REPOSITORY_ID",
			"Repository ID is required",
			nil))
		return
	}

	id, err := strconv.ParseUint(repositoryId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_REPOSITORY_ID",
			"Repository ID must be a valid number",
			err))
		return
	}

	ctx := c.Request.Context()
	repositoryDTO, err := h.services.Repository.GetRepository(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessRepository(ctx, authUser, &repositoryDTO) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this repository",
			nil))
		return
	}

	err = h.services.Repository.DeleteRepository(ctx, repositoryId)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, true)
}

func (h *Handler) GetUserRepositories(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()
	query := c.Query("q")
	with := c.QueryArray("with")
	page, limit := pagination.Paginate(c)

	response, err := h.services.Repository.GetUserRepositories(ctx, authUser.OrganizationID, query, with, page, limit)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, response)
}

func (h *Handler) CreateUserRepository(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	var req requests.CreateUserRepositoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	req.CreatedByID = authUser.ID
	req.OrganizationID = authUser.OrganizationID
	req.AddedByAdmin = false
	req.Status = "pending"

	ctx := c.Request.Context()
	repo, err := h.services.Repository.CreateUserRepository(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}
