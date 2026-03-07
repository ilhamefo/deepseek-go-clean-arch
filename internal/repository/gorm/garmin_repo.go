package gorm

import (
	"context"
	"event-registration/internal/core/domain"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *GarminRepo) UpsertSplits(activityID int64, splits *domain.ActivitySplitsResponse) (err error) {
	if splits == nil {
		r.logger.Info("no splits data to upsert")
		return nil
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return nil
	}

	// Delete existing laps for this activity
	if err := tx.Where("activity_id = ?", activityID).Delete(&domain.LapDTO{}).Error; err != nil {
		tx.Rollback()
		r.logger.Error("failed to delete existing laps", zap.Error(err), zap.Int64("activity_id", activityID))
		return err
	}

	// Insert new laps
	if len(splits.LapDTOs) > 0 {
		for i, lap := range splits.LapDTOs {
			lap.ActivityID = activityID
			if err := tx.Omit("id").Create(&lap).Error; err != nil {
				tx.Rollback()
				r.logger.Error("failed to create lap",
					zap.Error(err),
					zap.Int64("activity_id", activityID),
					zap.Int("lap_index", i))
				return err
			}
		}
		r.logger.Debug("laps inserted successfully",
			zap.Int64("activity_id", activityID),
			zap.Int("laps_count", len(splits.LapDTOs)))
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	r.logger.Info("splits upserted successfully",
		zap.Int64("activity_id", activityID),
		zap.Int("laps_count", len(splits.LapDTOs)))

	return nil
}

func (r *GarminRepo) Update(activities []*domain.Activity) (err error) {
	if len(activities) == 0 {
		r.logger.Info("no activities to update")
		return nil
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, activity := range activities {
		// Process nested objects first
		if err := r.processNestedObjects(tx, activity); err != nil {
			tx.Rollback()
			r.logger.Error("failed to process nested objects", zap.Error(err), zap.Int64("activity_id", activity.ActivityID))
			return err
		}

		// Upsert main activity
		if err := tx.Save(activity).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to upsert activity", zap.Error(err), zap.Int64("activity_id", activity.ActivityID))
			return err
		}

		// Process related data
		if err := r.processRelatedData(tx, activity); err != nil {
			tx.Rollback()
			r.logger.Error("failed to process related data", zap.Error(err), zap.Int64("activity_id", activity.ActivityID))
			return err
		}

		r.logger.Debug("activity processed successfully", zap.Int64("activity_id", activity.ActivityID))
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	r.logger.Info("activities updated successfully", zap.Int("count", len(activities)))
	return nil
}

func (r *GarminRepo) processNestedObjects(tx *gorm.DB, activity *domain.Activity) error {
	// Upsert ActivityType
	if activity.ActivityType.TypeID != 0 {
		if err := tx.Save(&activity.ActivityType).Error; err != nil {
			return err
		}
		activity.ActivityTypeID = activity.ActivityType.TypeID
	}

	// Upsert EventType
	if activity.EventType.TypeID != 0 {
		if err := tx.Save(&activity.EventType).Error; err != nil {
			return err
		}
		activity.EventTypeID = activity.EventType.TypeID
	}

	// Upsert Privacy
	if activity.Privacy.TypeID != 0 {
		if err := tx.Save(&activity.Privacy).Error; err != nil {
			return err
		}
		activity.PrivacyTypeID = activity.Privacy.TypeID
	}

	return nil
}

func (r *GarminRepo) processRelatedData(tx *gorm.DB, activity *domain.Activity) error {
	// Process UserRoles
	if len(activity.UserRoles) > 0 {
		// Delete existing user roles for this activity
		if err := tx.Where("activity_id = ?", activity.ActivityID).Delete(&domain.UserRole{}).Error; err != nil {
			return err
		}

		// Insert new user roles
		for _, roleName := range activity.UserRoles {
			userRole := &domain.UserRole{
				ActivityID: activity.ActivityID,
				RoleName:   roleName,
			}
			if err := tx.Select("ActivityID", "RoleName").Create(userRole).Error; err != nil {
				return err
			}
		}
	}

	// Process SplitSummaries
	if len(activity.SplitSummaries) > 0 {
		// Delete existing split summaries for this activity
		if err := tx.Where("activity_id = ?", activity.ActivityID).Delete(&domain.SplitSummary{}).Error; err != nil {
			return err
		}

		// Insert new split summaries
		for _, split := range activity.SplitSummaries {
			split.ActivityID = activity.ActivityID
			if err := tx.Omit("ID").Create(&split).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *GarminRepo) UpsertHeartRateByDate(ctx context.Context, data *domain.HeartRate) (err error) {
	if data == nil {
		r.logger.Info("no heart rate data to upsert")
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	if err = tx.Omit("ID").
		Clauses(
			clause.Returning{},
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_profile_pk"},
					{Name: "calendar_date"},
				},
				UpdateAll: true,
			}).
		Create(data).Error; err != nil {
		tx.Rollback()
		r.logger.Error("failed to create heart rate", zap.Error(err), zap.Int64("user_profile_pk", data.UserProfilePK), zap.String("date", data.CalendarDate))
		return err
	}

	// insert details
	if len(data.HeartRateValues) > 0 {

		var details []domain.HeartRateDetail

		for idx := range data.HeartRateValues {
			if len(data.HeartRateValues[idx]) < 2 {
				r.logger.Warn("heart rate value slice too short", zap.Int("index", idx))
				continue
			}
			detail := domain.HeartRateDetail{
				HeartRate:     data.HeartRateValues[idx][1],
				Timestamp:     data.HeartRateValues[idx][0],
				UserProfilePK: data.UserProfilePK,
				CalendarDate:  data.CalendarDate,
			}

			details = append(details, detail)
		}

		if err = tx.Clauses(
			clause.Returning{},
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_profile_pk"},
					{Name: "timestamp"},
				},
				UpdateAll: true,
			},
		).Create(&details).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to create heart rate details", zap.Error(err), zap.Int64("activity_id", data.UserProfilePK), zap.String("date", data.CalendarDate))
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	r.logger.Info("heart rate upserted successfully", zap.Int64("activity_id", data.UserProfilePK), zap.String("date", data.CalendarDate))

	return err
}

func (r *GarminRepo) UpsertUserSettings(data *domain.UserSetting) (err error) {
	if data == nil {
		return fmt.Errorf("no user settings data to upsert")
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	if err := r.upsertUserSetting(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.upsertUserData(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.upsertUserSleep(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.upsertUserSleepWindows(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.upsertPowerFormat(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.upsertHeartRateFormat(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.upsertHydrationContainers(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.upsertAvailableTrainingDays(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.upsertPreferredLongTrainingDays(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *GarminRepo) upsertUserSetting(tx *gorm.DB, data *domain.UserSetting) error {
	err := tx.Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(data).Error
	if err != nil {
		r.logger.Error("failed to upsert user settings", zap.Error(err), zap.Int64("id", data.ID))
	}
	data.UserData.UserProfilePK = data.ID
	data.UserSleep.UserProfilePK = data.ID
	for i := range data.UserSleepWindows {
		data.UserSleepWindows[i].UserProfilePK = data.ID
	}
	data.UserData.PowerFormat.UserProfilePK = data.ID
	data.UserData.HeartRateFormat.UserProfilePK = data.ID
	return err
}

func (r *GarminRepo) upsertUserData(tx *gorm.DB, data *domain.UserSetting) error {
	err := tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_profile_pk"}},
			UpdateAll: true,
		}).Create(&data.UserData).Error
	if err != nil {
		r.logger.Error("failed to upsert user data", zap.Error(err), zap.Int64("user_profile_pk", data.ID))
	}
	return err
}

func (r *GarminRepo) upsertUserSleep(tx *gorm.DB, data *domain.UserSetting) error {
	err := tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_profile_pk"}},
			UpdateAll: true,
		}).Create(&data.UserSleep).Error
	if err != nil {
		r.logger.Error("failed to upsert user sleep", zap.Error(err), zap.Int64("user_profile_pk", data.ID))
	}
	return err
}

func (r *GarminRepo) upsertUserSleepWindows(tx *gorm.DB, data *domain.UserSetting) error {
	err := tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_profile_pk"}, {Name: "sleep_window_frequency"}},
			UpdateAll: true,
		}).Create(&data.UserSleepWindows).Error
	if err != nil {
		r.logger.Error("failed to upsert user sleep windows", zap.Error(err), zap.Int64("user_profile_pk", data.ID))
	}
	return err
}

func (r *GarminRepo) upsertPowerFormat(tx *gorm.DB, data *domain.UserSetting) error {
	err := tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_profile_pk"}},
			UpdateAll: true,
		}).Create(&data.UserData.PowerFormat).Error
	if err != nil {
		r.logger.Error("failed to upsert user power format", zap.Error(err), zap.Int64("user_profile_pk", data.ID))
	}
	return err
}

func (r *GarminRepo) upsertHeartRateFormat(tx *gorm.DB, data *domain.UserSetting) error {
	err := tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_profile_pk"}},
			UpdateAll: true,
		}).Create(&data.UserData.HeartRateFormat).Error
	if err != nil {
		r.logger.Error("failed to upsert user heart rate format", zap.Error(err), zap.Int64("user_profile_pk", data.ID))
	}
	return err
}

func (r *GarminRepo) upsertHydrationContainers(tx *gorm.DB, data *domain.UserSetting) error {
	if len(data.UserData.HydrationContainers) > 0 {
		for i := range data.UserData.HydrationContainers {
			data.UserData.HydrationContainers[i].UserProfilePK = data.ID
		}
		err := tx.Omit("ID").Clauses(
			clause.Returning{Columns: []clause.Column{{Name: "id"}}},
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_profile_pk"},
					{Name: "volume"},
					{Name: "name"},
					{Name: "unit"},
				},
				UpdateAll: true,
			}).Create(&data.UserData.HydrationContainers).Error
		if err != nil {
			r.logger.Error("failed to upsert hydration containers", zap.Error(err), zap.Int64("user_profile_pk", data.ID))
		}
		return err
	}
	return nil
}

