package repositories

import (
	"clusterix-code/internal/data/models"
	"context"
	"gorm.io/gorm"
)

type WorkspaceConfigRepository struct {
	*Repository[models.WorkspaceConfig]
}

func NewWorkspaceConfigRepository(db *gorm.DB) *WorkspaceConfigRepository {
	return &WorkspaceConfigRepository{
		Repository: NewRepository[models.WorkspaceConfig](db),
	}
}

func (r *WorkspaceConfigRepository) GetByID(ctx context.Context, id uint64) (*models.WorkspaceConfig, error) {
	var workspaceConfig models.WorkspaceConfig
	err := r.db.WithContext(ctx).First(&workspaceConfig, id).Error
	if err != nil {
		return nil, err
	}
	return &workspaceConfig, nil
}
