package repositories

import (
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/utils/pagination"
	"context"

	"gorm.io/gorm"
)

type MachineConfigRepository struct {
	*Repository[models.MachineConfig]
}

func NewMachineConfigRepository(db *gorm.DB) *MachineConfigRepository {
	return &MachineConfigRepository{
		Repository: NewRepository[models.MachineConfig](db),
	}
}

func (r *MachineConfigRepository) GetByID(ctx context.Context, id uint64) (*models.MachineConfig, error) {
	var machineConfig models.MachineConfig
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&machineConfig).Error
	if err != nil {
		return nil, err
	}
	return &machineConfig, nil
}

func (r *MachineConfigRepository) GetByInstanceType(ctx context.Context, instanceType string) (*models.MachineConfig, error) {
	var config models.MachineConfig
	err := r.db.WithContext(ctx).Where("instance_type = ?", instanceType).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *MachineConfigRepository) CreateConfig(ctx context.Context, config *models.MachineConfig) error {
	return r.Create(ctx, config)
}

func (r *MachineConfigRepository) UpdateConfig(ctx context.Context, config *models.MachineConfig) error {
	return r.Update(ctx, config)
}

func (r *MachineConfigRepository) GetAll(ctx context.Context, page, limit int) (pagination.Pagination, error) {
	query := r.db.WithContext(ctx).Model(&models.MachineConfig{})
	return pagination.GormPaginate[models.MachineConfig](query, page, limit)
}

func (r *MachineConfigRepository) Search(ctx context.Context, query string, page, limit int) (pagination.Pagination, error) {
	db := r.db.WithContext(ctx).Model(&models.MachineConfig{})
	if query != "" {
		db = db.Where("instance_type ILIKE ?", "%"+query+"%")
	}

	return pagination.GormPaginate[models.MachineConfig](db, page, limit)
}
