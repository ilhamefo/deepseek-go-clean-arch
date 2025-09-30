package domain

import "github.com/lib/pq"

type Weight struct {
	Source      string  `json:"source" gorm:"column:source;size:50"`
	ParsedValue float64 `json:"parsedValue" gorm:"column:parsed_value"`
}

type Height struct {
	Source      string  `json:"source" gorm:"column:source;size:50"`
	ParsedValue float64 `json:"parsedValue" gorm:"column:parsed_value"`
}

type PowerFormat struct {
	ID            string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK int64  `json:"userProfilePK" gorm:"column:user_profile_pk;index"`
	Format
}

func (PowerFormat) TableName() string {
	return "user_power_formats"
}

type HeartRateFormat struct {
	ID            string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK int64  `json:"userProfilePK" gorm:"column:user_profile_pk;index"`
	Format
}

func (HeartRateFormat) TableName() string {
	return "user_heart_rate_formats"
}

type Format struct {
	FormatID      int     `json:"formatId" gorm:"column:format_id"`
	FormatKey     string  `json:"formatKey" gorm:"column:format_key;size:50"`
	MinFraction   int     `json:"minFraction" gorm:"column:min_fraction"`
	MaxFraction   int     `json:"maxFraction" gorm:"column:max_fraction"`
	GroupingUsed  bool    `json:"groupingUsed" gorm:"column:grouping_used"`
	DisplayFormat *string `json:"displayFormat" gorm:"column:display_format;size:100"`
}

type FirstDayOfWeek struct {
	DayID              int    `json:"dayId" gorm:"column:day_id"`
	DayName            string `json:"dayName" gorm:"column:day_name;size:20"`
	SortOrder          int    `json:"sortOrder" gorm:"column:sort_order"`
	IsPossibleFirstDay bool   `json:"isPossibleFirstDay" gorm:"column:is_possible_first_day"`
}

type VO2Max struct {
	Source      string  `json:"source" gorm:"column:source;size:50"`
	ParsedValue float64 `json:"parsedValue" gorm:"column:parsed_value"`
}

type HydrationContainer struct {
	ID            string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK int64  `json:"userProfilePK" gorm:"column:user_profile_pk;index"`
	Name          string `json:"name" gorm:"column:name;size:100"`
	Volume        int    `json:"volume" gorm:"column:volume"`
	Unit          string `json:"unit" gorm:"column:unit;size:20"`
}

func (HydrationContainer) TableName() string {
	return "hydration_containers"
}

type WeatherLocation struct {
	UseFixedLocation *bool    `json:"useFixedLocation" gorm:"column:use_fixed_location"`
	Latitude         *float64 `json:"latitude" gorm:"column:latitude"`
	Longitude        *float64 `json:"longitude" gorm:"column:longitude"`
	LocationName     *string  `json:"locationName" gorm:"column:location_name;size:100"`
	IsoCountryCode   *string  `json:"isoCountryCode" gorm:"column:iso_country_code;size:10"`
	PostalCode       *string  `json:"postalCode" gorm:"column:postal_code;size:20"`
}

