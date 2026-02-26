package gorm

import (
	"context"
	"event-registration/internal/core/domain"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GarminDashboardRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewGarminDashboardRepo(
	db *gorm.DB, // `name:"GarminDB"`
	logger *zap.Logger,
) domain.GarminDashboardRepository {
	return &GarminDashboardRepo{db: db, logger: logger}
}

func (r *GarminDashboardRepo) GetHeartRate(ctx context.Context, current time.Time) (res domain.HeartRate, err error) {

	err = r.db.WithContext(ctx).
		Table("heart_rates").
		Where("calendar_date = ?", current.Format("2006-01-02")).
		Scan(&res).Error
	if err != nil {
		r.logger.Error("failed to get heart rate", zap.Error(err))
		return res, err
	}

	return res, nil
}
