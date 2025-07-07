package gorm

import (
	"event-registration/internal/core/domain"

	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) domain.TransactionRepository {
	return &TransactionRepo{db: db}
}
