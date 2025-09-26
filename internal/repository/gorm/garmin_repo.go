package gorm

import (
	"event-registration/internal/core/domain"

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
		return err
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

	// Delete existing heart rate for this activity and date
	if err := tx.Where("user_profile_pk = ? AND calendar_date = ?", data.UserProfilePK, data.CalendarDate).Delete(&domain.HeartRate{}).Error; err != nil {
		tx.Rollback()
		r.logger.Error("failed to delete existing heart rate", zap.Error(err), zap.Int64("user_profile_pk", data.UserProfilePK), zap.String("calendar_date", data.CalendarDate))
		return err
	}

	// Insert new heart rate
	if err := tx.Omit("ID").Clauses(clause.Returning{}).Create(data).Error; err != nil {
		tx.Rollback()
		r.logger.Error("failed to create heart rate", zap.Error(err), zap.Int64("activity_id", data.UserProfilePK), zap.String("date", data.CalendarDate))
		return err
	}

	// insert details
	if len(data.HeartRateValues) > 0 {

		var details []domain.HeartRateDetail

		for detail := range data.HeartRateValues {
			detail := domain.HeartRateDetail{
				HeartRate:     data.HeartRateValues[detail][1],
				Timestamp:     data.HeartRateValues[detail][0],
				UserProfilePK: data.UserProfilePK,
				CalendarDate:  data.CalendarDate,
			}

			details = append(details, detail)
		}

		r.logger.Info("heart rate details to insert", zap.Int("count", len(details)), zap.Int64("activity_id", data.UserProfilePK), zap.String("date", data.CalendarDate), zap.Any("details", details))

		if err := tx.Where("user_profile_pk = ? AND calendar_date = ?", data.UserProfilePK, data.CalendarDate).Delete(&domain.HeartRateDetail{}).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to delete existing heart rate details", zap.Error(err), zap.Int64("activity_id", data.UserProfilePK), zap.String("date", data.CalendarDate))
			return err
		}

		if err := tx.Create(details).Error; err != nil {
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

	return nil
}