func (r *GarminRepo) upsertAvailableTrainingDays(tx *gorm.DB, data *domain.UserSetting) error {
	if len(data.UserData.AvailableTrainingDays) > 0 {
		model := domain.AvailableTrainingDays{
			UserProfilePK: data.ID,
			Days:          data.UserData.AvailableTrainingDays,
		}

		err := tx.Omit("ID").Clauses(
			clause.Returning{Columns: []clause.Column{{Name: "id"}}},
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_profile_pk"},
				},
				UpdateAll: true,
			}).Table((&domain.AvailableTrainingDays{}).TableName()).Create(&model).Error
		if err != nil {
			r.logger.Error("failed to upsert user available training days", zap.Error(err), zap.Int64("user_profile_pk", data.ID))
		}
		return err
	}
	return nil
}

func (r *GarminRepo) upsertPreferredLongTrainingDays(tx *gorm.DB, data *domain.UserSetting) error {
	if len(data.UserData.PreferredLongTrainingDays) > 0 {
		model := domain.PreferredLongTrainingDays{
			UserProfilePK: data.ID,
			Days:          data.UserData.PreferredLongTrainingDays,
		}

		err := tx.Omit("ID").Clauses(
			clause.Returning{Columns: []clause.Column{{Name: "id"}}},
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_profile_pk"},
				},
				UpdateAll: true,
			}).Table((&domain.PreferredLongTrainingDays{}).TableName()).Create(&model).Error
		if err != nil {
			r.logger.Error("failed to upsert user available training days", zap.Error(err), zap.Int64("user_profile_pk", data.ID))
		}
		return err
	}
	return nil
}

