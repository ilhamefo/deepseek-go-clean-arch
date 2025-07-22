package database

import (
	"event-registration/internal/common"
	"event-registration/internal/config"
	"event-registration/internal/core/domain"
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg *common.Config) (*gorm.DB, error) {
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

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func NewGormPlnMobileDB(cfg *common.Config, loggr *config.ZapLogger) (*gorm.DB, error) {

	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresPlnMobileUser,
		url.QueryEscape(cfg.PostgresPlnMobilePassword), // Use the encoded password
		cfg.PostgresPlnMobileHost,
		cfg.PostgresPlnMobilePort,
		cfg.PostgresPlnMobileDatabase,
	)

	db, err := gorm.Open(postgres.Open(connURL), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db.Logger = loggr

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if !cfg.IsProduction {
		return db.Debug(), nil
	}

	return db, nil
}

func NewGormDBAuth(cfg *common.Config, loggr *config.ZapLogger) (*gorm.DB, error) {

	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.AuthDBUser,
		url.QueryEscape(cfg.AuthDBPassword), // Use the encoded password
		cfg.AuthDBHost,
		cfg.AuthDBPort,
		cfg.AuthDB,
	)

	db, err := gorm.Open(postgres.Open(connURL), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db.Logger = loggr

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if !cfg.IsProduction {
		return db.Debug(), nil
	}

	return db, nil
}

func NewGarminGormDB(cfg *common.Config, loggr *config.ZapLogger) (*gorm.DB, error) {

	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.GarminDBUser,
		url.QueryEscape(cfg.GarminDBPassword), // Use the encoded password
		cfg.GarminDBHost,
		cfg.GarminDBPort,
		cfg.GarminDB,
	)

	db, err := gorm.Open(postgres.Open(connURL), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db.Logger = loggr

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if !cfg.IsProduction {
		return db.Debug(), nil
	}

	return db, nil
}

func NewGormDBVCC(cfg *common.Config, loggr *config.ZapLogger) (*gorm.DB, error) {

	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.VCCDBUser,
		url.QueryEscape(cfg.VCCDBPassword), // Use the encoded password
		cfg.VCCDBHost,
		cfg.VCCDBPort,
		cfg.VCCDBDatabase,
	)

	db, err := gorm.Open(postgres.Open(connURL), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db.Logger = loggr

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if !cfg.IsProduction {
		return db.Debug(), nil
	}

	return db, nil
}
