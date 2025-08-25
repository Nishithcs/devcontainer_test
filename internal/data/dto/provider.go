package dto

import "clusterix-code/internal/data/models"

type ProviderDTO struct {
	ID             uint64   `json:"id"`
	Title          string   `json:"title"`
	CreatedByID    *uint64  `json:"created_by_id,omitempty"`
	CreatedBy      *UserDto `json:"created_by,omitempty"`
	Icon           *string  `json:"icon"`
	WorkspaceCount int64    `json:"workspace_count"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

func ToProviderDTO(provider models.Provider, workspaceCounts map[uint64]int64) ProviderDTO {
	if workspaceCounts == nil {
		workspaceCounts = make(map[uint64]int64)
	}
	count := workspaceCounts[provider.ID]
	dto := ProviderDTO{
		ID:             provider.ID,
		Title:          provider.Title,
		CreatedByID:    provider.CreatedByID,
		Icon:           provider.Icon,
		WorkspaceCount: count,
		CreatedAt:      provider.CreatedAt.String(),
		UpdatedAt:      provider.UpdatedAt.String(),
	}

	if provider.User.ID != 0 {
		dto.CreatedBy = ToUserDTO(provider.User)
	}

	return dto
}

func ToProviderDTOs(providers []models.Provider, workspaceCounts map[uint64]int64) []ProviderDTO {
	result := make([]ProviderDTO, len(providers))
	for i, provider := range providers {
		result[i] = ToProviderDTO(provider, workspaceCounts)
	}
	return result
}
