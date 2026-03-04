package domain

import (
	"time"
)

// ============================================================
// Activity Details - Database Models
// Maps to garmin_activity_schema.sql
// ============================================================

// MetricUnit represents measurement units (meter, watt, bpm, etc.)
type MetricUnit struct {
	UnitID  int     `json:"unitId" gorm:"column:unit_id;primaryKey"`
	UnitKey string  `json:"unitKey" gorm:"column:unit_key;size:50;uniqueIndex"`
	Factor  float64 `json:"factor" gorm:"column:factor;type:numeric(10,2)"`
}

func (MetricUnit) TableName() string {
	return "metric_units"
}

// MetricDescriptor defines available metrics for each activity
type MetricDescriptor struct {
	ID           int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivityID   int64     `json:"activityId" gorm:"column:activity_id;not null;index:idx_md_activity"`
	MetricsIndex int16     `json:"metricsIndex" gorm:"column:metrics_index;not null"`
	MetricKey    string    `json:"key" gorm:"column:metric_key;size:100;not null;index:idx_md_key"`
	UnitID       int       `json:"unitId" gorm:"column:unit_id;not null"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;not null;default:now()"`

	// Relations
	Unit MetricUnit `json:"unit" gorm:"foreignKey:UnitID;references:UnitID"`
}

func (MetricDescriptor) TableName() string {
	return "metric_descriptors"
}

