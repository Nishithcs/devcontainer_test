package requests

type CreateWorkspaceRequest struct {
	Title            string   `json:"title" binding:"required"`
	Color            string   `json:"color" binding:"required"`
	IDE              string   `json:"ide" binding:"required"`
	RepositoryID     uint64   `json:"repository_id" binding:"required"`
	UserID           uint64   `json:"user_id" binding:"required"`
	GitAccessTokenID uint64   `json:"git_access_token_id" binding:"required"`
	OrganizationID   uint32   `json:"organization_id" binding:"required"`
	ProviderID       uint64   `json:"provider_id" binding:"required"`
	Status           string   `json:"status" binding:"required,oneof=running processing failed"`
	Tags             []string `json:"tags"`
}

type UpdateWorkspaceRequest struct {
	ID               uint64   `json:"id" binding:"required"`
	Title            string   `json:"title" binding:"omitempty"`
	Color            string   `json:"color" binding:"omitempty"`
	IDE              string   `json:"ide" binding:"omitempty"`
	RepositoryID     uint64   `json:"repository_id" binding:"omitempty"`
	GitAccessTokenID uint64   `json:"git_access_token_id" binding:"omitempty"`
	ProviderID       uint64   `json:"provider_id" binding:"omitempty"`
	Tags             []string `json:"tags"`
}

type WorkspaceActionRequest struct {
	ID     uint64 `json:"id" binding:"required"`
	UserID uint64 `json:"user_id" binding:"required"`
}
