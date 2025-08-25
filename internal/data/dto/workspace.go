package dto

import (
	"clusterix-code/internal/data/models"
	"time"
)

type WorkspaceDTO struct {
	ID                uint64              `json:"id"`
	Title             string              `json:"title"`
	Color             string              `json:"color"`
	IDE               string              `json:"ide"`
	URL               string              `json:"url"`
	Fingerprint       string              `json:"fingerprint"`
	RepositoryID      uint64              `json:"repository_id"`
	Repository        *RepositoryDTO      `json:"repository,omitempty"` // Optional nested
	UserID            uint64              `json:"user_id"`
	User              *UserDto            `json:"user,omitempty"` // Optional nested
	GitAccessTokenID  uint64              `json:"git_access_token_id"`
	GitAccessToken    *GitAccessTokenDTO  `json:"git_access_token,omitempty"` // Optional nested
	OrganizationID    uint32              `json:"organization_id"`
	ProviderID        uint64              `json:"provider_id,omitempty"`
	Provider          *ProviderDTO        `json:"provider,omitempty"` // Optional nested
	WorkspaceConfig   *WorkspaceConfigDTO `json:"workspace_config,omitempty"`
	WorkspaceConfigID uint64              `json:"workspace_config_id"`
	Status            string              `json:"status"`
	Tags              []string            `json:"tags"`
	LastRunAt         *time.Time          `json:"last_run_at"`
	CreatedAt         string              `json:"created_at"`
	UpdatedAt         string              `json:"updated_at"`
}

func ToWorkspaceDTO(workspace models.Workspace) WorkspaceDTO {
	dto := WorkspaceDTO{
		ID:               workspace.ID,
		Title:            workspace.Title,
		Color:            workspace.Color,
		IDE:              workspace.Ide,
		URL:              workspace.URL,
		Fingerprint:      workspace.Fingerprint,
		RepositoryID:     workspace.RepositoryID,
		GitAccessTokenID: workspace.GitPersonalAccessTokenID,
		OrganizationID:   workspace.OrganizationID,
		Status:           string(workspace.Status),
		Tags:             workspace.Tags,
		LastRunAt:        workspace.LastRunAt,
		CreatedAt:        workspace.CreatedAt.String(),
		UpdatedAt:        workspace.UpdatedAt.String(),
		UserID:           workspace.UserID,
	}

	if workspace.ProviderID != nil {
		dto.ProviderID = *workspace.ProviderID
	}

	if workspace.User.ID != 0 {
		dto.User = ToUserDTO(workspace.User)
	}

	if workspace.Repository.ID != 0 {
		repo := ToRepositoryDTO(workspace.Repository)
		dto.Repository = &repo
	}

	if workspace.GitPersonalAccessToken.ID != 0 {
		token := ToGitAccessTokenDTO(workspace.GitPersonalAccessToken)
		dto.GitAccessToken = &token
	}

	if workspace.Provider.ID != 0 {
		provider := ToProviderDTO(workspace.Provider, nil)
		dto.Provider = &provider
	}

	if workspace.WorkspaceConfig.ID != 0 {
		workspaceConfig := ToWorkspaceConfigDTO(&workspace.WorkspaceConfig)
		dto.WorkspaceConfig = workspaceConfig
	}

	return dto
}

func ToWorkspaceDTOs(workspaces []models.Workspace) []WorkspaceDTO {
	result := make([]WorkspaceDTO, len(workspaces))
	for i, r := range workspaces {
		result[i] = ToWorkspaceDTO(r)
	}
	return result
}
