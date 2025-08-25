package models

import (
	"clusterix-code/internal/data/enums"
	"gorm.io/gorm"
	"time"
)

type Workspace struct {
	ID                       uint64                `gorm:"primaryKey"`
	Title                    string                `gorm:"type:varchar(255);not null"`
	Color                    string                `gorm:"type:varchar(30);not null"`
	Ide                      string                `gorm:"type:varchar(50);default:'vscode'"`
	RepositoryID             uint64                `gorm:"not null"`
	UserID                   uint64                `gorm:"not null"`
	URL                      string                `gorm:"type:text"`
	Fingerprint              string                `gorm:"uniqueIndex;type:varchar(100)"`
	OrganizationID           uint32                `gorm:"not null"`
	GitPersonalAccessTokenID uint64                `gorm:"not null"`
	Status                   enums.WorkspaceStatus `gorm:"type:varchar(50);not null"`

	Tags              []string `gorm:"type:text[]"`
	ProviderID        *uint64
	WorkspaceConfigID *uint64
	LastRunAt         *time.Time

	Repository             Repository             `gorm:"foreignKey:RepositoryID"`
	User                   User                   `gorm:"foreignKey:UserID"`
	GitPersonalAccessToken GitPersonalAccessToken `gorm:"foreignKey:GitPersonalAccessTokenID"`
	Provider               Provider               `gorm:"foreignKey:ProviderID"`
	WorkspaceConfig        WorkspaceConfig        `gorm:"foreignKey:WorkspaceConfigID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