func (r *GarminRepo) UpsertSteps(ctx context.Context, data []*domain.Step) (err error) {
	if len(data) == 0 {
		r.logger.Info("no steps data to upsert")
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
	}

	err = tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_profile_pk"},
				{Name: "start_gmt"},
				{Name: "end_gmt"},
			},
			UpdateAll: true,
		}).Create(&data).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed to upsert steps", zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	r.logger.Info("steps upserted successfully", zap.Int("count", len(data)))

	return nil
}

func (r *GarminRepo) UpsertHRVByDate(ctx context.Context, data *domain.HRVData) (err error) {

	if data == nil {
		r.logger.Info("no hrv data to upsert")
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
	}

	err = tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_profile_pk"},
				{Name: "start_timestamp_gmt"},
				{Name: "end_timestamp_gmt"},
			},
			UpdateAll: true,
		}).Create(&data).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed to upsert hrv", zap.Error(err))
		return err
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
	}

	if len(data.HRVReadings) > 0 {
		for i := range data.HRVReadings {
			data.HRVReadings[i].UserProfilePK = data.UserProfilePK
			data.HRVReadings[i].ParentID = data.ID
		}
	}

	err = tx.Omit("ID").Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "parent_id"},
				{Name: "user_profile_pk"},
				{Name: "reading_time_gmt"},
			},
			UpdateAll: true,
		}).Create(&data.HRVReadings).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed to upsert hrv", zap.Error(err))
		return err
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
	}

	data.HRVSummary.UserProfilePK = data.UserProfilePK

	err = tx.Omit("ID").Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_profile_pk"},
				{Name: "calendar_date"},
			},
			UpdateAll: true,
		}).Create(&data.HRVSummary).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed to upsert hrv summary", zap.Error(err))
		return err
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
	}

	data.HRVSummary.Baseline.UserProfilePK = data.UserProfilePK
	data.HRVSummary.Baseline.CalendarDate = data.HRVSummary.CalendarDate

	err = tx.Omit("ID").Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_profile_pk"},
				{Name: "calendar_date"},
			},
			UpdateAll: true,
		}).Create(&data.HRVSummary.Baseline).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed to upsert hrv summary", zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (r *GarminRepo) UpsertActivityTypes(ctx context.Context, data []*domain.ActivityType) (err error) {
	if data == nil {
		r.logger.Info("no data to upsert")
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
	}

	err = tx.Omit("ID").Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "type_id"},
				{Name: "type_key"},
			},
			UpdateAll: true,
		}).Create(&data).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed to upsert", zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return err
}

