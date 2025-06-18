package gorm

import (
	"event-registration/internal/core/domain"

	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) domain.AuthRepository {
	return &AuthRepo{db: db}
}
