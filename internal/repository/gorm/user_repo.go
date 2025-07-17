package gorm

import (
	"errors"
	"event-registration/internal/common/constant"
	"event-registration/internal/core/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepo(
	db *gorm.DB, // `name:"VCCDB"`
	logger *zap.Logger,
) domain.UserRepository {
	return &UserRepo{db: db, logger: logger}
}

func (r *UserRepo) Search(key string) (user []*domain.UserVCC, err error) {

	key = "%" + key + "%"

	err = r.db.Model(&domain.UserVCC{}).
		Where("email ILIKE ?", key).
		Or("username ILIKE ?", key).
		Or("email ILIKE ?", key).
		Or("nip ILIKE ?", key).
		Or("full_name ILIKE ?", key).
		Or("id :: TEXT ILIKE ?", key).
		Preload("Roles").
		Limit(10).
		Find(&user).Error

	if err != nil {
		r.logger.Error(constant.SQL_ERROR, zap.Error(err))
		return user, handleGormError(err)
	}

	return user, nil
}

func (r *UserRepo) Roles() (user []*domain.Role, err error) {

	err = r.db.
		Model(&domain.Role{}).
		Scan(&user).Error

	if err != nil {
		r.logger.Error(constant.SQL_ERROR, zap.Error(err))
		return user, handleGormError(err)
	}

	return user, nil
}

func (r *UserRepo) Unit(level string) (units []*domain.UnitName, err error) {

	var query *gorm.DB

	switch level {
	case "1":
		query = r.db.Table("public.pln_unit_upi").
			Select("id_unit_upi as code, nama_unit_upi || ' - ' || id_unit_upi as label")
	case "2":
		query = r.db.Table("public.pln_unit_ap").
			Select("id_unit_ap as code, nama_unit_ap || ' - ' || id_unit_ap as label")
	case "3":
		query = r.db.Table("public.pln_unit_up").
			Select("id_unit_up as code, nama_unit_up || ' - ' || id_unit_up as label")
	default:
		return nil, errors.New("invalid level")
	}

	err = query.
		Order("label ASC").
		Scan(&units).Error
	if err != nil {
		r.logger.Error(constant.SQL_ERROR, zap.Error(err))
		return units, handleGormError(err)
	}

	return units, nil
}

// Update user record
func (r *UserRepo) Update(user *domain.UserVCC) (err error) {
	var updatableColumn []string = []string{
		"Email",
		"Username",
		"FullName",
		"Level",
		"Jabatan",
		"Company",
		"UnitCode",
		"UnitName",
		"Status",
	}

	roles := []map[string]interface{}{}

	tx := r.db.Begin()

	qry := tx.Select(updatableColumn).Updates(&user)

	// if qry.RowsAffected == 0 && qry.Error != nil {
	// 	tx.Rollback()
	// 	return errors.New("record_not_found")
	// }

	if qry.Error != nil {
		tx.Rollback()
		r.logger.Error(constant.SQL_ERROR, zap.Error(err))
		return handleGormError(qry.Error)
	}

	if len(user.Roles) > 0 {
		// hapus role user dlu
		err = tx.Delete(&domain.RoleUsers{}, "user_id = ?", user.ID).Error
		if err != nil {
			tx.Rollback()
			r.logger.Error(constant.SQL_ERROR, zap.Error(err))
			return handleGormError(err)
		}

		// build body untuk insert role baru
		for _, role := range user.Roles {
			roles = append(roles, map[string]interface{}{
				"UserID": user.ID,
				"RoleID": role.ID,
			})
		}

		// insert user roles
		err = tx.Model(&domain.RoleUsers{}).Create(roles).Error
		if err != nil {
			tx.Rollback()
			r.logger.Error(constant.SQL_ERROR, zap.Error(err))
			return handleGormError(err)
		}
	}

	return handleGormError(tx.Commit().Error)
}
