package domain

import (
	"context"
	"time"
)

type GarminDashboardRepository interface {
	GetHeartRate(ctx context.Context, current time.Time) (res HeartRate, err error)
	GetActivities(ctx context.Context, cursor int64, limit int) (res []ActivityVM, nextCursor int64, hasMore bool, err error)
}

type ActivityVM struct {
	ActivityID         int          `json:"activity_id" gorm:"column:activity_id"`
	ActivityName       string       `json:"activityName" gorm:"column:activity_name"`
	Distance           float64      `json:"distance" gorm:"column:distance"`
	ActivityTypeID     int          `json:"-" gorm:"column:activity_type_id;index"`
	ActivityType       ActivityType `json:"activityType" gorm:"foreignKey:ActivityTypeID;references:TypeID"`
	AverageSpeed       float64      `json:"averageSpeed" gorm:"column:average_speed;type:decimal(8,5)"`
	AverageMovingSpeed float64      `json:"averageMovingSpeed" gorm:"column:average_moving_speed;type:decimal(8,5)"`
	MaxSpeed           float64      `json:"maxSpeed" gorm:"column:max_speed;type:decimal(8,5)"`
	Calories           float64      `json:"calories" gorm:"column:calories"`
}

func (ActivityVM) TableName() string {
	return "activities"
}

type ActivityPaginatedResponse struct {
	Data       []ActivityVM `json:"data"`
	NextCursor int64        `json:"nextCursor,omitempty"`
	HasMore    bool         `json:"hasMore"`
	Limit      int          `json:"limit"`
}
