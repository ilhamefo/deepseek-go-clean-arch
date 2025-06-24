package gorm

import (
	"event-registration/internal/core/domain"

	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(
	db *gorm.DB, // `name:"authDB"`
) domain.AuthRepository {
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
	err = r.db.
		Model(&domain.User{}).
		Create(&user).
		Error

	return err
}

func (r *AuthRepo) FindByEmail(email string) (user *domain.User, err error) {
	err = r.db.Table("users").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
