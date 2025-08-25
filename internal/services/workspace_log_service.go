// internal/services/workspace_log_service.go

package services

import (
	"clusterix-code/internal/api/requests"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/data/repositories"
	"context"
	"time"
)

type WorkspaceLogServiceConfig struct {
	Repositories *repositories.Repositories
}

type WorkspaceLogService struct {
	workspaceLogRepository *repositories.WorkspaceLogRepository
}

func NewWorkspaceLogService(config *WorkspaceLogServiceConfig) *WorkspaceLogService {
	return &WorkspaceLogService{
		workspaceLogRepository: config.Repositories.WorkspaceLog,
	}
}

func (s *WorkspaceLogService) GetWorkspaceLogs(ctx context.Context, workspaceId uint64) ([]*dto.WorkspaceLogDTO, error) {
	logs, err := s.workspaceLogRepository.GetLatestLogs(ctx, workspaceId)
	if err != nil {
		return nil, err
	}

	return dto.ToWorkspaceLogDTOs(logs), nil
}

func (s *WorkspaceLogService) Create(ctx context.Context, log requests.CreateWorkspaceLogRequest) error {
	workspaceLog := models.WorkspaceLog{
		WorkspaceID: log.WorkspaceID,
		Text:        log.Text,
		Type:        log.Type,
		Time:        log.Time,
		CreatedAt:   time.Now(),
	}
	return s.workspaceLogRepository.Create(ctx, &workspaceLog)
}
