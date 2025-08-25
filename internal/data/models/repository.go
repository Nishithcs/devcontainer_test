package models

import (
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	ID              uint64 `gorm:"primaryKey"`
	Title           string `gorm:"type:varchar(255)"`
	MachineConfigID uint64 `gorm:"not null"`
	RepositoryURL   string `gorm:"type:text"`
	Status          string `gorm:"varchar(20);default:'pending'"` // confirmed, pending, ignored
	AddedByAdmin    bool   `gorm:"type:boolean"`
	CreatedByID     uint64 `gorm:"not null"`
	OrganizationID  uint32 `gorm:"not null"`

	User          User          `gorm:"foreignKey:CreatedByID"`
	MachineConfig MachineConfig `gorm:"foreignKey:MachineConfigID"`

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
}