// ActivityMetricsTimeseries stores detailed time-series metrics
// ~4000+ rows per activity
type ActivityMetricsTimeseries struct {
	ID         int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivityID int64 `json:"activityId" gorm:"column:activity_id;not null;index:idx_amts_activity"`
	Sequence   int   `json:"sequence" gorm:"column:sequence;not null"`

	// Metric index 0: sumDuration (second) - cumulative
	SumDuration *float64 `json:"sumDuration,omitempty" gorm:"column:sum_duration;type:numeric(12,2)"`

	// Metric index 1: directPower (watt)
	DirectPower *float64 `json:"directPower,omitempty" gorm:"column:direct_power;type:numeric(10,2)"`

	// Metric index 2: directGradeAdjustedSpeed (mps)
	DirectGradeAdjustedSpeed *float64 `json:"directGradeAdjustedSpeed,omitempty" gorm:"column:direct_grade_adjusted_speed;type:numeric(8,4)"`

	// Metric index 3: directAirTemperature (celsius)
	DirectAirTemperature *float64 `json:"directAirTemperature,omitempty" gorm:"column:direct_air_temperature;type:numeric(5,1)"`

	// Metric index 4: directHeartRate (bpm)
	DirectHeartRate *int16 `json:"directHeartRate,omitempty" gorm:"column:direct_heart_rate"`

	// Metric index 5: sumAccumulatedPower (watt) - cumulative
	SumAccumulatedPower *float64 `json:"sumAccumulatedPower,omitempty" gorm:"column:sum_accumulated_power;type:numeric(12,2)"`

	// Metric index 6: directFractionalCadence (stepsPerMinute)
	DirectFractionalCadence *float64 `json:"directFractionalCadence,omitempty" gorm:"column:direct_fractional_cadence;type:numeric(6,2)"`

	// Metric index 7: directBodyBattery (0-100)
	DirectBodyBattery *int16 `json:"directBodyBattery,omitempty" gorm:"column:direct_body_battery"`

	// Metric index 8: directElevation (meter)
	DirectElevation *float64 `json:"directElevation,omitempty" gorm:"column:direct_elevation;type:numeric(8,2)"`

	// Metric index 9: directRunCadence (stepsPerMinute)
	DirectRunCadence *int16 `json:"directRunCadence,omitempty" gorm:"column:direct_run_cadence"`

	// Metric index 10: directDoubleCadence (stepsPerMinute)
	DirectDoubleCadence *int16 `json:"directDoubleCadence,omitempty" gorm:"column:direct_double_cadence"`

	// Metric index 11: directSpeed (mps)
	DirectSpeed *float64 `json:"directSpeed,omitempty" gorm:"column:direct_speed;type:numeric(8,4)"`

	// Metric index 12: sumMovingDuration (second) - cumulative
	SumMovingDuration *float64 `json:"sumMovingDuration,omitempty" gorm:"column:sum_moving_duration;type:numeric(12,2)"`

	// Metric index 13: sumDistance (meter) - cumulative
	SumDistance *float64 `json:"sumDistance,omitempty" gorm:"column:sum_distance;type:numeric(12,2)"`

	// Metric index 14: sumElapsedDuration (second) - cumulative
	SumElapsedDuration *float64 `json:"sumElapsedDuration,omitempty" gorm:"column:sum_elapsed_duration;type:numeric(12,2)"`

	// Metric index 15: directTimestamp (unix epoch milliseconds)
	DirectTimestamp *int64 `json:"directTimestamp,omitempty" gorm:"column:direct_timestamp;index:idx_amts_timestamp"`

	// Metric index 16: directLongitude (decimal degrees)
	DirectLongitude *float64 `json:"directLongitude,omitempty" gorm:"column:direct_longitude"`

	// Metric index 17: directVerticalOscillation (centimeter)
	DirectVerticalOscillation *float64 `json:"directVerticalOscillation,omitempty" gorm:"column:direct_vertical_oscillation;type:numeric(6,2)"`

	// Metric index 18: directLatitude (decimal degrees)
	DirectLatitude *float64 `json:"directLatitude,omitempty" gorm:"column:direct_latitude"`

	// Metric index 19: directVerticalRatio (dimensionless)
	DirectVerticalRatio *float64 `json:"directVerticalRatio,omitempty" gorm:"column:direct_vertical_ratio;type:numeric(6,2)"`

	// Metric index 20: directStrideLength (centimeter)
	DirectStrideLength *float64 `json:"directStrideLength,omitempty" gorm:"column:direct_stride_length;type:numeric(6,2)"`

	// Metric index 21: directVerticalSpeed (mps)
	DirectVerticalSpeed *float64 `json:"directVerticalSpeed,omitempty" gorm:"column:direct_vertical_speed;type:numeric(8,4)"`

	// Metric index 22: directGroundContactTime (ms)
	DirectGroundContactTime *float64 `json:"directGroundContactTime,omitempty" gorm:"column:direct_ground_contact_time;type:numeric(6,2)"`

	// Computed columns (managed by DB)
	RecordedAt *time.Time `json:"recordedAt,omitempty" gorm:"column:recorded_at;->;index:idx_amts_recorded_at"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"column:created_at;not null;default:now()"`
}

func (ActivityMetricsTimeseries) TableName() string {
	return "activity_metrics_timeseries"
}

// GeoPolyline stores GPS route data
type GeoPolyline struct {
	ID           int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivityID   int64     `json:"activityId" gorm:"column:activity_id;not null;uniqueIndex"`
	StartPoint   *string   `json:"startPoint,omitempty" gorm:"column:start_point;type:jsonb"`
	EndPoint     *string   `json:"endPoint,omitempty" gorm:"column:end_point;type:jsonb"`
	MinLatitude  *float64  `json:"minLatitude,omitempty" gorm:"column:min_latitude"`
	MaxLatitude  *float64  `json:"maxLatitude,omitempty" gorm:"column:max_latitude"`
	MinLongitude *float64  `json:"minLongitude,omitempty" gorm:"column:min_longitude"`
	MaxLongitude *float64  `json:"maxLongitude,omitempty" gorm:"column:max_longitude"`
	Polyline     string    `json:"polyline" gorm:"column:polyline;type:jsonb;not null;default:'[]'"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;not null;default:now()"`
}

func (GeoPolyline) TableName() string {
	return "geo_polylines"
}

