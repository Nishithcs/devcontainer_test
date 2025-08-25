package models

import (
	"gorm.io/gorm"
	"time"
)

type WorkspaceConfig struct {
	ID            uint64  `gorm:"primaryKey"`
	WorkspaceID   *uint64 `gorm:""`
	DevpodMachine *string `gorm:"type:varchar(20)"`
	WorkerName    *string `gorm:"type:varchar(20)"`
	WorkerIP      *string `gorm:"type:varchar(20)"`
	WorkerPort    *int    `gorm:"type:integer"`

	Workspace *Workspace `gorm:"foreignKey:WorkspaceID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
