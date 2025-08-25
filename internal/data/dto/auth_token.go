package dto

import (
	"time"
)

// User struct to decode and encode user information
type User struct {
	ID             uint64      `json:"id"`
	Name           string      `json:"name"`
	Email          string      `json:"email"`
	IsActive       int         `json:"is_active"`
	Username       string      `json:"username"`
	OrganizationID uint32      `json:"organization_id"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DeletedAt      interface{} `json:"deleted_at"`
}

type Employee struct {
	Id           int32   `json:"id"`
	DepartmentId int32   `json:"department_id"`
	LocationId   int32   `json:"location_id"`
	TeamsIds     []int32 `json:"teams_ids"`
}

// TokenPayload struct to decode and encode the entire JSON payload
type TokenPayload struct {
	Iss          string   `json:"iss"`
	Iat          int64    `json:"iat"`
	Exp          int64    `json:"exp"`
	Nbf          int64    `json:"nbf"`
	Jti          string   `json:"jti"`
	Sub          int64    `json:"sub"`
	Prv          string   `json:"prv"`
	UserAgent    string   `json:"user_agent"`
	User         User     `json:"user"`
	Permissions  []any    `json:"permissions"`
	Roles        []string `json:"roles"`
	Employee     Employee `json:"employee"`
	PartnerNodes any      `json:"partner_nodes"` // will be a map[string]interface{} for registered users and just an empty array for anonymous users
}

func (t *TokenPayload) GetUserID() uint64 {
	return t.User.ID
}

func (t *TokenPayload) GetOrganizationID() uint32 {
	return t.User.OrganizationID
}
