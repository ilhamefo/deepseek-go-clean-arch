package domain

// SleepResponse represents the complete Garmin sleep data response
type SleepResponse struct {
	DailySleepDTO                        DailySleepDTO                     `json:"dailySleepDTO" gorm:"-"`
	SleepMovement                        []SleepMovement                   `json:"sleepMovement" gorm:"-"`
	RemSleepData                         bool                              `json:"remSleepData" gorm:"-"`
	SleepLevels                          []SleepLevel                      `json:"sleepLevels" gorm:"-"`
	SleepRestlessMoments                 []SleepRestlessMoment             `json:"sleepRestlessMoments" gorm:"-"`
	RestlessMomentsCount                 int                               `json:"restlessMomentsCount" gorm:"-"`
	WellnessSpO2SleepSummaryDTO          WellnessSpO2SleepSummaryDTO       `json:"wellnessSpO2SleepSummaryDTO" gorm:"-"`
	WellnessEpochSPO2DataDTOList         []WellnessEpochSPO2DataDTO        `json:"wellnessEpochSPO2DataDTOList" gorm:"-"`
	WellnessEpochRespirationDataDTOList  []WellnessEpochRespirationDataDTO `json:"wellnessEpochRespirationDataDTOList" gorm:"-"`
	WellnessEpochRespirationAveragesList []interface{}                     `json:"wellnessEpochRespirationAveragesList" gorm:"-"`
	RespirationVersion                   int                               `json:"respirationVersion" gorm:"-"`
	SleepHeartRate                       []SleepHeartRate                  `json:"sleepHeartRate" gorm:"-"`
	SleepStress                          []SleepStress                     `json:"sleepStress" gorm:"-"`
	SleepBodyBattery                     []SleepBodyBattery                `json:"sleepBodyBattery" gorm:"-"`
	SkinTempDataExists                   bool                              `json:"skinTempDataExists" gorm:"-"`
	HrvData                              []HrvData                         `json:"hrvData" gorm:"-"`
	BreathingDisruptionData              []BreathingDisruptionData         `json:"breathingDisruptionData" gorm:"-"`
	AvgOvernightHrv                      float64                           `json:"avgOvernightHrv" gorm:"-"`
	HrvStatus                            string                            `json:"hrvStatus" gorm:"-"`
	BodyBatteryChange                    int                               `json:"bodyBatteryChange" gorm:"-"`
	RestingHeartRate                     int                               `json:"restingHeartRate" gorm:"-"`
}

