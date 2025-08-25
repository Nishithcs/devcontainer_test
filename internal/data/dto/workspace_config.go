package dto

import (
	"clusterix-code/internal/data/models"
)

type WorkspaceConfigDTO struct {
	ID            uint64  `json:"id"`
	DevpodMachine *string `json:"devpod_machine,omitempty"`
	WorkerName    *string `json:"worker_name,omitempty"`
	WorkerIP      *string `json:"worker_ip,omitempty"`
	WorkerPort    *int    `json:"worker_port,omitempty"`
}

func ToWorkspaceConfigDTO(config *models.WorkspaceConfig) *WorkspaceConfigDTO {
	if config == nil {
		return nil
	}
	return &WorkspaceConfigDTO{
		ID:            config.ID,
		DevpodMachine: config.DevpodMachine,
		WorkerName:    config.WorkerName,
		WorkerIP:      config.WorkerIP,
		WorkerPort:    config.WorkerPort,
	}
}