func (r *GarminRepo) UpsertBodyBatteryByDate(ctx context.Context, data []*domain.StressData) (err error) {
	if len(data) == 0 {
		r.logger.Info("no_stress_data_to_upsert")
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	// exec here
	err = tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_profile_pk"},
				{Name: "calendar_date"},
				{Name: "event_start_time_gmt"},
				{Name: "activity_id"},
				{Name: "event_type"},
			},
			UpdateAll: true,
		}).
		Create(&data).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed to stress data", zap.Error(err))
		return err
	}

	var stressEvents []*domain.StressEvent

	for _, d := range data {
		stressEvents = append(stressEvents, &domain.StressEvent{
			ID:                     d.Event.ID,
			UserProfilePK:          d.UserProfilePK,
			EventType:              d.Event.EventType,
			EventStartTimeGmt:      d.Event.EventStartTimeGmt,
			TimezoneOffset:         d.Event.TimezoneOffset,
			DurationInMilliseconds: d.Event.DurationInMilliseconds,
			BodyBatteryImpact:      d.Event.BodyBatteryImpact,
			FeedbackType:           d.Event.FeedbackType,
			ShortFeedback:          d.Event.ShortFeedback,
			StressDataID:           d.ID,
		})
	}

	err = tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_profile_pk"},
				{Name: "event_start_time_gmt"},
			},
			UpdateAll: true,
		}).
		Create(&stressEvents).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed to stress event", zap.Error(err))
		return err
	}

	// exec end here

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (r *GarminRepo) HealthCheck(ctx context.Context) (err error) {
	// Check if database connection is alive
	sqlDB, err := r.db.DB()
	if err != nil {
		r.logger.Error("failed to get underlying sql.DB", zap.Error(err))
		return err
	}

	// Ping database with context
	if err := sqlDB.PingContext(ctx); err != nil {
		r.logger.Error("database ping failed", zap.Error(err))
		return err
	}

	return nil
}

