package domain

type HRVBaseline struct {
	LowUpper      int     `json:"lowUpper" gorm:"column:low_upper"`
	BalancedLow   int     `json:"balancedLow" gorm:"column:balanced_low"`
	BalancedUpper int     `json:"balancedUpper" gorm:"column:balanced_upper"`
	MarkerValue   float64 `json:"markerValue" gorm:"column:marker_value"`
	CalendarDate  string  `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	UserProfilePK int64   `json:"userProfilePK" gorm:"column:user_profile_pk"`
}

func (HRVBaseline) TableName() string {
	return "hrv_baselines"
}

type HRVSummary struct {
	CalendarDate      string      `json:"calendarDate" gorm:"column:calendar_date;type:date"`
	UserProfilePK     int64       `json:"userProfilePK" gorm:"column:user_profile_pk"`
	WeeklyAvg         int         `json:"weeklyAvg" gorm:"column:weekly_avg"`
	LastNightAvg      int         `json:"lastNightAvg" gorm:"column:last_night_avg"`
	LastNight5MinHigh int         `json:"lastNight5MinHigh" gorm:"column:last_night_5min_high"`
	Baseline          HRVBaseline `json:"baseline" gorm:"-"`
	Status            string      `json:"status" gorm:"column:status;size:50"`
	FeedbackPhrase    string      `json:"feedbackPhrase" gorm:"column:feedback_phrase;size:100"`
	CreateTimeStamp   string      `json:"createTimeStamp" gorm:"column:create_time_stamp;type:timestamp"`
}

func (HRVSummary) TableName() string {
	return "hrv_summaries"
}

type HRVReading struct {
	ParentID         string `json:"parent_id" gorm:"column:parent_id"`
	UserProfilePK    int64  `json:"userProfilePK" gorm:"column:user_profile_pk"`
	HRVValue         int    `json:"hrvValue" gorm:"column:hrv_value"`
	ReadingTimeGMT   string `json:"readingTimeGMT" gorm:"column:reading_time_gmt;type:timestamp"`
	ReadingTimeLocal string `json:"readingTimeLocal" gorm:"column:reading_time_local;type:timestamp"`
}

func (HRVReading) TableName() string {
	return "hrv_readings"
}

type HRVData struct {
	ID                       string       `json:"id" gorm:"column:id;primaryKey"`
	UserProfilePK            int64        `json:"userProfilePk" gorm:"column:user_profile_pk"`
	HRVSummary               HRVSummary   `json:"hrvSummary" gorm:"-"`
	StartTimestampGMT        string       `json:"startTimestampGMT" gorm:"column:start_timestamp_gmt;type:timestamp"`
	EndTimestampGMT          string       `json:"endTimestampGMT" gorm:"column:end_timestamp_gmt;type:timestamp"`
	StartTimestampLocal      string       `json:"startTimestampLocal" gorm:"column:start_timestamp_local;type:timestamp"`
	EndTimestampLocal        string       `json:"endTimestampLocal" gorm:"column:end_timestamp_local;type:timestamp"`
	SleepStartTimestampGMT   string       `json:"sleepStartTimestampGMT" gorm:"column:sleep_start_timestamp_gmt;type:timestamp"`
	SleepEndTimestampGMT     string       `json:"sleepEndTimestampGMT" gorm:"column:sleep_end_timestamp_gmt;type:timestamp"`
	SleepStartTimestampLocal string       `json:"sleepStartTimestampLocal" gorm:"column:sleep_start_timestamp_local;type:timestamp"`
	SleepEndTimestampLocal   string       `json:"sleepEndTimestampLocal" gorm:"column:sleep_end_timestamp_local;type:timestamp"`
	HRVReadings              []HRVReading `json:"hrvReadings" gorm:"-"`
}

func (HRVData) TableName() string {
	return "hrv_data"
}
