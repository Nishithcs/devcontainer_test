package models

import (
	"time"
)

type WorkspaceLog struct {
	WorkspaceID uint64    `gorm:"index;not null"`
	Text        string    `gorm:"type:text"`
	Type        string    `gorm:"type:type"`
	Time        string    `gorm:"type:time"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