func (r *GarminRepo) UpsertSleepByDate(ctx context.Context, data *domain.SleepResponse) (err error) {

	if data == nil {
		r.logger.Info("no sleep data to upsert")
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
	}

	// 1. Upsert DailySleepDTO (main sleep data)
	err = tx.Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_profile_pk"},
				{Name: "calendar_date"},
			},
			UpdateAll: true,
		}).
		Create(&data.DailySleepDTO).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed upsert daily sleep", zap.Error(err))
		return err
	}

	sleepID := data.DailySleepDTO.ID
	r.logger.Debug("daily sleep upserted", zap.Int64("sleep_id", sleepID))

	// 2. Upsert SleepScores and details
	if err := r.upsertSleepScores(tx, sleepID, &data.DailySleepDTO.SleepScores); err != nil {
		tx.Rollback()
		return err
	}

	// 3. Upsert SleepNeed (current and next)
	if err := r.upsertSleepNeeds(tx, sleepID, &data.DailySleepDTO.SleepNeed, &data.DailySleepDTO.NextSleepNeed); err != nil {
		tx.Rollback()
		return err
	}

	// 4. Upsert SleepMovements
	if len(data.SleepMovement) > 0 {
		if err := r.upsertSleepMovements(tx, sleepID, data.SleepMovement); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 5. Upsert SleepLevels
	if len(data.SleepLevels) > 0 {
		if err := r.upsertSleepLevels(tx, sleepID, data.SleepLevels); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 6. Upsert SleepRestlessMoments
	if len(data.SleepRestlessMoments) > 0 {
		if err := r.upsertSleepRestlessMoments(tx, sleepID, data.SleepRestlessMoments); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 7. Upsert SpO2 Summary
	if err := r.upsertSpO2Summary(tx, sleepID, &data.WellnessSpO2SleepSummaryDTO); err != nil {
		tx.Rollback()
		return err
	}

	// 8. Upsert SpO2 Data
	if len(data.WellnessEpochSPO2DataDTOList) > 0 {
		if err := r.upsertSpO2Data(tx, sleepID, data.WellnessEpochSPO2DataDTOList); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 9. Upsert Respiration Data
	if len(data.WellnessEpochRespirationDataDTOList) > 0 {
		if err := r.upsertRespirationData(tx, sleepID, data.WellnessEpochRespirationDataDTOList); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 10. Upsert Heart Rate Data
	if len(data.SleepHeartRate) > 0 {
		if err := r.upsertSleepHeartRate(tx, sleepID, data.SleepHeartRate); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 11. Upsert Stress Data
	if len(data.SleepStress) > 0 {
		if err := r.upsertSleepStress(tx, sleepID, data.SleepStress); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 12. Upsert Body Battery Data
	if len(data.SleepBodyBattery) > 0 {
		if err := r.upsertSleepBodyBattery(tx, sleepID, data.SleepBodyBattery); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 13. Upsert HRV Data
	if len(data.HrvData) > 0 {
		if err := r.upsertSleepHRV(tx, sleepID, data.HrvData); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 14. Upsert Breathing Disruption Data
	if len(data.BreathingDisruptionData) > 0 {
		if err := r.upsertBreathingDisruption(tx, sleepID, data.BreathingDisruptionData); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	r.logger.Info("sleep data upserted successfully", zap.Int64("sleep_id", sleepID))
	return nil
}

func (r *GarminRepo) GetActivity(ctx context.Context, id string) (activity *domain.Activity, err error) {
	err = r.db.WithContext(ctx).
		Model(&activity).
		Where("activity_id = ?", id).Scan(&activity).Error
	if err != nil {
		r.logger.Error("failed to get activity", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	return activity, nil
}

// UpsertActivityDetails inserts or updates activity details from Garmin API
func (r *GarminRepo) UpsertActivityDetails(ctx context.Context, data *domain.ActivityDetailsResponse) (err error) {
	if data == nil {
		r.logger.Info("no activity details to upsert")
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic during activity details upsert: %v", r)
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	summary := &domain.ActivityDetailsSummary{
		ActivityID:        data.ActivityID,
		MeasurementCount:  data.MeasurementCount,
		MetricsCount:      data.MetricsCount,
		TotalMetricsCount: data.TotalMetricsCount,
		DetailsAvailable:  data.DetailsAvailable,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "activity_id"}},
		UpdateAll: true,
	}).Create(summary).Error; err != nil {
		tx.Rollback()
		r.logger.Error("failed to upsert activity_details_summary", zap.Error(err), zap.Int64("activity_id", data.ActivityID))
		return err
	}

	// 2. Upsert metric_descriptors
	if len(data.MetricDescriptors) > 0 {
		// Delete existing descriptors
		if err := tx.Where("activity_id = ?", data.ActivityID).Delete(&domain.MetricDescriptor{}).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to delete existing metric_descriptors", zap.Error(err))
			return err
		}

		// Insert new descriptors
		var descriptors []domain.MetricDescriptor
		for _, desc := range data.MetricDescriptors {
			descriptors = append(descriptors, domain.MetricDescriptor{
				ActivityID:   data.ActivityID,
				MetricsIndex: int16(desc.MetricsIndex),
				MetricKey:    desc.Key,
				UnitID:       desc.Unit.ID,
			})
		}

		if err := tx.CreateInBatches(descriptors, 100).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to insert metric_descriptors", zap.Error(err))
			return err
		}
	}

	// 3. Upsert activity_metrics_timeseries
	if len(data.ActivityDetailMetrics) > 0 {
		// Delete existing metrics for this activity
		if err := tx.Where("activity_id = ?", data.ActivityID).Delete(&domain.ActivityMetricsTimeseries{}).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to delete existing activity_metrics_timeseries", zap.Error(err))
			return err
		}

		// Build timeseries records
		var timeseries []domain.ActivityMetricsTimeseries
		for seq, metrics := range data.ActivityDetailMetrics {
			record := domain.ActivityMetricsTimeseries{
				ActivityID: data.ActivityID,
				Sequence:   seq,
			}

			// Map metrics by index (based on metricDescriptors order)
			for i, value := range metrics.Metrics {
				if i >= len(data.MetricDescriptors) {
					break
				}

				// Skip null values
				if value == 0 {
					continue
				}

				key := data.MetricDescriptors[i].Key
				factor := data.MetricDescriptors[i].Unit.Factor

				// Apply factor and assign to correct field
				val := value / factor

				switch key {
				case "sumDuration":
					record.SumDuration = &val
				case "directPower":
					record.DirectPower = &val
				case "directGradeAdjustedSpeed":
					record.DirectGradeAdjustedSpeed = &val
				case "directAirTemperature":
					record.DirectAirTemperature = &val
				case "directHeartRate":
					hr := int16(val)
					record.DirectHeartRate = &hr
				case "sumAccumulatedPower":
					record.SumAccumulatedPower = &val
				case "directFractionalCadence":
					record.DirectFractionalCadence = &val
				case "directBodyBattery":
					bb := int16(val)
					record.DirectBodyBattery = &bb
				case "directElevation":
					record.DirectElevation = &val
				case "directRunCadence":
					rc := int16(val)
					record.DirectRunCadence = &rc
				case "directDoubleCadence":
					dc := int16(val)
					record.DirectDoubleCadence = &dc
				case "directSpeed":
					record.DirectSpeed = &val
				case "sumMovingDuration":
					record.SumMovingDuration = &val
				case "sumDistance":
					record.SumDistance = &val
				case "sumElapsedDuration":
					record.SumElapsedDuration = &val
				case "directTimestamp":
					ts := int64(value)
					record.DirectTimestamp = &ts
				case "directLongitude":
					record.DirectLongitude = &val
				case "directVerticalOscillation":
					record.DirectVerticalOscillation = &val
				case "directLatitude":
					record.DirectLatitude = &val
				case "directVerticalRatio":
					record.DirectVerticalRatio = &val
				case "directStrideLength":
					record.DirectStrideLength = &val
				case "directVerticalSpeed":
					record.DirectVerticalSpeed = &val
				case "directGroundContactTime":
					record.DirectGroundContactTime = &val
				}
			}

			timeseries = append(timeseries, record)
		}

		// Insert in batches (1000 records at a time)
		if err := tx.CreateInBatches(timeseries, 1000).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to insert activity_metrics_timeseries", zap.Error(err))
			return err
		}

		r.logger.Info("inserted activity_metrics_timeseries",
			zap.Int64("activity_id", data.ActivityID),
			zap.Int("count", len(timeseries)))
	}

	// 4. Upsert geo_polylines
	if len(data.GeoPolylineDTO.Polyline) > 0 || data.GeoPolylineDTO.MinLat != nil {
		polyline := &domain.GeoPolyline{
			ActivityID:   data.ActivityID,
			MinLatitude:  data.GeoPolylineDTO.MinLat,
			MaxLatitude:  data.GeoPolylineDTO.MaxLat,
			MinLongitude: data.GeoPolylineDTO.MinLon,
			MaxLongitude: data.GeoPolylineDTO.MaxLon,
			Polyline:     "[]",
		}

		// Convert polyline to JSON string if available
		if len(data.GeoPolylineDTO.Polyline) > 0 {
			// Polyline is already []interface{}, convert to JSON
			polyline.Polyline = fmt.Sprintf("%v", data.GeoPolylineDTO.Polyline)
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "activity_id"}},
			UpdateAll: true,
		}).Create(polyline).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to upsert geo_polyline", zap.Error(err))
			return err
		}
	}

	// 5. Upsert heart_rate_timeseries
	if len(data.HeartRateDTOs) > 0 {
		// Delete existing HR data
		if err := tx.Where("activity_id = ?", data.ActivityID).Delete(&domain.HeartRateTimeseries{}).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to delete existing heart_rate_timeseries", zap.Error(err))
			return err
		}

		// Insert new HR data
		var hrData []domain.HeartRateTimeseries
		for _, hr := range data.HeartRateDTOs {
			hrData = append(hrData, domain.HeartRateTimeseries{
				ActivityID:  data.ActivityID,
				TimestampMs: hr.TimestampMs,
				HeartRate:   hr.HeartRate,
			})
		}

		if err := tx.CreateInBatches(hrData, 1000).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to insert heart_rate_timeseries", zap.Error(err))
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	r.logger.Info("activity details upserted successfully",
		zap.Int64("activity_id", data.ActivityID),
		zap.Int("metrics_count", len(data.ActivityDetailMetrics)))

	return nil
}

// upsertSleepScores handles inserting/updating sleep scores and their details
func (r *GarminRepo) upsertSleepScores(tx *gorm.DB, sleepID int64, scores *domain.SleepScores) error {
	if scores == nil {
		return nil
	}

	scores.SleepID = sleepID

	// Upsert main sleep scores record
	if err := tx.Omit("ID").Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "sleep_id"}},
			UpdateAll: true,
		}).Create(scores).Error; err != nil {
		r.logger.Error("failed to upsert sleep_scores", zap.Error(err), zap.Int64("sleep_id", sleepID))
		return err
	}

	scoresID := scores.ID
	r.logger.Debug("sleep scores upserted", zap.String("scores_id", scoresID), zap.Int64("sleep_id", sleepID))

	// Delete existing score details
	if err := tx.Where("sleep_scores_id = ?", scoresID).Delete(&domain.SleepScore{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep score details", zap.Error(err))
		return err
	}

	// Insert all score details
	var scoreDetails []domain.SleepScore
	scoreTypes := []struct {
		scoreType string
		score     domain.SleepScore
	}{
		{"totalDuration", scores.TotalDuration},
		{"stress", scores.Stress},
		{"awakeCount", scores.AwakeCount},
		{"remPercentage", scores.RemPercentage},
		{"restlessness", scores.Restlessness},
		{"lightPercentage", scores.LightPercentage},
		{"deepPercentage", scores.DeepPercentage},
	}

	for _, st := range scoreTypes {
		detail := st.score
		detail.SleepScoresID = scoresID
		detail.ScoreType = st.scoreType
		scoreDetails = append(scoreDetails, detail)
	}

	if len(scoreDetails) > 0 {
		if err := tx.Omit("ID").Create(&scoreDetails).Error; err != nil {
			r.logger.Error("failed to insert sleep score details", zap.Error(err))
			return err
		}
		r.logger.Debug("sleep score details inserted", zap.Int("count", len(scoreDetails)))
	}

	return nil
}

