package machine_config

import (
	"clusterix-code/internal/api/api_context"
	"clusterix-code/internal/api/handlers"
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

func (h *Handler) GetMachineConfigs(c *gin.Context) {
	ctx := c.Request.Context()
	query := c.Query("q")
	page, limit := pagination.Paginate(c)

	//authUser, err := api_context.AuthUser(c)
	//if err != nil {
	//	handlers.ErrorResponse(c, errors.NewAuthenticationError(err.Error()))
	//	return
	//}

	response, err := h.services.MachineConfig.GetMachineConfigs(ctx, query, page, limit)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	handlers.SuccessResponse(c, response)
}

func (h *Handler) GetMachineConfig(c *gin.Context) {
	authUser, err := api_context.AuthUser(c)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	fmt.Println(authUser)

	machineConfigId := c.Param("id")
	if machineConfigId == "" {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"MISSING_MACHINE_CONFIG_ID",
			"Machine Config ID is required",
			nil))
		return
	}
	id, err := strconv.ParseUint(machineConfigId, 10, 64)
	if err != nil {
		handlers.ErrorResponse(c, errors.NewError(
			errors.ErrorTypeBadRequest,
			"INVALID_REPOSITORY_ID",
			"Machine Config ID must be a valid number",
			err))
		return
	}

	ctx := c.Request.Context()
	res, err := h.services.MachineConfig.GetMachineConfig(ctx, authUser.ID, id)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}
	handlers.SuccessResponse(c, res)
}