// DailySleepDTO represents the main sleep data
type DailySleepDTO struct {
	ID                            int64       `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK                 int64       `json:"userProfilePK" gorm:"column:user_profile_pk"`
	CalendarDate                  string      `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	SleepTimeSeconds              int         `json:"sleepTimeSeconds" gorm:"column:sleep_time_seconds"`
	NapTimeSeconds                int         `json:"napTimeSeconds" gorm:"column:nap_time_seconds"`
	SleepWindowConfirmed          bool        `json:"sleepWindowConfirmed" gorm:"column:sleep_window_confirmed"`
	SleepWindowConfirmationType   string      `json:"sleepWindowConfirmationType" gorm:"column:sleep_window_confirmation_type;size:100"`
	SleepStartTimestampGMT        *int64      `json:"sleepStartTimestampGMT" gorm:"column:sleep_start_timestamp_gmt"`
	SleepEndTimestampGMT          *int64      `json:"sleepEndTimestampGMT" gorm:"column:sleep_end_timestamp_gmt"`
	SleepStartTimestampLocal      *int64      `json:"sleepStartTimestampLocal" gorm:"column:sleep_start_timestamp_local"`
	SleepEndTimestampLocal        *int64      `json:"sleepEndTimestampLocal" gorm:"column:sleep_end_timestamp_local"`
	AutoSleepStartTimestampGMT    *int64      `json:"autoSleepStartTimestampGMT" gorm:"column:auto_sleep_start_timestamp_gmt"`
	AutoSleepEndTimestampGMT      *int64      `json:"autoSleepEndTimestampGMT" gorm:"column:auto_sleep_end_timestamp_gmt"`
	SleepQualityTypePK            *int        `json:"sleepQualityTypePK" gorm:"column:sleep_quality_type_pk"`
	SleepResultTypePK             *int        `json:"sleepResultTypePK" gorm:"column:sleep_result_type_pk"`
	UnmeasurableSleepSeconds      int         `json:"unmeasurableSleepSeconds" gorm:"column:unmeasurable_sleep_seconds"`
	DeepSleepSeconds              int         `json:"deepSleepSeconds" gorm:"column:deep_sleep_seconds"`
	LightSleepSeconds             int         `json:"lightSleepSeconds" gorm:"column:light_sleep_seconds"`
	RemSleepSeconds               int         `json:"remSleepSeconds" gorm:"column:rem_sleep_seconds"`
	AwakeSleepSeconds             int         `json:"awakeSleepSeconds" gorm:"column:awake_sleep_seconds"`
	DeviceRemCapable              bool        `json:"deviceRemCapable" gorm:"column:device_rem_capable"`
	Retro                         bool        `json:"retro" gorm:"column:retro"`
	SleepFromDevice               bool        `json:"sleepFromDevice" gorm:"column:sleep_from_device"`
	AverageSpO2Value              float64     `json:"averageSpO2Value" gorm:"embedded;embeddedPrefix:avg_spo2_"`
	LowestSpO2Value               int         `json:"lowestSpO2Value" gorm:"column:lowest_spo2_value"`
	HighestSpO2Value              int         `json:"highestSpO2Value" gorm:"column:highest_spo2_value"`
	AverageSpO2HRSleep            float64     `json:"averageSpO2HRSleep" gorm:"embedded;embeddedPrefix:avg_spo2_hr_"`
	AverageRespirationValue       float64     `json:"averageRespirationValue" gorm:"embedded;embeddedPrefix:avg_respiration_"`
	LowestRespirationValue        float64     `json:"lowestRespirationValue" gorm:"embedded;embeddedPrefix:lowest_respiration_"`
	HighestRespirationValue       float64     `json:"highestRespirationValue" gorm:"embedded;embeddedPrefix:highest_respiration_"`
	AwakeCount                    int         `json:"awakeCount" gorm:"column:awake_count"`
	AvgSleepStress                float64     `json:"avgSleepStress" gorm:"embedded;embeddedPrefix:avg_sleep_stress_"`
	AgeGroup                      string      `json:"ageGroup" gorm:"column:age_group;size:20"`
	AvgHeartRate                  float64     `json:"avgHeartRate" gorm:"embedded;embeddedPrefix:avg_heart_rate_"`
	SleepScoreFeedback            string      `json:"sleepScoreFeedback" gorm:"column:sleep_score_feedback;size:50"`
	SleepScoreInsight             string      `json:"sleepScoreInsight" gorm:"column:sleep_score_insight;size:50"`
	SleepScorePersonalizedInsight string      `json:"sleepScorePersonalizedInsight" gorm:"column:sleep_score_personalized_insight;size:100"`
	SleepScores                   SleepScores `json:"sleepScores" gorm:"-"`
	SleepVersion                  int         `json:"sleepVersion" gorm:"column:sleep_version"`
	SleepNeed                     SleepNeed   `json:"sleepNeed" gorm:"-"`
	NextSleepNeed                 SleepNeed   `json:"nextSleepNeed" gorm:"-"`
	BreathingDisruptionSeverity   string      `json:"breathingDisruptionSeverity" gorm:"column:breathing_disruption_severity;size:20"`
}

func (DailySleepDTO) TableName() string {
	return "daily_sleep"
}

// SleepScores represents sleep scoring metrics
type SleepScores struct {
	ID              string     `json:"id" gorm:"column:id;primaryKey"`
	SleepID         int64      `json:"sleepId" gorm:"column:sleep_id;index"`
	TotalDuration   SleepScore `json:"totalDuration" gorm:"-"`
	Stress          SleepScore `json:"stress" gorm:"-"`
	AwakeCount      SleepScore `json:"awakeCount" gorm:"-"`
	Overall         Overall    `json:"overall" gorm:"embedded;embeddedPrefix:overall_"`
	RemPercentage   SleepScore `json:"remPercentage" gorm:"-"`
	Restlessness    SleepScore `json:"restlessness" gorm:"-"`
	LightPercentage SleepScore `json:"lightPercentage" gorm:"-"`
	DeepPercentage  SleepScore `json:"deepPercentage" gorm:"-"`
}

