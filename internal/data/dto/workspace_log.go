package dto

import (
	"clusterix-code/internal/data/models"
)

type WorkspaceLogDTO struct {
	Text string `json:"text"`
	Type string `json:"type"`
	Time string `json:"time"`
}

func ToWorkspaceLogDTO(config *models.WorkspaceLog) *WorkspaceLogDTO {
	if config == nil {
		return nil
	}
	return &WorkspaceLogDTO{
		Text: config.Text,
		Type: config.Type,
		Time: config.Time,
	}
}

func ToWorkspaceLogDTOs(logs []*models.WorkspaceLog) []*WorkspaceLogDTO {
	dtos := make([]*WorkspaceLogDTO, 0, len(logs))
	for _, log := range logs {
		dtos = append(dtos, ToWorkspaceLogDTO(log))
	}
	return dtos
}
