package requests

type CreateWorkspaceLogRequest struct {
	WorkspaceID uint64 `json:"workspace_id" binding:"required"`
	Text        string `json:"text" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Time        string `json:"time" binding:"required"`
}
