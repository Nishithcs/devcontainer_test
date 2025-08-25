package provider

import (
	"clusterix-code/internal/api/handlers"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/pagination"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) GetProviders(c *gin.Context) {
	ctx := c.Request.Context()
	query := c.Query("q")
	with := c.QueryArray("with")
	page, limit := pagination.Paginate(c)

	response, err := h.services.Provider.GetProviders(ctx, query, with, page, limit)
	if err != nil {
		handlers.ErrorResponse(c, err)
		return
	}

	handlers.SuccessResponse(c, response)
}
