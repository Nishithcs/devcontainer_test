package services

import (
	"context"
	"os"
	"strings"

	"clusterix-code/internal/data/dto"
)

type PermissionService struct {
	superRoles map[string]struct{}
}

func NewPermissionService() *PermissionService {
	rolesEnv := os.Getenv("SUPER_ROLES")
	rolesEnv = strings.Trim(rolesEnv, "[]")
	roles := strings.Split(rolesEnv, ",")
	superRoles := make(map[string]struct{})
	for _, r := range roles {
		if r != "" {
			superRoles[r] = struct{}{}
		}
	}

	return &PermissionService{
		superRoles: superRoles,
	}
}

func (p *PermissionService) IsAdmin(userRoles []string) bool {
	for _, role := range userRoles {
		role = strings.TrimSpace(role)
		if _, ok := p.superRoles[role]; ok {
			return true
		}
	}
	return false
}

// Checks if the user has access to the workspace
func (p *PermissionService) CanAccessWorkspace(ctx context.Context, authUser *dto.User, workspace *dto.WorkspaceDTO) bool {
	if workspace.UserID != authUser.ID {
		return false
	}
	if workspace.OrganizationID != authUser.OrganizationID {
		return false
	}
	return true
}

// Checks if the user has access to the Repository
func (p *PermissionService) CanAccessRepository(ctx context.Context, authUser *dto.User, repo *dto.RepositoryDTO) bool {
    if repo.CreatedBy.ID != authUser.ID {
        return false
    }
    if repo.OrganizationID != authUser.OrganizationID {
        return false
    }
    return true
}

// Checks if the user has access to the GitToken
func (p *PermissionService) CanAccessGitToken(ctx context.Context, authUser *dto.User, token *dto.GitAccessTokenDTO) bool {
    if token.UserID != authUser.ID {
        return false
    }
    if uint32(token.User.OrganizationID) != authUser.OrganizationID {
        return false
    }
    return true
}
