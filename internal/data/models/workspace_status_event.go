package models

import (
	"clusterix-code/internal/data/enums"
	"time"
)

type WorkspaceStatusEvent struct {
	ID          uint64                `gorm:"primaryKey"`
	WorkspaceID uint64                `gorm:"index;not null"`
	Status      enums.WorkspaceStatus `gorm:"type:varchar(50);not null"`
	Message     string                `gorm:"type:text"`
	CreatedAt   time.Time             `gorm:"autoCreateTime"`
}
