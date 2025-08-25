package dto

import (
	"clusterix-code/internal/data/models"
)

type GitAccessTokenDTO struct {
	ID        uint64   `json:"id"`
	Title     string   `json:"title"`
	UserID    uint64   `json:"user_id"`
	User      *UserDto `json:"user,omitempty"`
	IsDefault bool     `json:"is_default"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func ToGitAccessTokenDTO(gitAccessToken models.GitPersonalAccessToken) GitAccessTokenDTO {
	dto := GitAccessTokenDTO{
		ID:        gitAccessToken.ID,
		Title:     gitAccessToken.Title,
		UserID:    gitAccessToken.UserID,
		IsDefault: gitAccessToken.IsDefault,
		CreatedAt: gitAccessToken.CreatedAt.String(),
		UpdatedAt: gitAccessToken.UpdatedAt.String(),
	}

	if gitAccessToken.User.ID != 0 {
		dto.User = ToUserDTO(gitAccessToken.User)
	}

	return dto
}

func ToGitAccessTokenDTOs(tokens []models.GitPersonalAccessToken) []GitAccessTokenDTO {
	result := make([]GitAccessTokenDTO, len(tokens))
	for i, token := range tokens {
		result[i] = ToGitAccessTokenDTO(token)
	}
	return result
}
