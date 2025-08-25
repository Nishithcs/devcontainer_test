package repositories

import (
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/utils/pagination"
	"clusterix-code/internal/utils/preload"

	"context"
	"gorm.io/gorm"
)

type GitRepository struct {
	*Repository[models.Repository]
}

func NewGitRepository(db *gorm.DB) *GitRepository {
	return &GitRepository{
		Repository: NewRepository[models.Repository](db),
	}
}

func (r *GitRepository) GetAll(ctx context.Context) ([]models.Repository, error) {
	var repos []models.Repository
	err := r.db.WithContext(ctx).Preload("User").Preload("MachineConfig").Find(&repos).Error
	return repos, err
}

func (r *GitRepository) GetByID(ctx context.Context, id uint64, with []string) (*models.Repository, error) {
	var repo models.Repository

	query := r.db.WithContext(ctx).Model(&models.Repository{})

	query = preload.ApplyPreloads(query, with)

	err := query.First(&repo, id).Error
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *GitRepository) GetRepositories(ctx context.Context, organizationId uint32, status string, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Repository{}).
		Where("organization_id = ?", organizationId)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.Repository](query, page, limit)
}

func (r *GitRepository) Search(ctx context.Context, organizationId uint32, search, status string, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Repository{}).
		Preload("User").
		Preload("MachineConfig").
		Where("organization_id = ?", organizationId).
		Where("title ILIKE ? OR repository_url ILIKE ?", "%"+search+"%", "%"+search+"%")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.Repository](query, page, limit)
}

func (r *GitRepository) GetUserRepositories(ctx context.Context, organizationId uint32, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Repository{}).
		Where("organization_id = ? AND status = ?", organizationId, "confirmed")

	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.Repository](query, page, limit)
}

func (r *GitRepository) SearchUserRepositories(ctx context.Context, organizationId uint32, search string, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Repository{}).
		Where("organization_id = ? AND status = ?", organizationId, "confirmed").
		Where("title ILIKE ? OR repository_url ILIKE ?", "%"+search+"%", "%"+search+"%")

	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.Repository](query, page, limit)
}

func (r *GitRepository) DeleteRepository(ctx context.Context, repositoryId string) error {
	if err := r.db.WithContext(ctx).Delete(&models.Repository{}, repositoryId).Error; err != nil {
		return err
	}
	return nil
}
