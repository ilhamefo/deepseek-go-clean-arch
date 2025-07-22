package gorm

import (
	"event-registration/internal/core/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GarminRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewGarminRepo(
	db *gorm.DB, // `name:"GarminDB"`
	logger *zap.Logger,
) domain.GarminRepository {
	return &GarminRepo{db: db, logger: logger}
}
