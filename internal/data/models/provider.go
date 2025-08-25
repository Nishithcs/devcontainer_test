package models

import (
	"gorm.io/gorm"
	"time"
)

type Provider struct {
	ID          uint64 `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(255)"`
	CreatedByID *uint64
	Icon        *string `gorm:"type:text"`

	User User `gorm:"foreignKey:CreatedByID"`

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
}