// HeartRateTimeseries stores dedicated heart rate measurements
type HeartRateTimeseries struct {
	ID          int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivityID  int64      `json:"activityId" gorm:"column:activity_id;not null;index:idx_hrts_activity"`
	TimestampMs int64      `json:"timestampMs" gorm:"column:timestamp_ms;not null"`
	HeartRate   int16      `json:"heartRate" gorm:"column:heart_rate;not null"`
	RecordedAt  *time.Time `json:"recordedAt,omitempty" gorm:"column:recorded_at;->;index:idx_hrts_recorded_at"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;not null;default:now()"`
}

func (HeartRateTimeseries) TableName() string {
	return "heart_rate_timeseries"
}

// ActivityDetailsSummary stores top-level metadata
type ActivityDetailsSummary struct {
	ID                int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivityID        int64     `json:"activityId" gorm:"column:activity_id;not null;uniqueIndex"`
	MeasurementCount  int       `json:"measurementCount" gorm:"column:measurement_count;not null"`
	MetricsCount      int       `json:"metricsCount" gorm:"column:metrics_count;not null"`
	TotalMetricsCount int       `json:"totalMetricsCount" gorm:"column:total_metrics_count;not null"`
	DetailsAvailable  bool      `json:"detailsAvailable" gorm:"column:details_available;not null;default:true"`
	PendingData       *string   `json:"pendingData,omitempty" gorm:"column:pending_data;type:jsonb"`
	FetchedAt         time.Time `json:"fetchedAt" gorm:"column:fetched_at;not null;default:now();index:idx_ads_fetched"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"column:updated_at;not null;default:now()"`
}

func (ActivityDetailsSummary) TableName() string {
	return "activity_details_summary"
}

// ============================================================
// JSON API Response Models
// For parsing Garmin API responses
// ============================================================

// ActivityDetailsResponse represents the complete JSON response
type ActivityDetailsResponse struct {
	ActivityID            int64                 `json:"activityId"`
	MeasurementCount      int                   `json:"measurementCount"`
	MetricsCount          int                   `json:"metricsCount"`
	TotalMetricsCount     int                   `json:"totalMetricsCount"`
	MetricDescriptors     []MetricDescriptorDTO `json:"metricDescriptors"`
	ActivityDetailMetrics []Metrics             `json:"activityDetailMetrics"`
	GeoPolylineDTO        GeoPolylineDTO        `json:"geoPolylineDTO"`
	HeartRateDTOs         []HeartRateDTO        `json:"heartRateDTOs"`
	PendingData           interface{}           `json:"pendingData"`
	DetailsAvailable      bool                  `json:"detailsAvailable"`
	ActivitySummary       *ActivityVM           `json:"activitySummary,omitempty"`
}

type Metrics struct {
	Metrics []float64 `json:"metrics"`
}

// MetricDescriptorDTO from JSON metricDescriptors[]
type MetricDescriptorDTO struct {
	MetricsIndex int     `json:"metricsIndex"`
	Key          string  `json:"key"`
	Unit         UnitDTO `json:"unit"`
}

// UnitDTO represents unit information in JSON
type UnitDTO struct {
	ID     int     `json:"id"`
	Key    string  `json:"key"`
	Factor float64 `json:"factor"`
}

// GeoPolylineDTO from JSON geoPolylineDTO
type GeoPolylineDTO struct {
	StartPoint interface{}   `json:"startPoint"`
	EndPoint   interface{}   `json:"endPoint"`
	MinLat     *float64      `json:"minLat"`
	MaxLat     *float64      `json:"maxLat"`
	MinLon     *float64      `json:"minLon"`
	MaxLon     *float64      `json:"maxLon"`
	Polyline   []interface{} `json:"polyline"`
}

// HeartRateDTO from JSON heartRateDTOs[]
type HeartRateDTO struct {
	TimestampMs int64 `json:"timestampMs"`
	HeartRate   int16 `json:"heartRate"`
}
