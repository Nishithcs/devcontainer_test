package requests

type CreateRepositoryRequest struct {
	Title           string `json:"title" binding:"required"`
	MachineConfigID uint64 `json:"machine_config_id" binding:"required"`
	RepositoryURL   string `json:"repository_url" binding:"required"`
	OrganizationID  uint32 `json:"organization_id"`
	AddedByAdmin    bool   `json:"added_by_admin"`
	CreatedByID     uint64 `json:"created_by_id"`
	Status          string `json:"status" binding:"omitempty,oneof=confirmed pending ignored"`
}

type UpdateRepositoryRequest struct {
	ID              uint64 `json:"id" binding:"required"`
	Title           string `json:"title" binding:"omitempty"`
	MachineConfigID uint64 `json:"machine_config_id" binding:"omitempty"`
	RepositoryURL   string `json:"repository_url" binding:"omitempty"`
	Status          string `json:"status" binding:"omitempty,oneof=confirmed pending ignored"`
}

type CreateUserRepositoryRequest struct {
	Title          string `json:"title" binding:"required"`
	RepositoryURL  string `json:"repository_url" binding:"required"`
	OrganizationID uint32 `json:"organization_id"`
	AddedByAdmin   bool   `json:"added_by_admin"`
	CreatedByID    uint64 `json:"created_by_id"`
	Status         string `json:"status" binding:"omitempty,oneof=confirmed pending ignored"`
}
