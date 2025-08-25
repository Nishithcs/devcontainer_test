package workspace_log

import (
	"clusterix-code/internal/api/api_context"
	"clusterix-code/internal/api/handlers"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/errors"
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

func (h *Handler) GetWorkspaceLogs(c *gin.Context) {
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

	res, err := h.services.WorkspaceLog.GetWorkspaceLogs(ctx, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, res)
}
