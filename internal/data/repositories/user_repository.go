package repositories

import (
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/data/models"
	"context"
	"gorm.io/gorm/clause"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	*Repository[models.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository: NewRepository[models.User](db),
	}
}

func (r *UserRepository) SyncUser(ctx context.Context, userDto *dto.AuthUserDto) (*models.User, error) {
	user := models.User{
		ID:             userDto.ID,
		FirstName:      userDto.Profile.FirstName,
		LastName:       userDto.Profile.LastName,
		FullName:       userDto.Name,
		Email:          userDto.Email,
		OrganizationID: userDto.OrganizationID,
		Avatar:         userDto.Profile.Avatar,
		IsActive:       userDto.IsActive == 1,
		DeletedAt: func() gorm.DeletedAt {
			if userDto.IsActive == 0 {
				return gorm.DeletedAt{Time: time.Now(), Valid: true}
			}
			return gorm.DeletedAt{}
		}(),
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "full_name", "email", "organization_id", "avatar", "is_active", "updated_at"}),
		}).
		Create(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SoftDelete(ctx context.Context, userID uint64) error {
	return r.db.
		WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"is_active":  false,
		}).Error
}
