package repositories

import (
	"clusterix-code/internal/data/enums"
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/utils/pagination"
	"clusterix-code/internal/utils/preload"
	"context"

	"gorm.io/gorm"
)

type ProviderRepository struct {
	*Repository[models.Provider]
}

func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{
		Repository: NewRepository[models.Provider](db),
	}
}

func (r *ProviderRepository) GetAll(ctx context.Context, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).Model(&models.Provider{})
	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.Provider](query, page, limit)
}

func (r *ProviderRepository) Search(ctx context.Context, search string, with []string, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).Model(&models.Provider{})

	if search != "" {
		query = query.Where("instance_type ILIKE ?", "%"+search+"%")
	}

	query = preload.ApplyPreloads(query, with)

	return pagination.GormPaginate[models.Provider](query, page, limit)
}

func (r *ProviderRepository) GetWorkspaceCounts(ctx context.Context) (map[uint64]int64, error) {
	type result struct {
		ProviderID uint64
		Count      int64
	}

	var results []result
	err := r.db.WithContext(ctx).
		Table("workspaces").
		Select("provider_id, COUNT(*) as count").
		Where("provider_id IS NOT NULL").
		Where("status != ?", enums.WorkspaceStatusTerminated).
		Group("provider_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[uint64]int64, len(results))
	for _, r := range results {
		counts[r.ProviderID] = r.Count
	}
	return counts, nil
}
