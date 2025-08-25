package server

import (
	"clusterix-code/internal/api/handlers/auth"
	"clusterix-code/internal/api/handlers/git_access_token"
	"clusterix-code/internal/api/handlers/health"
	"clusterix-code/internal/api/handlers/machine_config"
	"clusterix-code/internal/api/handlers/metrics"
	"clusterix-code/internal/api/handlers/provider"
	"clusterix-code/internal/api/handlers/repository"
	"clusterix-code/internal/api/handlers/websocket"
	"clusterix-code/internal/api/handlers/workspace"
	"clusterix-code/internal/api/handlers/workspace_log"
	"clusterix-code/internal/api/middleware"
	"clusterix-code/internal/config"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/logger"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine   *gin.Engine
	services *services.Services
}

// NewRouter initializes and returns a new Router instance.
func NewRouter(services *services.Services, cfg *config.Config) *Router {
	engine := InitRouter(cfg)

	router := &Router{
		engine:   engine,
		services: services,
	}

	router.SetupRoutes()

	return router
}

// InitRouter initializes the Gin router with middleware.
func InitRouter(cfg *config.Config) *gin.Engine {
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Apply middleware
	router.Use(
		gin.Recovery(),
		metrics.Middleware(),
		middleware.Logger(),
		middleware.RequestID(),
		middleware.CORS(),
		middleware.ErrorHandler(),
	)

	logger.Info("Router initialized successfully")
	return router
}

// SetupRoutes defines all application routes.
func (r *Router) SetupRoutes() {
	healthHandler := health.NewHandler()
	machineConfigHandler := machine_config.NewHandler(r.services)
	providerHandler := provider.NewHandler(r.services)
	gitAccessTokenHandler := git_access_token.NewHandler(r.services)
	repositoryHandler := repository.NewHandler(r.services)
	workspaceHandler := workspace.NewHandler(r.services)
	authHandler := auth.NewHandler(r.services)
	socketHandler := websocket.NewHandler(r.services)
	workspaceLogHandler := workspace_log.NewHandler(r.services)

	// Metrics and Health Check Endpoints
	r.engine.GET("/metrics", metrics.Handler())
	r.engine.GET("/health", healthHandler.Health)

	// User-related routes
	protected := r.engine.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/auth", authHandler.GenerateShortAuthToken)

		protected.GET("/machine-configs", machineConfigHandler.GetMachineConfigs)
		protected.GET("/machine-configs/:id", machineConfigHandler.GetMachineConfig)
		protected.GET("/providers", providerHandler.GetProviders)

		protected.GET("/git-access-tokens", gitAccessTokenHandler.GetUserAccessTokens)
		protected.GET("/git-access-tokens/:id", gitAccessTokenHandler.GetUserAccessToken)
		protected.POST("/git-access-tokens", gitAccessTokenHandler.CreateUserAccessToken)
		protected.PATCH("/git-access-tokens/:id", gitAccessTokenHandler.UpdateUserAccessToken)
		protected.DELETE("/git-access-tokens/:id", gitAccessTokenHandler.DeleteUserAccessToken)

		protected.GET("/repositories", middleware.AdminOnly(), repositoryHandler.GetRepositories)
		protected.GET("/repositories/:id", middleware.AdminOnly(), repositoryHandler.GetRepository)
		protected.POST("/repositories", middleware.AdminOnly(), repositoryHandler.CreateRepository)
		protected.PATCH("/repositories/:id", middleware.AdminOnly(), repositoryHandler.UpdateRepository)
		protected.DELETE("/repositories/:id", middleware.AdminOnly(), repositoryHandler.DeleteRepository)
		protected.GET("/user-repositories", repositoryHandler.GetUserRepositories)
		protected.POST("/user-repositories", repositoryHandler.CreateUserRepository)

		protected.GET("/workspaces", workspaceHandler.GetWorkspaces)
		protected.GET("/workspaces/:id", workspaceHandler.GetWorkspace)
		protected.POST("/workspaces", workspaceHandler.CreateWorkspace)
		protected.PATCH("/workspaces/:id", workspaceHandler.UpdateWorkspace)
		protected.DELETE("/workspaces/:id", workspaceHandler.DeleteWorkspace)
		protected.GET("/workspaces/fingerprint/:fingerprint", workspaceHandler.GetWorkspaceByFingerprint)
		protected.GET("/workspaces/:id/logs", workspaceLogHandler.GetWorkspaceLogs)

		protected.POST("/workspaces/:id/start", workspaceHandler.StartWorkspace)
		protected.POST("/workspaces/:id/stop", workspaceHandler.StopWorkspace)
		protected.POST("/workspaces/:id/restart", workspaceHandler.RestartWorkspace)
		protected.POST("/workspaces/:id/rebuild", workspaceHandler.RebuildWorkspace)
		protected.POST("/workspaces/:id/terminate", workspaceHandler.TerminateWorkspace)
	}

	// Websocket
	socket := r.engine.Group("")
	socket.Use(
		middleware.AuthMiddlewareWithQueryParam(),
	)
	{
		socket.GET("/ws", socketHandler.WebSocket)
	}
}
