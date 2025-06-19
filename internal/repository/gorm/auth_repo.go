package gorm

import (
	"event-registration/internal/core/domain"
	"event-registration/internal/helper"

	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) domain.AuthRepository {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) IsRegistered(email string) (isRegistered bool, err error) {
	var count int64
	err = r.db.Table("users").
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		return isRegistered, err
	}

	return count > 0, err
}

func (r *AuthRepo) Register(user domain.User) (err error) {
	err = r.db.Model(&domain.User{}).Omit("ID").Create(&user).Error

	helper.PrettyPrint(user)

	return err
}
