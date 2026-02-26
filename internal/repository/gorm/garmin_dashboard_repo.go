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
	return &GarminDashboardRepo{db: db.Debug(), logger: logger}
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

func (r *GarminDashboardRepo) GetActivities(ctx context.Context, cursor int64, limit int) (res []domain.ActivityVM, nextCursor int64, hasMore bool, err error) {

	query := r.db.WithContext(ctx).
		Select([]string{
			"activity_id",
			"activity_type_id",
			"activity_name",
			"distance",
			"average_speed",
			"max_speed",
			"calories",
			"begin_timestamp",
		}).
		Preload("ActivityType").
		Order("begin_timestamp DESC, activity_id DESC")

	if cursor > 0 {
		query = query.Where("activity_id < ?", cursor)
	}

	err = query.Limit(limit + 1).Find(&res).Error
	if err != nil {
		r.logger.Error("failed to get activities", zap.Error(err))
		return res, 0, false, err
	}

	hasMore = len(res) > limit
	if hasMore {
		res = res[:limit]
		nextCursor = int64(res[len(res)-1].ActivityID)
	} else if len(res) > 0 {
		nextCursor = int64(res[len(res)-1].ActivityID)
	}

	return res, nextCursor, hasMore, nil
}
