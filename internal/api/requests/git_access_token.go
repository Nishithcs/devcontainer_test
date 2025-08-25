package requests

type CreateGitAccessTokenRequest struct {
	Title     string `json:"title" binding:"required"`
	Token     string `json:"token" binding:"required"`
	UserID    uint64 `json:"user_id" binding:"required"`
	IsDefault *bool  `json:"is_default" binding:"required"`
}

type UpdateGitAccessTokenRequest struct {
	ID        uint64 `json:"id" binding:"required"`
	Title     string `json:"title" binding:"omitempty"`
	IsDefault *bool  `json:"is_default" binding:"omitempty"`
}
