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

func (r *GarminRepo) UpsertHeartRateByDate(data *domain.HeartRate) (err error) {
	if data == nil {
		r.logger.Info("no heart rate data to upsert")
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
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "type_id"},
				{Name: "type_key"},
			},
			UpdateAll: true,
		}).Create(&data).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error("failed to upsert hrv", zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return err
}
