package domain

import (
	"context"
	"time"
)

type GarminDashboardRepository interface {
	GetHeartRate(ctx context.Context, current time.Time) (res HeartRate, err error)
}

type HeartRateResponse struct {
	ID                               string     `json:"id"`
	UserProfilePK                    int64      `json:"user_profile_pk"`
	CalendarDate                     time.Time  `json:"calendar_date"`
	StartTimestampGMT                time.Time  `json:"start_timestamp_gmt"`
	EndTimestampGMT                  time.Time  `json:"end_timestamp_gmt"`
	StartTimestampLocal              time.Time  `json:"start_timestamp_local"`
	EndTimestampLocal                time.Time  `json:"end_timestamp_local"`
	MaxHeartRate                     int        `json:"max_heart_rate"`
	MinHeartRate                     int        `json:"min_heart_rate"`
	RestingHeartRate                 int        `json:"resting_heart_rate"`
	LastSevenDaysAvgRestingHeartRate int        `json:"last_seven_days_avg_resting_heart_rate"`
	CreatedAt                        *time.Time `json:"created_at"`
	UpdatedAt                        *time.Time `json:"updated_at"`
}