// upsertSleepNeeds handles inserting/updating current and next sleep needs
func (r *GarminRepo) upsertSleepNeeds(tx *gorm.DB, sleepID int64, current *domain.SleepNeed, next *domain.SleepNeed) error {
	var needs []domain.SleepNeed

	if current != nil && current.UserProfilePk > 0 {
		current.SleepID = sleepID
		current.NeedType = "current"
		needs = append(needs, *current)
	}

	if next != nil && next.UserProfilePk > 0 {
		next.SleepID = sleepID
		next.NeedType = "next"
		needs = append(needs, *next)
	}

	if len(needs) == 0 {
		return nil
	}

	for _, need := range needs {
		if err := tx.Omit("ID").Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "sleep_id"},
					{Name: "need_type"},
				},
				UpdateAll: true,
			}).Create(&need).Error; err != nil {
			r.logger.Error("failed to upsert sleep_need", zap.Error(err), zap.String("need_type", need.NeedType))
			return err
		}
	}

	r.logger.Debug("sleep needs upserted", zap.Int("count", len(needs)))
	return nil
}

// upsertSleepMovements handles inserting/updating sleep movements
func (r *GarminRepo) upsertSleepMovements(tx *gorm.DB, sleepID int64, movements []domain.SleepMovement) error {
	// Delete existing movements for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.SleepMovement{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_movements", zap.Error(err))
		return err
	}

	// Set sleep_id for all movements
	for i := range movements {
		movements[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(movements, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_movements", zap.Error(err))
		return err
	}

	r.logger.Debug("sleep movements inserted", zap.Int("count", len(movements)))
	return nil
}

// upsertSleepLevels handles inserting/updating sleep levels
func (r *GarminRepo) upsertSleepLevels(tx *gorm.DB, sleepID int64, levels []domain.SleepLevel) error {
	// Delete existing levels for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.SleepLevel{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_levels", zap.Error(err))
		return err
	}

	// Set sleep_id for all levels
	for i := range levels {
		levels[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(levels, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_levels", zap.Error(err))
		return err
	}

	r.logger.Debug("sleep levels inserted", zap.Int("count", len(levels)))
	return nil
}

// upsertSleepRestlessMoments handles inserting/updating restless moments
func (r *GarminRepo) upsertSleepRestlessMoments(tx *gorm.DB, sleepID int64, moments []domain.SleepRestlessMoment) error {
	// Delete existing restless moments for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.SleepRestlessMoment{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_restless_moments", zap.Error(err))
		return err
	}

	// Set sleep_id for all moments
	for i := range moments {
		moments[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(moments, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_restless_moments", zap.Error(err))
		return err
	}

	r.logger.Debug("sleep restless moments inserted", zap.Int("count", len(moments)))
	return nil
}

// upsertSpO2Summary handles inserting/updating SpO2 summary
func (r *GarminRepo) upsertSpO2Summary(tx *gorm.DB, sleepID int64, summary *domain.WellnessSpO2SleepSummaryDTO) error {
	if summary == nil || summary.UserProfilePk == 0 {
		return nil
	}

	summary.SleepID = sleepID

	if err := tx.Omit("ID").Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "sleep_id"}},
			UpdateAll: true,
		}).Create(summary).Error; err != nil {
		r.logger.Error("failed to upsert sleep_spo2_summary", zap.Error(err))
		return err
	}

	r.logger.Debug("SpO2 summary upserted", zap.Int64("sleep_id", sleepID))
	return nil
}