func (SleepScores) TableName() string {
	return "sleep_scores"
}

// SleepScore represents individual sleep score metrics
type SleepScore struct {
	ID                  string   `json:"id" gorm:"column:id;primaryKey"`
	SleepScoresID       string   `json:"sleepScoresId" gorm:"column:sleep_scores_id;index"`
	ScoreType           string   `json:"scoreType" gorm:"column:score_type;size:50"`
	Value               *int     `json:"value,omitempty" gorm:"column:value"`
	QualifierKey        string   `json:"qualifierKey" gorm:"column:qualifier_key;size:50"`
	OptimalStart        *float64 `json:"optimalStart,omitempty" gorm:"embedded;embeddedPrefix:optimal_start_"`
	OptimalEnd          *float64 `json:"optimalEnd,omitempty" gorm:"embedded;embeddedPrefix:optimal_end_"`
	IdealStartInSeconds *float64 `json:"idealStartInSeconds,omitempty" gorm:"column:ideal_start_in_seconds;type:decimal(10,2)"`
	IdealEndInSeconds   *float64 `json:"idealEndInSeconds,omitempty" gorm:"column:ideal_end_in_seconds;type:decimal(10,2)"`
}

func (SleepScore) TableName() string {
	return "sleep_score_details"
}

// Overall represents overall sleep score
type Overall struct {
	Value        int    `json:"value" gorm:"column:value"`
	QualifierKey string `json:"qualifierKey" gorm:"column:qualifier_key;size:50"`
}

// SleepNeed represents sleep need information
type SleepNeed struct {
	ID                       string `json:"id" gorm:"column:id;primaryKey"`
	SleepID                  int64  `json:"sleepId" gorm:"column:sleep_id;index"`
	UserProfilePk            int64  `json:"userProfilePk" gorm:"column:user_profile_pk"`
	CalendarDate             string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	DeviceId                 int64  `json:"deviceId" gorm:"column:device_id"`
	TimestampGmt             string `json:"timestampGmt" gorm:"column:timestamp_gmt;type:timestamp"`
	Baseline                 int    `json:"baseline" gorm:"column:baseline"`
	Actual                   int    `json:"actual" gorm:"column:actual"`
	Feedback                 string `json:"feedback" gorm:"column:feedback;size:50"`
	TrainingFeedback         string `json:"trainingFeedback" gorm:"column:training_feedback;size:50"`
	SleepHistoryAdjustment   string `json:"sleepHistoryAdjustment" gorm:"column:sleep_history_adjustment;size:50"`
	HrvAdjustment            string `json:"hrvAdjustment" gorm:"column:hrv_adjustment;size:50"`
	NapAdjustment            string `json:"napAdjustment" gorm:"column:nap_adjustment;size:50"`
	DisplayedForTheDay       bool   `json:"displayedForTheDay" gorm:"column:displayed_for_the_day"`
	PreferredActivityTracker bool   `json:"preferredActivityTracker" gorm:"column:preferred_activity_tracker"`
	NeedType                 string `json:"needType" gorm:"column:need_type;size:20"` // untuk membedakan current dan next
}

func (SleepNeed) TableName() string {
	return "sleep_needs"
}

// SleepMovement represents sleep movement data
type SleepMovement struct {
	ID            string  `json:"id" gorm:"column:id;primaryKey"`
	SleepID       int64   `json:"sleepId" gorm:"column:sleep_id;index"`
	StartGMT      string  `json:"startGMT" gorm:"column:start_gmt;type:timestamp"`
	EndGMT        string  `json:"endGMT" gorm:"column:end_gmt;type:timestamp"`
	ActivityLevel float64 `json:"activityLevel" gorm:"column:activity_level;type:decimal(15,12)"`
}

func (SleepMovement) TableName() string {
	return "sleep_movements"
}

// SleepLevel represents sleep level data
type SleepLevel struct {
	ID            string  `json:"id" gorm:"column:id;primaryKey"`
	SleepID       int64   `json:"sleepId" gorm:"column:sleep_id;index"`
	StartGMT      string  `json:"startGMT" gorm:"column:start_gmt;type:timestamp"`
	EndGMT        string  `json:"endGMT" gorm:"column:end_gmt;type:timestamp"`
	ActivityLevel float64 `json:"activityLevel" gorm:"embedded;embeddedPrefix:activity_level_"`
}

