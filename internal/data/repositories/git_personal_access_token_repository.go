package repositories

import (
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/utils/pagination"
	"clusterix-code/internal/utils/preload"
	"context"
	"gorm.io/gorm"
)

type GitPersonalAccessTokenRepository struct {
	*Repository[models.GitPersonalAccessToken]
}

func NewGitPersonalAccessTokenRepository(db *gorm.DB) *GitPersonalAccessTokenRepository {
	return &GitPersonalAccessTokenRepository{
		Repository: NewRepository[models.GitPersonalAccessToken](db),
	}
}

func (r *GitPersonalAccessTokenRepository) GetByID(ctx context.Context, userId uint64, id uint64) (*models.GitPersonalAccessToken, error) {
	var repo models.GitPersonalAccessToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND user_id = ?", id, userId).
		First(&repo).Error
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *GitPersonalAccessTokenRepository) Search(ctx context.Context, userId uint64, search string, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).
		Model(&models.GitPersonalAccessToken{}).
		Where("user_id = ?", userId).
		Where("title ILIKE ?", "%"+search+"%")

	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.GitPersonalAccessToken](query, page, limit)
}

func (r *GitPersonalAccessTokenRepository) GetUserAccessTokens(ctx context.Context, userId uint64, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).
		Model(&models.GitPersonalAccessToken{}).
		Where("user_id = ?", userId)

	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.GitPersonalAccessToken](query, page, limit)
}

func (r *GitPersonalAccessTokenRepository) DeleteUserAccessTokens(ctx context.Context, gitAccessTokenId string) error {
	if err := r.db.WithContext(ctx).Delete(&models.GitPersonalAccessToken{}, gitAccessTokenId).Error; err != nil {
		return err
	}
	return nil
}
