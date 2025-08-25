package db

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"clusterix-code/internal/config"
	"clusterix-code/internal/utils/di"
	internalLogger "clusterix-code/internal/utils/logger"
)

func Provider(c *di.Container) (*gorm.DB, error) {
	cfg := di.Make[*config.Config](c)
	db, err := NewDatabase(cfg.Database)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		internalLogger.Fatal("Failed to get database instance", zap.Error(err))
	}
	if err := sqlDB.Ping(); err != nil {
		internalLogger.Fatal("Failed to ping database", zap.Error(err))
	}
	internalLogger.Info("Database connection established successfully")

	return db, nil
}

func NewDatabase(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.Name, config.Port, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(255);uniqueIndex"`
	AppliedAt time.Time
}

func (Migration) TableName() string {
	return "migrations"
}