type UserData struct {
	ID                             string   `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK                  int64    `json:"userProfilePK" gorm:"column:user_profile_pk;index"`
	Gender                         string   `json:"gender" gorm:"column:gender;size:10"`
	TimeFormat                     string   `json:"timeFormat" gorm:"column:time_format;size:50"`
	BirthDate                      string   `json:"birthDate" gorm:"column:birth_date;type:date"`
	MeasurementSystem              string   `json:"measurementSystem" gorm:"column:measurement_system;size:20"`
	ActivityLevel                  *string  `json:"activityLevel" gorm:"column:activity_level;size:50"`
	Handedness                     string   `json:"handedness" gorm:"column:handedness;size:10"`
	IntensityMinutesCalcMethod     string   `json:"intensityMinutesCalcMethod" gorm:"column:intensity_minutes_calc_method;size:20"`
	ModerateIntensityMinutesHrZone int      `json:"moderateIntensityMinutesHrZone" gorm:"column:moderate_intensity_minutes_hr_zone"`
	VigorousIntensityMinutesHrZone int      `json:"vigorousIntensityMinutesHrZone" gorm:"column:vigorous_intensity_minutes_hr_zone"`
	HydrationMeasurementUnit       string   `json:"hydrationMeasurementUnit" gorm:"column:hydration_measurement_unit;size:20"`
	HydrationAutoGoalEnabled       bool     `json:"hydrationAutoGoalEnabled" gorm:"column:hydration_auto_goal_enabled"`
	FirstbeatMaxStressScore        *float64 `json:"firstbeatMaxStressScore" gorm:"column:firstbeat_max_stress_score"`
	FirstbeatCyclingLtTimestamp    *string  `json:"firstbeatCyclingLtTimestamp" gorm:"column:firstbeat_cycling_lt_timestamp;type:timestamp"`
	FirstbeatRunningLtTimestamp    *string  `json:"firstbeatRunningLtTimestamp" gorm:"column:firstbeat_running_lt_timestamp;type:timestamp"`
	ThresholdHeartRateAutoDetected *bool    `json:"thresholdHeartRateAutoDetected" gorm:"column:threshold_heart_rate_auto_detected"`
	FtpAutoDetected                bool     `json:"ftpAutoDetected" gorm:"column:ftp_auto_detected"`
	TrainingStatusPausedDate       *string  `json:"trainingStatusPausedDate" gorm:"column:training_status_paused_date;type:date"`
	GolfDistanceUnit               string   `json:"golfDistanceUnit" gorm:"column:golf_distance_unit;size:20"`
	GolfElevationUnit              *string  `json:"golfElevationUnit" gorm:"column:golf_elevation_unit;size:20"`
	GolfSpeedUnit                  *string  `json:"golfSpeedUnit" gorm:"column:golf_speed_unit;size:20"`
	ExternalBottomTime             *string  `json:"externalBottomTime" gorm:"column:external_bottom_time;type:timestamp"`
	VirtualCaddieDataSource        *string  `json:"virtualCaddieDataSource" gorm:"column:virtual_caddie_data_source;size:50"`
	NumberDivesAutomatically       *bool    `json:"numberDivesAutomatically" gorm:"column:number_dives_automatically"`
	DiveNumber                     *int     `json:"diveNumber" gorm:"column:dive_number"`
	LactateThresholdSpeed          *float64 `json:"lactateThresholdSpeed" gorm:"column:lactate_threshold_speed"`
	LactateThresholdHeartRate      *float64 `json:"lactateThresholdHeartRate" gorm:"column:lactate_threshold_heart_rate"`
	Weight                         *float32 `json:"weight" gorm:"weight"`
	Height                         *float32 `json:"height" gorm:"height"`
	VO2MaxRunning                  *float32 `json:"vo2MaxRunning" gorm:"vo2_max_running"`
	VO2MaxCycling                  *float32 `json:"vo2MaxCycling" gorm:"vo2_max_cycling"`

	PowerFormat               PowerFormat          `json:"powerFormat" gorm:"-"`
	HeartRateFormat           HeartRateFormat      `json:"heartRateFormat" gorm:"-"`
	FirstDayOfWeek            FirstDayOfWeek       `json:"firstDayOfWeek" gorm:"-"`
	HydrationContainers       []HydrationContainer `json:"hydrationContainers" gorm:"-"`
	WeatherLocation           WeatherLocation      `json:"weatherLocation" gorm:"-"`
	AvailableTrainingDays     []string             `json:"availableTrainingDays" gorm:"-"`
	PreferredLongTrainingDays []string             `json:"preferredLongTrainingDays" gorm:"-"`
}

func (UserData) TableName() string {
	return "user_data"
}

type PreferredLongTrainingDays struct {
	ID            string         `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK int64          `json:"userProfilePK" gorm:"column:user_profile_pk;index"`
	Days          pq.StringArray `json:"days" gorm:"column:days;type:varchar[]"`
}

func (PreferredLongTrainingDays) TableName() string {
	return "user_preferred_long_training_days"
}

type AvailableTrainingDays struct {
	ID            string         `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK int64          `json:"userProfilePK" gorm:"column:user_profile_pk;index"`
	Days          pq.StringArray `json:"days" gorm:"column:days;type:varchar[]"`
}

func (AvailableTrainingDays) TableName() string {
	return "user_available_training_days"
}

type UserSleep struct {
	ID               string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK    int64  `json:"userProfilePK" gorm:"column:user_profile_pk;index"`
	SleepTime        int    `json:"sleepTime" gorm:"column:sleep_time"`
	DefaultSleepTime bool   `json:"defaultSleepTime" gorm:"column:default_sleep_time"`
	WakeTime         int    `json:"wakeTime" gorm:"column:wake_time"`
	DefaultWakeTime  bool   `json:"defaultWakeTime" gorm:"column:default_wake_time"`
}

func (UserSleep) TableName() string {
	return "user_sleep"
}

type UserSleepWindow struct {
	ID                                string `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK                     int64  `json:"userProfilePK" gorm:"column:user_profile_pk;index"`
	SleepWindowFrequency              string `json:"sleepWindowFrequency" gorm:"column:sleep_window_frequency;size:20"`
	StartSleepTimeSecondsFromMidnight int    `json:"startSleepTimeSecondsFromMidnight" gorm:"column:start_sleep_time_seconds_from_midnight"`
	EndSleepTimeSecondsFromMidnight   int    `json:"endSleepTimeSecondsFromMidnight" gorm:"column:end_sleep_time_seconds_from_midnight"`
}

func (UserSleepWindow) TableName() string {
	return "user_sleep_windows"
}

// UserSetting represents the complete user profile response from Garmin API
type UserSetting struct {
	ID               int64             `json:"id" gorm:"column:id;primaryKey"`
	ConnectDate      *string           `json:"connectDate" gorm:"column:connect_date;type:timestamp"`
	SourceType       *string           `json:"sourceType" gorm:"column:source_type;size:50"`
	UserData         UserData          `json:"userData" gorm:"-"`
	UserSleep        UserSleep         `json:"userSleep" gorm:"-"`
	UserSleepWindows []UserSleepWindow `json:"userSleepWindows" gorm:"-"`

	// UserDataRef          UserData          `gorm:"foreignKey:UserProfilePK;references:ID"`
	// UserSleepRef         UserSleep         `gorm:"foreignKey:UserProfilePK;references:ID"`
	// UserSleepWindowsList []UserSleepWindow `gorm:"foreignKey:UserProfilePK;references:ID"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}