func (SleepLevel) TableName() string {
	return "sleep_levels"
}

// SleepRestlessMoment represents restless moments during sleep
type SleepRestlessMoment struct {
	ID       string `json:"id" gorm:"column:id;primaryKey"`
	SleepID  int64  `json:"sleepId" gorm:"column:sleep_id;index"`
	StartGMT int64  `json:"startGMT" gorm:"column:start_gmt;type:timestamp"`
	EndGMT   string `json:"endGMT" gorm:"column:end_gmt;type:timestamp"`
}

func (SleepRestlessMoment) TableName() string {
	return "sleep_restless_moments"
}

// WellnessSpO2SleepSummaryDTO represents SpO2 summary during sleep
type WellnessSpO2SleepSummaryDTO struct {
	ID                             string  `json:"id" gorm:"column:id;primaryKey"`
	SleepID                        int64   `json:"sleepId" gorm:"column:sleep_id;index"`
	UserProfilePk                  int64   `json:"userProfilePk" gorm:"column:user_profile_pk"`
	DeviceId                       int64   `json:"deviceId" gorm:"column:device_id"`
	SleepMeasurementStartGMT       string  `json:"sleepMeasurementStartGMT" gorm:"column:sleep_measurement_start_gmt;type:timestamp"`
	SleepMeasurementEndGMT         string  `json:"sleepMeasurementEndGMT" gorm:"column:sleep_measurement_end_gmt;type:timestamp"`
	AlertThresholdValue            *int    `json:"alertThresholdValue" gorm:"column:alert_threshold_value"`
	NumberOfEventsBelowThreshold   *int    `json:"numberOfEventsBelowThreshold" gorm:"column:number_of_events_below_threshold"`
	DurationOfEventsBelowThreshold *int    `json:"durationOfEventsBelowThreshold" gorm:"column:duration_of_events_below_threshold"`
	AverageSPO2                    float64 `json:"averageSPO2" gorm:"-"`
	AverageSpO2HR                  float64 `json:"averageSpO2HR" gorm:"-"`
	LowestSPO2                     int     `json:"lowestSPO2" gorm:"column:lowest_spo2"`
}

func (WellnessSpO2SleepSummaryDTO) TableName() string {
	return "sleep_spo2_summary"
}

// WellnessEpochSPO2DataDTO represents individual SpO2 readings
type WellnessEpochSPO2DataDTO struct {
	ID                string `json:"id" gorm:"column:id;primaryKey"`
	SleepID           int64  `json:"sleepId" gorm:"column:sleep_id;index"`
	UserProfilePK     int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	EpochTimestamp    string `json:"epochTimestamp" gorm:"column:epoch_timestamp;type:timestamp"`
	DeviceId          int64  `json:"deviceId" gorm:"column:device_id"`
	CalendarDate      string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	EpochDuration     int    `json:"epochDuration" gorm:"column:epoch_duration"`
	Spo2Reading       int    `json:"spo2Reading" gorm:"column:spo2_reading"`
	ReadingConfidence int    `json:"readingConfidence" gorm:"column:reading_confidence"`
}

func (WellnessEpochSPO2DataDTO) TableName() string {
	return "sleep_spo2_data"
}

// WellnessEpochRespirationDataDTO represents individual respiration readings
type WellnessEpochRespirationDataDTO struct {
	ID                 string `json:"id" gorm:"column:id;primaryKey"`
	SleepID            int64  `json:"sleepId" gorm:"column:sleep_id;index"`
	UserProfilePK      int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	EpochTimestamp     string `json:"epochTimestamp" gorm:"column:epoch_timestamp;type:timestamp"`
	DeviceId           int64  `json:"deviceId" gorm:"column:device_id"`
	CalendarDate       string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	EpochDuration      int    `json:"epochDuration" gorm:"column:epoch_duration"`
	RespirationReading int    `json:"respirationReading" gorm:"column:respiration_reading"`
	ReadingConfidence  int    `json:"readingConfidence" gorm:"column:reading_confidence"`
}

func (WellnessEpochRespirationDataDTO) TableName() string {
	return "sleep_respiration_data"
}

