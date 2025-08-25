package repositories

import (
	"clusterix-code/internal/utils/di"
	"context"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base repository for common CRUD operations
type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *Repository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *Repository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}

func (r *Repository[T]) RawQuery(ctx context.Context, query string, args ...interface{}) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Repository registry
type Repositories struct {
	User                   *UserRepository
	MachineConfig          *MachineConfigRepository
	Provider               *ProviderRepository
	GitPersonalAccessToken *GitPersonalAccessTokenRepository
	GitRepository          *GitRepository
	Workspace              *WorkspaceRepository
	WorkspaceConfig        *WorkspaceConfigRepository
	WorkspaceStatusEvent   *WorkspaceStatusEventRepository
	WorkspaceLog           *WorkspaceLogRepository
}

func Provider(c *di.Container) (*Repositories, error) {
	db := di.Make[*gorm.DB](c)
	mongoDB := di.Make[*mongo.Database](c)

	return NewRepositories(db, mongoDB), nil
}

func NewRepositories(db *gorm.DB, mongoDB *mongo.Database) *Repositories {
	return &Repositories{
		User:                   NewUserRepository(db),
		MachineConfig:          NewMachineConfigRepository(db),
		Provider:               NewProviderRepository(db),
		GitPersonalAccessToken: NewGitPersonalAccessTokenRepository(db),
		GitRepository:          NewGitRepository(db),
		Workspace:              NewWorkspaceRepository(db),
		WorkspaceConfig:        NewWorkspaceConfigRepository(db),
		WorkspaceStatusEvent:   NewWorkspaceStatusEventRepository(db),
		WorkspaceLog:           NewWorkspaceLogRepository(mongoDB),
	}
}
