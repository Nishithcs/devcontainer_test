package models

type MachineConfig struct {
	ID                 uint64  `gorm:"primaryKey"`
	InstanceType       string  `gorm:"type:varchar(255);unique;not null"`
	Category           string  `gorm:"type:varchar(255);not null"`
	CPUCores           int     `gorm:"not null"`
	MemoryGB           float64 `gorm:"not null"`
	StorageType        string  `gorm:"type:varchar(255)"`
	StorageSizeGB      int
	NetworkPerformance string `gorm:"type:varchar(255)"`
	Architecture       string `gorm:"type:varchar(255)"`
	Hypervisor         string `gorm:"type:varchar(255)"`
	Generation         string `gorm:"type:varchar(255)"`
	EnhancedNetworking bool   `gorm:"type:boolean"`
	GPU                string `gorm:"type:varchar(255)"`
	AdditionalFeatures string `gorm:"type:text"`
}

func (MachineConfig) TableName() string {
	return "machine_configs"
}
