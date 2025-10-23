package domain

import (
	"context"
)

type GarminRepository interface {
	Update(activities []*Activity) error
	UpsertHeartRateByDate(ctx context.Context, data *HeartRate) error
	UpsertUserSettings(data *UserSetting) (err error)
	UpsertSteps(ctx context.Context, data []*Step) (err error)
	UpsertHRVByDate(ctx context.Context, data *HRVData) (err error)
	UpsertActivityTypes(ctx context.Context, data []*ActivityType) (err error)
	UpsertBodyBatteryByDate(ctx context.Context, data []*StressData) (err error)
	UpsertSleepByDate(ctx context.Context, data *SleepResponse) (err error)
}

type ActivityType struct {
	TypeID       int    `json:"typeId" gorm:"column:type_id;primaryKey"`
	TypeKey      string `json:"typeKey" gorm:"column:type_key;size:50"`
	ParentTypeID int    `json:"parentTypeId" gorm:"column:parent_type_id"`
	IsHidden     bool   `json:"isHidden" gorm:"column:is_hidden"`
	Restricted   bool   `json:"restricted" gorm:"column:restricted"`
	Trimmable    bool   `json:"trimmable" gorm:"column:trimmable"`
}

func (ActivityType) TableName() string {
	return "activity_types"
}

type EventType struct {
	TypeID    int    `json:"typeId" gorm:"column:type_id;primaryKey"`
	TypeKey   string `json:"typeKey" gorm:"column:type_key;size:50"`
	SortOrder int    `json:"sortOrder" gorm:"column:sort_order"`
}

func (EventType) TableName() string {
	return "event_types"
}

type Privacy struct {
	TypeID  int    `json:"typeId" gorm:"column:type_id;primaryKey"`
	TypeKey string `json:"typeKey" gorm:"column:type_key;size:50"`
}

func (Privacy) TableName() string {
	return "privacy_settings"
}

type SummarizedDiveInfo struct {
	SummarizedDiveGases []interface{} `json:"summarizedDiveGases"`
}