// upsertSpO2Data handles inserting/updating SpO2 readings
func (r *GarminRepo) upsertSpO2Data(tx *gorm.DB, sleepID int64, data []domain.WellnessEpochSPO2DataDTO) error {
	// Delete existing SpO2 data for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.WellnessEpochSPO2DataDTO{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_spo2_data", zap.Error(err))
		return err
	}

	// Set sleep_id for all data points
	for i := range data {
		data[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(data, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_spo2_data", zap.Error(err))
		return err
	}

	r.logger.Debug("SpO2 data inserted", zap.Int("count", len(data)))
	return nil
}

// upsertRespirationData handles inserting/updating respiration data
func (r *GarminRepo) upsertRespirationData(tx *gorm.DB, sleepID int64, data []domain.WellnessEpochRespirationDataDTO) error {
	// Delete existing respiration data for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.WellnessEpochRespirationDataDTO{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_respiration_data", zap.Error(err))
		return err
	}

	// Set sleep_id for all data points
	for i := range data {
		data[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(data, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_respiration_data", zap.Error(err))
		return err
	}

	r.logger.Debug("respiration data inserted", zap.Int("count", len(data)))
	return nil
}

// upsertSleepHeartRate handles inserting/updating heart rate data
func (r *GarminRepo) upsertSleepHeartRate(tx *gorm.DB, sleepID int64, data []domain.SleepHeartRate) error {
	// Delete existing heart rate data for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.SleepHeartRate{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_heart_rate", zap.Error(err))
		return err
	}

	// Set sleep_id for all data points
	for i := range data {
		data[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(data, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_heart_rate", zap.Error(err))
		return err
	}

	r.logger.Debug("heart rate data inserted", zap.Int("count", len(data)))
	return nil
}

// upsertSleepStress handles inserting/updating stress data
func (r *GarminRepo) upsertSleepStress(tx *gorm.DB, sleepID int64, data []domain.SleepStress) error {
	// Delete existing stress data for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.SleepStress{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_stress", zap.Error(err))
		return err
	}

	// Set sleep_id for all data points
	for i := range data {
		data[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(data, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_stress", zap.Error(err))
		return err
	}

	r.logger.Debug("stress data inserted", zap.Int("count", len(data)))
	return nil
}

// upsertSleepBodyBattery handles inserting/updating body battery data
func (r *GarminRepo) upsertSleepBodyBattery(tx *gorm.DB, sleepID int64, data []domain.SleepBodyBattery) error {
	// Delete existing body battery data for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.SleepBodyBattery{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_body_battery", zap.Error(err))
		return err
	}

	// Set sleep_id for all data points
	for i := range data {
		data[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(data, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_body_battery", zap.Error(err))
		return err
	}

	r.logger.Debug("body battery data inserted", zap.Int("count", len(data)))
	return nil
}

// upsertSleepHRV handles inserting/updating HRV data
func (r *GarminRepo) upsertSleepHRV(tx *gorm.DB, sleepID int64, data []domain.HrvData) error {
	// Delete existing HRV data for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.HrvData{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_hrv_data", zap.Error(err))
		return err
	}

	// Set sleep_id for all data points
	for i := range data {
		data[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(data, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_hrv_data", zap.Error(err))
		return err
	}

	r.logger.Debug("HRV data inserted", zap.Int("count", len(data)))
	return nil
}

// upsertBreathingDisruption handles inserting/updating breathing disruption data
func (r *GarminRepo) upsertBreathingDisruption(tx *gorm.DB, sleepID int64, data []domain.BreathingDisruptionData) error {
	// Delete existing breathing disruption data for this sleep
	if err := tx.Where("sleep_id = ?", sleepID).Delete(&domain.BreathingDisruptionData{}).Error; err != nil {
		r.logger.Error("failed to delete existing sleep_breathing_disruption", zap.Error(err))
		return err
	}

	// Set sleep_id for all data points
	for i := range data {
		data[i].SleepID = sleepID
	}

	// Insert in batches
	if err := tx.Omit("ID").CreateInBatches(data, 500).Error; err != nil {
		r.logger.Error("failed to insert sleep_breathing_disruption", zap.Error(err))
		return err
	}

	r.logger.Debug("breathing disruption data inserted", zap.Int("count", len(data)))
	return nil
}
