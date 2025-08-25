package health

import (
	"clusterix-code/internal/api/handlers"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/utils/helpers"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

var startTime = time.Now()

func (h *Handler) Health(c *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	systemInfo := dto.SystemInfo{
		Version:      "1.0.0", // You can make this configurable
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		MemoryUsage:  helpers.FormatBytes(mem.Alloc),
		NumCPU:       runtime.NumCPU(),
		Uptime:       helpers.FormatUptime(time.Since(startTime)),
	}

	response := dto.HealthResponse{
		Status:      "OK",
		Timestamp:   time.Now().Unix(),
		Environment: helpers.GetEnvironment(),
		System:      systemInfo,
	}

	handlers.SuccessResponse(c, response)
}
