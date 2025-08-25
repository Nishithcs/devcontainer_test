package models

import (
	"gorm.io/gorm"
	"time"
)

type GitPersonalAccessToken struct {
	ID        uint64 `gorm:"primaryKey"`
	Title     string `gorm:"type:varchar(255);not null"`
	Token     string `gorm:"type:varchar(255);not null"`
	UserID    uint64 `gorm:"not null"`
	IsDefault bool   `gorm:"type:boolean"`

	User User `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (GitPersonalAccessToken) TableName() string {
	return "git_personal_access_tokens"
}
