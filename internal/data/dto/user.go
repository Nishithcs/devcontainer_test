package dto

import "clusterix-code/internal/data/models"

type UserDto struct {
	ID             uint64 `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	FullName       string `json:"full_name"`
	Avatar         string `json:"avatar"`
	OrganizationID uint64 `json:"organization_id"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type AuthUserDto struct {
	ID             uint64           `json:"id"`
	Email          string           `json:"email"`
	Name           string           `json:"name"`
	Username       string           `json:"username"`
	OrganizationID uint64           `json:"organization_id"`
	IsActive       int              `json:"is_active"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
	Profile        *AuthUserProfile `json:"profile,omitempty"`
}

type AuthUserProfile struct {
	FirstName string  `json:"f_name"`
	LastName  string  `json:"l_name"`
	Avatar    *string `json:"avatar"`
	Language  string  `json:"language"`
}

type AnonymousUserDto struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	OrganizationID uint64 `json:"organization_id"`
	AccessToken    string `json:"access_token"`
	ExpiresIn      uint64 `json:"expires_in"`
	TokenType      string `json:"token_type"`
}

type AuthApplicationDto struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

func ToUserDTO(user models.User) *UserDto {
	var avatar string
	if user.Avatar != nil {
		avatar = *user.Avatar
	}
	return &UserDto{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		FullName:       user.FullName,
		Avatar:         avatar,
		OrganizationID: user.OrganizationID,
		IsActive:       user.IsActive,
		CreatedAt:      user.CreatedAt.String(),
		UpdatedAt:      user.UpdatedAt.String(),
	}
}
