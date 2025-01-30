package database

import (
	"event-registration/internal/config"
	"event-registration/internal/core/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.PostgresURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the Event model
	err = db.AutoMigrate(&domain.Event{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)           // Set maximum idle connections
	sqlDB.SetMaxOpenConns(100)          // Set maximum open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Set maximum connection lifetime

	return db, nil
}
