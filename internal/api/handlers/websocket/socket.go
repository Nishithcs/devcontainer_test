// clusterix-code/internal/api/handlers/socket/socket.go

package websocket

import (
	"clusterix-code/internal/services"
	"clusterix-code/internal/websocket"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) WebSocket(c *gin.Context) {
	websocket.ServeWs(h.services.Socket.Config.Hub, c.Writer, c.Request)
}
