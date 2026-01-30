package gorm

import (
	"event-registration/internal/common/constant"
	"event-registration/internal/core/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAuthRepo(
	db *gorm.DB, // `name:"authDB"`
	logger *zap.Logger,
) domain.AuthRepository {
	return &AuthRepo{db: db, logger: logger}
}

func (r *AuthRepo) IsRegistered(email string) (isRegistered bool, err error) {
	var count int64
	err = r.db.Table("users").
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		r.logger.Error(constant.SQL_ERROR, zap.Error(err))
		return isRegistered, handleGormError(err)
	}

	return count > 0, nil
}

func (r *AuthRepo) Register(user domain.User) (err error) {
	err = r.db.
		Model(&domain.User{}).
		Create(&user).
		Error

	if err != nil {
		r.logger.Error(constant.SQL_ERROR, zap.Error(err))
		return handleGormError(err)
	}

	return nil
}

func (r *AuthRepo) FindByEmail(email string) (user *domain.User, err error) {
	err = r.db.Table("users").
		Select("id", "email", "password", "name", "picture").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		r.logger.Error(constant.SQL_ERROR, zap.Error(err))
		return user, handleGormError(err)
	}

	return user, nil
}
