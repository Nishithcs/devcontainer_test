package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             uint64  `gorm:"primaryKey"`
	FirstName      string  `gorm:"type:varchar(255)"`
	LastName       string  `gorm:"type:varchar(255)"`
	FullName       string  `gorm:"type:varchar(255);index"`
	Email          string  `gorm:"type:varchar(255)"`
	Avatar         *string `gorm:"type:text"`
	OrganizationID uint64  `gorm:"type:int"`
	IsActive       bool    `gorm:"type:boolean"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}
