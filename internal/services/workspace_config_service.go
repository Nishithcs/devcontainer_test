package services

import (
	"clusterix-code/internal/api/requests"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/data/repositories"
	"context"
)

type WorkspaceConfigServiceConfig struct {
	Repositories *repositories.Repositories
}

type WorkspaceConfigService struct {
	workspaceConfigRepository *repositories.WorkspaceConfigRepository
}

func NewWorkspaceConfigService(config *WorkspaceConfigServiceConfig) *WorkspaceConfigService {
	return &WorkspaceConfigService{
		workspaceConfigRepository: config.Repositories.WorkspaceConfig,
	}
}

func (s *WorkspaceConfigService) CreateWorkspaceConfig(ctx context.Context, req requests.CreateWorkspaceConfigRequest) (dto.WorkspaceConfigDTO, error) {
	workspaceConfig := models.WorkspaceConfig{
		DevpodMachine: &req.DevpodMachine,
		WorkerName:    &req.WorkerName,
		WorkerIP:      &req.WorkerIP,
		WorkerPort:    &req.WorkerPort,
	}

	if err := s.workspaceConfigRepository.Create(ctx, &workspaceConfig); err != nil {
		return dto.WorkspaceConfigDTO{}, err
	}

	createdWorkspaceConfig, _ := s.workspaceConfigRepository.GetByID(ctx, workspaceConfig.ID)

	return *dto.ToWorkspaceConfigDTO(createdWorkspaceConfig), nil
}

func (s *WorkspaceConfigService) UpdateWorkspaceConfig(ctx context.Context, req requests.UpdateWorkspaceConfigRequest) (dto.WorkspaceConfigDTO, error) {
	workspaceConfig, err := s.workspaceConfigRepository.GetByID(ctx, req.ID)
	if err != nil {
		return dto.WorkspaceConfigDTO{}, err
	}

	if req.DevpodMachine != "" {
		workspaceConfig.DevpodMachine = &req.DevpodMachine
	}
	if req.WorkerName != "" {
		workspaceConfig.WorkerName = &req.WorkerName
	}
	if req.WorkerIP != "" {
		workspaceConfig.WorkerIP = &req.WorkerIP
	}
	if req.WorkerPort != 0 {
		workspaceConfig.WorkerPort = &req.WorkerPort
	}

	if err := s.workspaceConfigRepository.Update(ctx, workspaceConfig); err != nil {
		return dto.WorkspaceConfigDTO{}, err
	}
	return *dto.ToWorkspaceConfigDTO(workspaceConfig), nil
}
