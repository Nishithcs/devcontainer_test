package server

import (
	"clusterix-code/internal/config"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	cfg        *config.Config
}

// Provider initializes the server with dependencies.
func Provider(c *di.Container) (*Server, error) {
	cfg := di.Make[*config.Config](c)
	services := di.Make[*services.Services](c)

	router := NewRouter(services, cfg)

	server := &Server{
		router: router.engine,
		cfg:    cfg,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: router.engine,
		},
	}
	return server, nil
}

// Run starts the HTTP server and handles graceful shutdown.
func (s *Server) Run() {
	// Start the server in a goroutine
	go func() {
		logger.Info("Starting server",
			zap.Int("port", s.cfg.Server.Port),
			zap.String("environment", s.cfg.Server.Environment),
		)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Server failed unexpectedly", zap.Error(err))
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.Shutdown()
}

// Shutdown handles graceful shutdown of the server.
func (s *Server) Shutdown() {
	logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server shutdown completed successfully")
}
