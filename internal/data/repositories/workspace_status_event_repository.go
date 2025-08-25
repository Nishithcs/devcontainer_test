package repositories

import (
	"clusterix-code/internal/data/models"
	"context"

	"gorm.io/gorm"
)

type WorkspaceStatusEventRepository struct {
	db *gorm.DB
}

func NewWorkspaceStatusEventRepository(db *gorm.DB) *WorkspaceStatusEventRepository {
	return &WorkspaceStatusEventRepository{db: db}
}

func (r *WorkspaceStatusEventRepository) Create(ctx context.Context, event *models.WorkspaceStatusEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}