type UserRole struct {
	ID         string `json:"id" gorm:"column:id;primaryKey"`
	ActivityID int64  `json:"activityId" gorm:"column:activity_id;index"`
	RoleName   string `json:"roleName" gorm:"column:role_name;size:100"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type SplitSummary struct {
	ID                   string  `json:"id" gorm:"column:id;primaryKey"`
	ActivityID           int64   `json:"activityId" gorm:"column:activity_id;index"`
	NoOfSplits           int     `json:"noOfSplits" gorm:"column:no_of_splits"`
	TotalAscent          float64 `json:"totalAscent" gorm:"column:total_ascent"`
	Duration             float64 `json:"duration" gorm:"column:duration"`
	SplitType            string  `json:"splitType" gorm:"column:split_type;size:50"`
	NumClimbSends        int     `json:"numClimbSends" gorm:"column:num_climb_sends"`
	MaxElevationGain     float64 `json:"maxElevationGain" gorm:"column:max_elevation_gain"`
	AverageElevationGain float64 `json:"averageElevationGain" gorm:"column:average_elevation_gain"`
	MaxDistance          int     `json:"maxDistance" gorm:"column:max_distance"`
	Distance             float64 `json:"distance" gorm:"column:distance;type:decimal(10,6)"`
	AverageSpeed         float64 `json:"averageSpeed" gorm:"column:average_speed;type:decimal(8,6)"`
	MaxSpeed             float64 `json:"maxSpeed" gorm:"column:max_speed;type:decimal(8,6)"`
	NumFalls             int     `json:"numFalls" gorm:"column:num_falls"`
	ElevationLoss        float64 `json:"elevationLoss" gorm:"column:elevation_loss"`
}

func (SplitSummary) TableName() string {
	return "split_summaries"
}

type Activity struct {
	ActivityID                            int64   `json:"activityId" gorm:"column:activity_id;primaryKey"`
	ActivityName                          string  `json:"activityName" gorm:"column:activity_name"`
	StartTimeLocal                        string  `json:"startTimeLocal" gorm:"column:start_time_local;type:timestamp"`
	StartTimeGMT                          string  `json:"startTimeGMT" gorm:"column:start_time_gmt;type:timestamp"`
	EndTimeGMT                            string  `json:"endTimeGMT" gorm:"column:end_time_gmt;type:timestamp"`
	Distance                              float64 `json:"distance" gorm:"column:distance"`
	Duration                              float64 `json:"duration" gorm:"column:duration"`
	ElapsedDuration                       float64 `json:"elapsedDuration" gorm:"column:elapsed_duration"`
	MovingDuration                        float64 `json:"movingDuration" gorm:"column:moving_duration"`
	ElevationGain                         float64 `json:"elevationGain" gorm:"column:elevation_gain"`
	ElevationLoss                         float64 `json:"elevationLoss" gorm:"column:elevation_loss"`
	AverageSpeed                          float64 `json:"averageSpeed" gorm:"column:average_speed"`
	MaxSpeed                              float64 `json:"maxSpeed" gorm:"column:max_speed"`
	StartLatitude                         float64 `json:"startLatitude" gorm:"column:start_latitude"`
	StartLongitude                        float64 `json:"startLongitude" gorm:"column:start_longitude"`
	EndLatitude                           float64 `json:"endLatitude" gorm:"column:end_latitude"`
	EndLongitude                          float64 `json:"endLongitude" gorm:"column:end_longitude"`
	HasPolyline                           bool    `json:"hasPolyline" gorm:"column:has_polyline"`
	HasImages                             bool    `json:"hasImages" gorm:"column:has_images"`
	OwnerID                               int64   `json:"ownerId" gorm:"column:owner_id"`
	OwnerDisplayName                      string  `json:"ownerDisplayName" gorm:"column:owner_display_name"`
	OwnerFullName                         string  `json:"ownerFullName" gorm:"column:owner_full_name"`
	OwnerProfileImageURLSmall             string  `json:"ownerProfileImageUrlSmall" gorm:"column:owner_profile_image_url_small"`
	OwnerProfileImageURLMedium            string  `json:"ownerProfileImageUrlMedium" gorm:"column:owner_profile_image_url_medium"`
	OwnerProfileImageURLLarge             string  `json:"ownerProfileImageUrlLarge" gorm:"column:owner_profile_image_url_large"`
	Calories                              float64 `json:"calories" gorm:"column:calories"`
	BMRCalories                           float64 `json:"bmrCalories" gorm:"column:bmr_calories"`
	AverageHR                             float64 `json:"averageHR" gorm:"column:average_hr"`
	MaxHR                                 float64 `json:"maxHR" gorm:"column:max_hr"`
	AverageRunningCadenceInStepsPerMinute float64 `json:"averageRunningCadenceInStepsPerMinute" gorm:"column:average_running_cadence"`
	MaxRunningCadenceInStepsPerMinute     float64 `json:"maxRunningCadenceInStepsPerMinute" gorm:"column:max_running_cadence"`
	Steps                                 int     `json:"steps" gorm:"column:steps"`
	UserPro                               bool    `json:"userPro" gorm:"column:user_pro"`
	HasVideo                              bool    `json:"hasVideo" gorm:"column:has_video"`
	TimeZoneID                            int     `json:"timeZoneId" gorm:"column:time_zone_id"`
	BeginTimestamp                        int64   `json:"beginTimestamp" gorm:"column:begin_timestamp"`
	SportTypeID                           int     `json:"sportTypeId" gorm:"column:sport_type_id"`
	AvgPower                              float64 `json:"avgPower" gorm:"column:avg_power"`
	MaxPower                              float64 `json:"maxPower" gorm:"column:max_power"`
	AerobicTrainingEffect                 float64 `json:"aerobicTrainingEffect" gorm:"column:aerobic_training_effect"`
	AnaerobicTrainingEffect               float64 `json:"anaerobicTrainingEffect" gorm:"column:anaerobic_training_effect"`
	NormPower                             float64 `json:"normPower" gorm:"column:norm_power"`
	AvgVerticalOscillation                float64 `json:"avgVerticalOscillation" gorm:"column:avg_vertical_oscillation"`
	AvgGroundContactTime                  float64 `json:"avgGroundContactTime" gorm:"column:avg_ground_contact_time"`
	AvgStrideLength                       float64 `json:"avgStrideLength" gorm:"column:avg_stride_length"`
	VO2MaxValue                           float64 `json:"vO2MaxValue" gorm:"column:vo2_max_value"`
	AvgVerticalRatio                      float64 `json:"avgVerticalRatio" gorm:"column:avg_vertical_ratio"`
	DeviceID                              int64   `json:"deviceId" gorm:"column:device_id"`
	MinTemperature                        float64 `json:"minTemperature" gorm:"column:min_temperature"`
	MaxTemperature                        float64 `json:"maxTemperature" gorm:"column:max_temperature"`
	MinElevation                          float64 `json:"minElevation" gorm:"column:min_elevation;type:decimal(8,6)"`
	MaxElevation                          float64 `json:"maxElevation" gorm:"column:max_elevation"`
	MaxDoubleCadence                      float64 `json:"maxDoubleCadence" gorm:"column:max_double_cadence"`
	MaxVerticalSpeed                      float64 `json:"maxVerticalSpeed" gorm:"column:max_vertical_speed;type:float"`
	Manufacturer                          string  `json:"manufacturer" gorm:"column:manufacturer;size:50"`
	LocationName                          string  `json:"locationName" gorm:"column:location_name;size:100"`
	LapCount                              int     `json:"lapCount" gorm:"column:lap_count"`
	WaterEstimated                        float64 `json:"waterEstimated" gorm:"column:water_estimated"`
	TrainingEffectLabel                   string  `json:"trainingEffectLabel" gorm:"column:training_effect_label;size:50"`
	MinActivityLapDuration                float64 `json:"minActivityLapDuration" gorm:"column:min_activity_lap_duration;type:decimal(8,6)"`
	AerobicTrainingEffectMessage          string  `json:"aerobicTrainingEffectMessage" gorm:"column:aerobic_training_effect_message;size:100"`
	AnaerobicTrainingEffectMessage        string  `json:"anaerobicTrainingEffectMessage" gorm:"column:anaerobic_training_effect_message;size:100"`
	HasSplits                             bool    `json:"hasSplits" gorm:"column:has_splits"`
	ModerateIntensityMinutes              int     `json:"moderateIntensityMinutes" gorm:"column:moderate_intensity_minutes"`
	VigorousIntensityMinutes              int     `json:"vigorousIntensityMinutes" gorm:"column:vigorous_intensity_minutes"`
	AvgGradeAdjustedSpeed                 float64 `json:"avgGradeAdjustedSpeed" gorm:"column:avg_grade_adjusted_speed;type:decimal(8,6)"`
	DifferenceBodyBattery                 int     `json:"differenceBodyBattery" gorm:"column:difference_body_battery"`
	HasHeatMap                            bool    `json:"hasHeatMap" gorm:"column:has_heat_map"`
	FastestSplit1000                      float64 `json:"fastestSplit_1000" gorm:"column:fastest_split_1000;type:decimal(8,6)"`
	FastestSplit1609                      float64 `json:"fastestSplit_1609" gorm:"column:fastest_split_1609;type:decimal(8,6)"`
	HRTimeInZone1                         float64 `json:"hrTimeInZone_1" gorm:"column:hr_time_in_zone_1;type:decimal(8,3)"`
	HRTimeInZone2                         float64 `json:"hrTimeInZone_2" gorm:"column:hr_time_in_zone_2;type:decimal(8,3)"`
	HRTimeInZone3                         float64 `json:"hrTimeInZone_3" gorm:"column:hr_time_in_zone_3;type:decimal(8,3)"`
	HRTimeInZone4                         float64 `json:"hrTimeInZone_4" gorm:"column:hr_time_in_zone_4;type:decimal(8,3)"`
	HRTimeInZone5                         float64 `json:"hrTimeInZone_5" gorm:"column:hr_time_in_zone_5;type:decimal(8,3)"`
	PowerTimeInZone1                      float64 `json:"powerTimeInZone_1" gorm:"column:power_time_in_zone_1;type:decimal(8,3)"`
	PowerTimeInZone2                      float64 `json:"powerTimeInZone_2" gorm:"column:power_time_in_zone_2;type:decimal(8,3)"`
	PowerTimeInZone3                      float64 `json:"powerTimeInZone_3" gorm:"column:power_time_in_zone_3;type:decimal(8,3)"`
	PowerTimeInZone4                      float64 `json:"powerTimeInZone_4" gorm:"column:power_time_in_zone_4;type:decimal(8,3)"`
	PowerTimeInZone5                      float64 `json:"powerTimeInZone_5" gorm:"column:power_time_in_zone_5;type:decimal(8,3)"`
	QualifyingDive                        bool    `json:"qualifyingDive" gorm:"column:qualifying_dive"`
	Parent                                bool    `json:"parent" gorm:"column:parent"`
	PR                                    bool    `json:"pr" gorm:"column:pr"`
	Favorite                              bool    `json:"favorite" gorm:"column:favorite"`
	Purposeful                            bool    `json:"purposeful" gorm:"column:purposeful"`
	DecoDive                              bool    `json:"decoDive" gorm:"column:deco_dive"`
	ManualActivity                        bool    `json:"manualActivity" gorm:"column:manual_activity"`
	AutoCalcCalories                      bool    `json:"autoCalcCalories" gorm:"column:auto_calc_calories"`
	ElevationCorrected                    bool    `json:"elevationCorrected" gorm:"column:elevation_corrected"`
	ATPActivity                           bool    `json:"atpActivity" gorm:"column:atp_activity"`

	// Foreign key columns
	ActivityTypeID int `gorm:"column:activity_type_id;index"`
	EventTypeID    int `gorm:"column:event_type_id;index"`
	PrivacyTypeID  int `gorm:"column:privacy_type_id;index"`

	// Embedded structs (tidak disimpan di database, hanya untuk parsing JSON)
	ActivityType       ActivityType       `json:"activityType" gorm:"-"`
	EventType          EventType          `json:"eventType" gorm:"-"`
	Privacy            Privacy            `json:"privacy" gorm:"-"`
	UserRoles          []string           `json:"userRoles" gorm:"-"`
	SummarizedDiveInfo SummarizedDiveInfo `json:"summarizedDiveInfo" gorm:"-"`
	SplitSummaries     []SplitSummary     `json:"splitSummaries" gorm:"-"`

	// Relationships
	ActivityTypeRef ActivityType   `gorm:"foreignKey:ActivityTypeID;references:TypeID"`
	EventTypeRef    EventType      `gorm:"foreignKey:EventTypeID;references:TypeID"`
	PrivacyRef      Privacy        `gorm:"foreignKey:PrivacyTypeID;references:TypeID"`
	UserRolesList   []UserRole     `gorm:"foreignKey:ActivityID;references:ActivityID"`
	SplitsList      []SplitSummary `gorm:"foreignKey:ActivityID;references:ActivityID"`
}

func (Activity) TableName() string {
	return "activities"
}

// SectionTypeDTO represents the section type information
type SectionTypeDTO struct {
	ID             int    `json:"id" gorm:"column:id;primaryKey"`
	Key            string `json:"key" gorm:"column:key;size:50"`
	SectionTypeKey string `json:"sectionTypeKey" gorm:"column:section_type_key;size:50"`
}

func (SectionTypeDTO) TableName() string {
	return "section_types"
}

type LapDTO struct {
	ID                    string  `json:"id" gorm:"column:id"`
	ActivityID            int64   `json:"activityId" gorm:"column:activity_id;index"`
	StartTimeGMT          string  `json:"startTimeGMT" gorm:"column:start_time_gmt;size:50"`
	StartLatitude         float64 `json:"startLatitude" gorm:"column:start_latitude"`
	StartLongitude        float64 `json:"startLongitude" gorm:"column:start_longitude"`
	EndLatitude           float64 `json:"endLatitude" gorm:"column:end_latitude"`
	EndLongitude          float64 `json:"endLongitude" gorm:"column:end_longitude"`
	Distance              float64 `json:"distance" gorm:"column:distance;type:decimal(10,5)"`
	Duration              float64 `json:"duration" gorm:"column:duration;type:decimal(10,5)"`
	MovingDuration        float64 `json:"movingDuration" gorm:"column:moving_duration;type:decimal(10,5)"`
	ElapsedDuration       float64 `json:"elapsedDuration" gorm:"column:elapsed_duration;type:decimal(10,5)"`
	ElevationGain         float64 `json:"elevationGain" gorm:"column:elevation_gain"`
	ElevationLoss         float64 `json:"elevationLoss" gorm:"column:elevation_loss"`
	MaxElevation          float64 `json:"maxElevation" gorm:"column:max_elevation"`
	MinElevation          float64 `json:"minElevation" gorm:"column:min_elevation"`
	AverageSpeed          float64 `json:"averageSpeed" gorm:"column:average_speed;type:decimal(8,5)"`
	AverageMovingSpeed    float64 `json:"averageMovingSpeed" gorm:"column:average_moving_speed;type:decimal(8,5)"`
	MaxSpeed              float64 `json:"maxSpeed" gorm:"column:max_speed;type:decimal(8,5)"`
	Calories              float64 `json:"calories" gorm:"column:calories"`
	BMRCalories           float64 `json:"bmrCalories" gorm:"column:bmr_calories"`
	AverageHR             float64 `json:"averageHR" gorm:"column:average_hr;type:decimal(5,2)"`
	MaxHR                 float64 `json:"maxHR" gorm:"column:max_hr;type:decimal(5,2)"`
	AverageRunCadence     float64 `json:"averageRunCadence" gorm:"column:average_run_cadence;type:decimal(8,4)"`
	MaxRunCadence         float64 `json:"maxRunCadence" gorm:"column:max_run_cadence"`
	AverageTemperature    float64 `json:"averageTemperature" gorm:"column:average_temperature;type:decimal(5,2)"`
	MaxTemperature        float64 `json:"maxTemperature" gorm:"column:max_temperature;type:decimal(5,2)"`
	MinTemperature        float64 `json:"minTemperature" gorm:"column:min_temperature;type:decimal(5,2)"`
	AveragePower          float64 `json:"averagePower" gorm:"column:average_power"`
	MaxPower              float64 `json:"maxPower" gorm:"column:max_power"`
	MinPower              float64 `json:"minPower" gorm:"column:min_power"`
	NormalizedPower       float64 `json:"normalizedPower" gorm:"column:normalized_power"`
	TotalWork             float64 `json:"totalWork" gorm:"column:total_work;type:decimal(12,6)"`
	GroundContactTime     float64 `json:"groundContactTime" gorm:"column:ground_contact_time;type:decimal(8,6)"`
	StrideLength          float64 `json:"strideLength" gorm:"column:stride_length;type:decimal(8,6)"`
	VerticalOscillation   float64 `json:"verticalOscillation" gorm:"column:vertical_oscillation;type:decimal(8,6)"`
	VerticalRatio         float64 `json:"verticalRatio" gorm:"column:vertical_ratio;type:decimal(8,6)"`
	MaxVerticalSpeed      float64 `json:"maxVerticalSpeed" gorm:"column:max_vertical_speed;type:decimal(8,12)"`
	AvgGradeAdjustedSpeed float64 `json:"avgGradeAdjustedSpeed" gorm:"column:avg_grade_adjusted_speed;type:decimal(8,6)"`
	LapIndex              int     `json:"lapIndex" gorm:"column:lap_index"`
	IntensityType         string  `json:"intensityType" gorm:"column:intensity_type;size:20"`
	MessageIndex          int     `json:"messageIndex" gorm:"column:message_index"`

	// Arrays (tidak disimpan di database, hanya untuk parsing JSON)
	LengthDTOs           []interface{} `json:"lengthDTOs" gorm:"-"`
	ConnectIQMeasurement []interface{} `json:"connectIQMeasurement" gorm:"-"`
}

func (LapDTO) TableName() string {
	return "activity_laps"
}

// ActivitySplitsResponse represents the complete response from Garmin splits API
type ActivitySplitsResponse struct {
	ActivityID int64     `json:"activityId"`
	LapDTOs    []*LapDTO `json:"lapDTOs"`
}

// HeartRateValueDescriptor represents the descriptor for heart rate values
type HeartRateValueDescriptor struct {
	Key   string `json:"key" gorm:"column:key;size:50"`
	Index int    `json:"index" gorm:"column:index"`
}

// HeartRateResponse represents the complete heart rate response from Garmin API
type HeartRate struct {
	ID                               string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK                    int64  `json:"userProfilePK" gorm:"column:user_profile_pk;primaryKey"`
	CalendarDate                     string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	StartTimestampGMT                string `json:"startTimestampGMT" gorm:"column:start_timestamp_gmt;type:timestamp"`
	EndTimestampGMT                  string `json:"endTimestampGMT" gorm:"column:end_timestamp_gmt;type:timestamp"`
	StartTimestampLocal              string `json:"startTimestampLocal" gorm:"column:start_timestamp_local;type:timestamp"`
	EndTimestampLocal                string `json:"endTimestampLocal" gorm:"column:end_timestamp_local;type:timestamp"`
	MaxHeartRate                     int    `json:"maxHeartRate" gorm:"column:max_heart_rate"`
	MinHeartRate                     int    `json:"minHeartRate" gorm:"column:min_heart_rate"`
	RestingHeartRate                 int    `json:"restingHeartRate" gorm:"column:resting_heart_rate"`
	LastSevenDaysAvgRestingHeartRate int    `json:"lastSevenDaysAvgRestingHeartRate" gorm:"column:last_seven_days_avg_resting_heart_rate"`

	// Arrays (tidak disimpan di database, hanya untuk parsing JSON)
	HeartRateValues           [][]int                    `json:"heartRateValues" gorm:"-"`
	HeartRateValueDescriptors []HeartRateValueDescriptor `json:"heartRateValueDescriptors" gorm:"-"`
}

func (HeartRate) TableName() string {
	return "heart_rates"
}

type HeartRateDetail struct {
	HeartRate     int    `json:"heartrate" gorm:"column:heartrate"`
	Timestamp     int    `json:"timestamp" gorm:"column:timestamp"`
	UserProfilePK int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	CalendarDate  string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
}

func (HeartRateDetail) TableName() string {
	return "heart_rate_details"
}

type Step struct {
	ID                    string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK         int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	StartGMT              string `json:"startGMT" gorm:"column:start_gmt;type:timestamp"`
	EndGMT                string `json:"endGMT" gorm:"column:end_gmt;type:timestamp"`
	Steps                 int    `json:"steps" gorm:"column:steps"`
	Pushes                int    `json:"pushes" gorm:"column:pushes"`
	PrimaryActivityLevel  string `json:"primaryActivityLevel" gorm:"column:primary_activity_level;size:50"`
	ActivityLevelConstant bool   `json:"activityLevelConstant" gorm:"column:activity_level_constant"`
}

func (Step) TableName() string {
	return "steps"
}

// StressEvent represents the stress event information
type StressEvent struct {
	ID                     string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK          int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	StressDataID           string `json:"-" gorm:"column:daily_stress_id"`
	EventType              string `json:"eventType" gorm:"column:event_type;size:50"`
	EventStartTimeGmt      string `json:"eventStartTimeGmt" gorm:"column:event_start_time_gmt;type:timestamp"`
	TimezoneOffset         int64  `json:"timezoneOffset" gorm:"column:timezone_offset"`
	DurationInMilliseconds int64  `json:"durationInMilliseconds" gorm:"column:duration_in_milliseconds"`
	BodyBatteryImpact      int    `json:"bodyBatteryImpact" gorm:"column:body_battery_impact"`
	FeedbackType           string `json:"feedbackType" gorm:"column:feedback_type;size:50"`
	ShortFeedback          string `json:"shortFeedback" gorm:"column:short_feedback;size:100"`
}

func (StressEvent) TableName() string {
	return "stress_events"
}

// StressValueDescriptor represents the descriptor for stress values
type StressValueDescriptor struct {
	Key   string `json:"key" gorm:"column:key;size:50"`
	Index int    `json:"index" gorm:"column:index"`
}

// BodyBatteryValueDescriptor represents the descriptor for body battery values
type BodyBatteryValueDescriptor struct {
	BodyBatteryValueDescriptorIndex int    `json:"bodyBatteryValueDescriptorIndex" gorm:"column:body_battery_value_descriptor_index"`
	BodyBatteryValueDescriptorKey   string `json:"bodyBatteryValueDescriptorKey" gorm:"column:body_battery_value_descriptor_key;size:50"`
}

// BodyBatteryVersion represents the version information in body battery values
type BodyBatteryVersion struct {
	Source      string `json:"source" gorm:"column:source;size:50"`
	ParsedValue int    `json:"parsedValue" gorm:"column:parsed_value"`
}

// StressData represents the complete stress data response from Garmin API
type StressData struct {
	ID                                 string                       `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK                      int64                        `json:"userProfilePK" gorm:"column:user_profile_pk"`
	CalendarDate                       string                       `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	ActivityName                       string                       `json:"activityName" gorm:"column:activity_name;size:255"`
	EventStartTimeGmt                  string                       `json:"eventStartTimeGmt" gorm:"column:event_start_time_gmt;type:timestamp"`
	ActivityType                       string                       `json:"activityType" gorm:"column:activity_type;size:100"`
	ActivityID                         int64                        `json:"activityId" gorm:"column:activity_id"`
	AverageStress                      float64                      `json:"averageStress" gorm:"column:average_stress;type:decimal(8,4)"`
	EventType                          string                       `json:"eventType" gorm:"column:event_type;size:50"`
	Event                              StressEvent                  `json:"event" gorm:"-"`
	StressValueDescriptorsDTOList      []StressValueDescriptor      `json:"stressValueDescriptorsDTOList" gorm:"-"`
	StressValuesArray                  [][]int64                    `json:"stressValuesArray" gorm:"-"`
	BodyBatteryValueDescriptorsDTOList []BodyBatteryValueDescriptor `json:"bodyBatteryValueDescriptorsDTOList" gorm:"-"`
	BodyBatteryValuesArray             [][]interface{}              `json:"bodyBatteryValuesArray" gorm:"-"`
}

func (StressData) TableName() string {
	return "user_daily_stress"
}

// StressDetail represents individual stress measurements
type StressDetail struct {
	ID            string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	CalendarDate  string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	Timestamp     int64  `json:"timestamp" gorm:"column:timestamp"`
	StressLevel   int    `json:"stressLevel" gorm:"column:stress_level"`
}

func (StressDetail) TableName() string {
	return "stress_details"
}

// BodyBatteryDetail represents individual body battery measurements
type BodyBatteryDetail struct {
	ID                string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK     int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	CalendarDate      string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	Timestamp         int64  `json:"timestamp" gorm:"column:timestamp"`
	BodyBatteryStatus string `json:"bodyBatteryStatus" gorm:"column:body_battery_status;size:20"`
	BodyBatteryLevel  int    `json:"bodyBatteryLevel" gorm:"column:body_battery_level"`
	VersionSource     string `json:"versionSource" gorm:"column:version_source;size:50"`
	VersionParsed     int    `json:"versionParsed" gorm:"column:version_parsed"`
}

func (BodyBatteryDetail) TableName() string {
	return "body_battery_details"
}
