package repositories

import (
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/utils/pagination"
	"clusterix-code/internal/utils/preload"

	"gorm.io/gorm"

	"context"
)

type WorkspaceRepository struct {
	*Repository[models.Workspace]
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{
		Repository: NewRepository[models.Workspace](db),
	}
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id uint64) (*models.Workspace, error) {
	var workspace models.Workspace
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Repository").
		Preload("Repository.MachineConfig").
		Preload("WorkspaceConfig").
		First(&workspace, id).Error
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *WorkspaceRepository) GetByIDIncludingDeleted(ctx context.Context, id uint64) (*models.Workspace, error) {
	var workspace models.Workspace
	err := r.db.WithContext(ctx).
		Preload("Repository").
		Preload("Repository.MachineConfig").
		Preload("WorkspaceConfig").
		Preload("GitPersonalAccessToken").
		First(&workspace, id).Error
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *WorkspaceRepository) GetByFingerprint(ctx context.Context, fingerprint string) (*models.Workspace, error) {
	var workspace models.Workspace
	err := r.db.WithContext(ctx).
		Preload("WorkspaceConfig").
		Where("fingerprint = ?", fingerprint).
		First(&workspace).Error

	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (r *WorkspaceRepository) GetWorkspaces(ctx context.Context, userId uint64, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Workspace{}).
		Where("user_id = ?", userId)

	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.Workspace](query, page, limit)
}

func (r *WorkspaceRepository) Search(ctx context.Context, userId uint64, search string, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Workspace{}).
		Where("user_id = ?", userId).
		Where("title ILIKE ?", "%"+search+"%")

	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.Workspace](query, page, limit)
}

func (r *WorkspaceRepository) DeleteWorkspace(ctx context.Context, workspaceId string) error {
	if err := r.db.WithContext(ctx).Delete(&models.Workspace{}, workspaceId).Error; err != nil {
		return err
	}
	return nil
}

func (r *WorkspaceRepository) UpdateStatus(ctx context.Context, workspaceID uint64, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Workspace{}).
		Where("id = ?", workspaceID).
		Update("status", status).Error
}

func (r *WorkspaceRepository) UpdateURL(ctx context.Context, workspaceID uint64, url string) error {
	return r.db.WithContext(ctx).
		Model(&models.Workspace{}).
		Where("id = ?", workspaceID).
		Update("url", url).Error
}