// SleepHeartRate represents heart rate data during sleep
type SleepHeartRate struct {
	ID                string `json:"id" gorm:"column:id;primaryKey"`
	SleepID           int64  `json:"sleepId" gorm:"column:sleep_id;index"`
	UserProfilePK     int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	EpochTimestamp    string `json:"epochTimestamp" gorm:"column:epoch_timestamp;type:timestamp"`
	DeviceId          int64  `json:"deviceId" gorm:"column:device_id"`
	CalendarDate      string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	EpochDuration     int    `json:"epochDuration" gorm:"column:epoch_duration"`
	HeartRateReading  int    `json:"heartRateReading" gorm:"column:heart_rate_reading"`
	ReadingConfidence int    `json:"readingConfidence" gorm:"column:reading_confidence"`
}

func (SleepHeartRate) TableName() string {
	return "sleep_heart_rate"
}

// SleepStress represents stress data during sleep
type SleepStress struct {
	ID                string `json:"id" gorm:"column:id;primaryKey"`
	SleepID           int64  `json:"sleepId" gorm:"column:sleep_id;index"`
	UserProfilePK     int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	EpochTimestamp    string `json:"epochTimestamp" gorm:"column:epoch_timestamp;type:timestamp"`
	DeviceId          int64  `json:"deviceId" gorm:"column:device_id"`
	CalendarDate      string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	EpochDuration     int    `json:"epochDuration" gorm:"column:epoch_duration"`
	StressReading     int    `json:"stressReading" gorm:"column:stress_reading"`
	ReadingConfidence int    `json:"readingConfidence" gorm:"column:reading_confidence"`
}

func (SleepStress) TableName() string {
	return "sleep_stress"
}

// SleepBodyBattery represents body battery data during sleep
type SleepBodyBattery struct {
	ID                 string `json:"id" gorm:"column:id;primaryKey"`
	SleepID            int64  `json:"sleepId" gorm:"column:sleep_id;index"`
	UserProfilePK      int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	EpochTimestamp     string `json:"epochTimestamp" gorm:"column:epoch_timestamp;type:timestamp"`
	DeviceId           int64  `json:"deviceId" gorm:"column:device_id"`
	CalendarDate       string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	EpochDuration      int    `json:"epochDuration" gorm:"column:epoch_duration"`
	BodyBatteryReading int    `json:"bodyBatteryReading" gorm:"column:body_battery_reading"`
	ReadingConfidence  int    `json:"readingConfidence" gorm:"column:reading_confidence"`
}

func (SleepBodyBattery) TableName() string {
	return "sleep_body_battery"
}

// HrvData represents heart rate variability data
type HrvData struct {
	ID                string `json:"id" gorm:"column:id;primaryKey"`
	SleepID           int64  `json:"sleepId" gorm:"column:sleep_id;index"`
	UserProfilePK     int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	EpochTimestamp    string `json:"epochTimestamp" gorm:"column:epoch_timestamp;type:timestamp"`
	DeviceId          int64  `json:"deviceId" gorm:"column:device_id"`
	CalendarDate      string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	EpochDuration     int    `json:"epochDuration" gorm:"column:epoch_duration"`
	HrvReading        int    `json:"hrvReading" gorm:"column:hrv_reading"`
	ReadingConfidence int    `json:"readingConfidence" gorm:"column:reading_confidence"`
}

func (HrvData) TableName() string {
	return "sleep_hrv_data"
}

// BreathingDisruptionData represents breathing disruption data
type BreathingDisruptionData struct {
	ID                         string `json:"id" gorm:"column:id;primaryKey"`
	SleepID                    int64  `json:"sleepId" gorm:"column:sleep_id;index"`
	UserProfilePK              int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	EpochTimestamp             string `json:"epochTimestamp" gorm:"column:epoch_timestamp;type:timestamp"`
	DeviceId                   int64  `json:"deviceId" gorm:"column:device_id"`
	CalendarDate               string `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	EpochDuration              int    `json:"epochDuration" gorm:"column:epoch_duration"`
	BreathingDisruptionReading int    `json:"breathingDisruptionReading" gorm:"column:breathing_disruption_reading"`
	ReadingConfidence          int    `json:"readingConfidence" gorm:"column:reading_confidence"`
}

func (BreathingDisruptionData) TableName() string {
	return "sleep_breathing_disruption"
}
