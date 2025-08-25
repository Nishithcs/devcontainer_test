package workspace

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

func (h *Handler) GetWorkspaces(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()
	query := c.Query("q")
	with := c.QueryArray("with")
	page, limit := pagination.Paginate(c)

	response, err := h.services.Workspace.GetWorkspaces(ctx, authUser.ID, query, with, page, limit)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, response)
}

func (h *Handler) GetWorkspace(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	fmt.Println(authUser)

	workspaceId := c.Param("id")
	if workspaceId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_WORKSPACE_ID",
			"Workspace ID is required",
			nil))
		return
	}
	id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_WORKSPACE_ID",
			"Workspace ID must be a valid number",
			err))
		return
	}

	ctx := c.Request.Context()
	workspace, err := h.services.Workspace.GetWorkspace(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessWorkspace(ctx, authUser, &workspace) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this workspace",
			nil))
		return
	}

	handlers.SuccessResponse(c, workspace)
}

func (h *Handler) GetWorkspaceByFingerprint(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	fmt.Println(authUser)

	workspaceFingerprint := c.Param("fingerprint")
	if workspaceFingerprint == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_WORKSPACE_FINGERPRINT",
			"Workspace Fingerprint is required",
			nil))
		return
	}

	ctx := c.Request.Context()
	workspace, err := h.services.Workspace.GetWorkspaceByFingerprint(ctx, workspaceFingerprint)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessWorkspace(ctx, authUser, &workspace) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this workspace",
			nil))
		return
	}


	handlers.SuccessResponse(c, workspace)
}

func (h *Handler) CreateWorkspace(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	var req requests.CreateWorkspaceRequest
	req.UserID = authUser.ID
	req.OrganizationID = authUser.OrganizationID
	req.Status = "processing"
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()
	repo, err := h.services.Workspace.CreateWorkspace(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}

func (h *Handler) UpdateWorkspace(c *gin.Context) {
	workspaceId := c.Param("id")
	if workspaceId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_WORKSPACE_ID",
			"Workspace ID is required",
			nil))
		return
	}
	var req requests.UpdateWorkspaceRequest
	id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_WORKSPACE_ID",
			"Workspace ID must be a valid number",
			err))
		return
	}
	req.ID = id
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()
	workspace, err := h.services.Workspace.GetWorkspace(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessWorkspace(ctx, authUser, &workspace) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this workspace",
			nil))
		return
	}

	repo, err := h.services.Workspace.UpdateWorkspace(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}

func (h *Handler) DeleteWorkspace(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	fmt.Println(authUser.ID)

	workspaceId := c.Param("id")
	if workspaceId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_REPOSITORY_ID",
			"Repository ID is required",
			nil))
		return
	}

	ctx := c.Request.Context()

	id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_WORKSPACE_ID",
			"Workspace ID must be a valid number",
			err))
		return
	}
	workspace, err := h.services.Workspace.GetWorkspace(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessWorkspace(ctx, authUser, &workspace) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this workspace",
			nil))
		return
	}

	err = h.services.Workspace.DeleteWorkspace(ctx, workspaceId)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, true)
}

func (h *Handler) StartWorkspace(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	workspaceId := c.Param("id")
	if workspaceId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_WORKSPACE_ID",
			"Workspace ID is required",
			nil))
		return
	}
	var req requests.WorkspaceActionRequest
	id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_WORKSPACE_ID",
			"Workspace ID must be a valid number",
			err))
		return
	}
	req.ID = id
	req.UserID = authUser.ID

	ctx := c.Request.Context()
	workspace, err := h.services.Workspace.GetWorkspace(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessWorkspace(ctx, authUser, &workspace) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this workspace",
			nil))
		return
	}

	repo, err := h.services.Workspace.StartWorkspace(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}

func (h *Handler) StopWorkspace(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	workspaceId := c.Param("id")
	if workspaceId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_WORKSPACE_ID",
			"Workspace ID is required",
			nil))
		return
	}
	var req requests.WorkspaceActionRequest
	id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_WORKSPACE_ID",
			"Workspace ID must be a valid number",
			err))
		return
	}
	req.ID = id
	req.UserID = authUser.ID

	ctx := c.Request.Context()
	workspace, err := h.services.Workspace.GetWorkspace(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessWorkspace(ctx, authUser, &workspace) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this workspace",
			nil))
		return
	}

	repo, err := h.services.Workspace.StopWorkspace(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}

func (h *Handler) RestartWorkspace(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	workspaceId := c.Param("id")
	if workspaceId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_WORKSPACE_ID",
			"Workspace ID is required",
			nil))
		return
	}
	var req requests.WorkspaceActionRequest
	id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_WORKSPACE_ID",
			"Workspace ID must be a valid number",
			err))
		return
	}
	req.ID = id
	req.UserID = authUser.ID

	ctx := c.Request.Context()
	workspace, err := h.services.Workspace.GetWorkspace(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessWorkspace(ctx, authUser, &workspace) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this workspace",
			nil))
		return
	}

	repo, err := h.services.Workspace.RestartWorkspace(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}

func (h *Handler) RebuildWorkspace(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	workspaceId := c.Param("id")
	if workspaceId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_WORKSPACE_ID",
			"Workspace ID is required",
			nil))
		return
	}
	var req requests.WorkspaceActionRequest
	id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_WORKSPACE_ID",
			"Workspace ID must be a valid number",
			err))
		return
	}
	req.ID = id
	req.UserID = authUser.ID

	ctx := c.Request.Context()
	workspace, err := h.services.Workspace.GetWorkspace(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessWorkspace(ctx, authUser, &workspace) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this workspace",
			nil))
		return
	}

	repo, err := h.services.Workspace.RebuildWorkspace(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}

func (h *Handler) TerminateWorkspace(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	workspaceId := c.Param("id")
	if workspaceId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_WORKSPACE_ID",
			"Workspace ID is required",
			nil))
		return
	}
	var req requests.WorkspaceActionRequest
	id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_WORKSPACE_ID",
			"Workspace ID must be a valid number",
			err))
		return
	}
	req.ID = id
	req.UserID = authUser.ID

	ctx := c.Request.Context()
	workspace, err := h.services.Workspace.GetWorkspace(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	// Validate user permission
	permissionService := services.NewPermissionService()
	if !permissionService.CanAccessWorkspace(ctx, authUser, &workspace) {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeAuth,
			"FORBIDDEN",
			"You do not have access to this workspace",
			nil))
		return
	}

	repo, err := h.services.Workspace.TerminateWorkspace(ctx, req)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, repo)
}
