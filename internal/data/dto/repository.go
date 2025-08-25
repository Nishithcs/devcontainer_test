package dto

import (
	"clusterix-code/internal/data/models"
)

type RepositoryDTO struct {
	ID              uint64            `json:"id"`
	Title           string            `json:"title"`
	RepositoryURL   string            `json:"repository_url"`
	Status          string            `json:"status"`
	AddedByAdmin    bool              `json:"added_by_admin"`
	CreatedByID     uint64            `json:"created_by_id"`
	CreatedBy       *UserDto          `json:"created_by,omitempty"`
	MachineConfigID uint64            `json:"machine_config_id,omitempty"`
	MachineConfig   *MachineConfigDTO `json:"machine_config,omitempty"`
	OrganizationID  uint32            `json:"organization_id"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
}

func ToRepositoryDTO(repo models.Repository) RepositoryDTO {
	dto := RepositoryDTO{
		ID:              repo.ID,
		Title:           repo.Title,
		RepositoryURL:   repo.RepositoryURL,
		Status:          repo.Status,
		AddedByAdmin:    repo.AddedByAdmin,
		CreatedByID:     repo.CreatedByID,
		MachineConfigID: repo.MachineConfigID,
		OrganizationID:  repo.OrganizationID,
		CreatedAt:       repo.CreatedAt.String(),
		UpdatedAt:       repo.UpdatedAt.String(),
	}

	if repo.MachineConfig.ID != 0 {
		mc := ToMachineConfigDTO(repo.MachineConfig)
		dto.MachineConfig = &mc
	}

	if repo.User.ID != 0 {
		dto.CreatedBy = ToUserDTO(repo.User)
	}

	return dto
}

func ToRepositoryDTOs(repos []models.Repository) []RepositoryDTO {
	result := make([]RepositoryDTO, len(repos))
	for i, r := range repos {
		result[i] = ToRepositoryDTO(r)
	}
	return result
}
