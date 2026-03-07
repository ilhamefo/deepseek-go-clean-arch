package gorm

import (
	"context"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"strconv"
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

func (r *GarminDashboardRepo) GetActivities(ctx context.Context, payload *request.ActivityDashboardRequest) (res []domain.ActivityVM, nextCursor int64, hasMore bool, err error) {

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
			"max_hr",
			"average_hr",
			"duration",
			"start_time_local",
			"total_sets",
		}).
		Preload("ActivityType")

	if payload.Cursor > 0 {
		query = query.Where("activity_id < ?", payload.Cursor)
	}

	// sort
	switch payload.SortBy {
	case "date":
		query = query.Order("begin_timestamp " + payload.SortOrder)
	case "distance":
		query = query.Order("distance " + payload.SortOrder)
	case "duration":
		query = query.Order("duration " + payload.SortOrder)
	case "calories":
		query = query.Order("calories " + payload.SortOrder)
	case "maxHr":
		query = query.Order("max_hr " + payload.SortOrder)
	case "avgPace":
		query = query.Order("average_speed " + payload.SortOrder)
	case "name":
		query = query.Order("activity_name " + payload.SortOrder)
	default:
		query = query.Order("begin_timestamp DESC, activity_id DESC")
	}

	if payload.Type != nil {
		switch strconv.Itoa(*payload.Type) {
		case "0":
			break
		default:
			query = query.Where("activity_type_id = ?", payload.Type)
		}
	}

	if len(payload.Search) > 0 {
		query = query.Where("activity_name ILIKE ?", "%"+payload.Search+"%")
	}

	err = query.Limit(payload.Limit + 1).Find(&res).Error
	if err != nil {
		r.logger.Error("failed to get activities", zap.Error(err))
		return res, 0, false, err
	}

	hasMore = len(res) > payload.Limit
	if hasMore {
		res = res[:payload.Limit]
		nextCursor = int64(res[len(res)-1].ActivityID)
	} else if len(res) > 0 {
		nextCursor = int64(res[len(res)-1].ActivityID)
	}

	return res, nextCursor, hasMore, nil
}

func (r *GarminDashboardRepo) GetActivityDetails(ctx context.Context, activityID int64) (metrics []*domain.ActivityMetricsTimeseries, res *domain.ActivityVM, err error) {
	err = r.db.WithContext(ctx).
		Where("activity_id = ?", activityID).
		Order("sequence ASC").
		Find(&metrics).Error

	if err != nil {
		r.logger.Error("failed to get activity metrics", zap.Error(err), zap.Int64("activity_id", activityID))
		return nil, nil, err
	}

	err = r.db.WithContext(ctx).
		Select([]string{
			"activity_id",
			"activity_type_id",
			"activity_name",
			"distance",
			"average_speed",
			"max_speed",
			"calories",
			"begin_timestamp",
			"max_hr",
			"average_hr",
			"duration",
			"start_time_local",
			"total_sets",
		}).
		Where("activity_id = ?", activityID).
		First(&res).Error

	if err != nil {
		r.logger.Error("failed to get activity summary", zap.Error(err), zap.Int64("activity_id", activityID))
		return nil, nil, err
	}

	return metrics, res, nil
}
