-- PostgreSQL script to create Garmin Sleep tracking tables
-- Generated from garmin_sleep.go domain models

-- Create database if not exists (uncomment if needed)
-- CREATE DATABASE garmin_data;

-- Main sleep data table
CREATE TABLE IF NOT EXISTS daily_sleep (
    id BIGINT PRIMARY KEY,
    user_profile_pk BIGINT NOT NULL,
    calendar_date DATE NOT NULL,
    sleep_time_seconds INTEGER NOT NULL DEFAULT 0,
    nap_time_seconds INTEGER NOT NULL DEFAULT 0,
    sleep_window_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    sleep_window_confirmation_type VARCHAR(100),
    sleep_start_timestamp_gmt BIGINT,
    sleep_end_timestamp_gmt BIGINT,
    sleep_start_timestamp_local BIGINT,
    sleep_end_timestamp_local BIGINT,
    auto_sleep_start_timestamp_gmt BIGINT,
    auto_sleep_end_timestamp_gmt BIGINT,
    sleep_quality_type_pk INTEGER,
    sleep_result_type_pk INTEGER,
    unmeasurable_sleep_seconds INTEGER NOT NULL DEFAULT 0,
    deep_sleep_seconds INTEGER NOT NULL DEFAULT 0,
    light_sleep_seconds INTEGER NOT NULL DEFAULT 0,
    rem_sleep_seconds INTEGER NOT NULL DEFAULT 0,
    awake_sleep_seconds INTEGER NOT NULL DEFAULT 0,
    device_rem_capable BOOLEAN NOT NULL DEFAULT FALSE,
    retro BOOLEAN NOT NULL DEFAULT FALSE,
    sleep_from_device BOOLEAN NOT NULL DEFAULT FALSE,
    -- Average SpO2 Value (embedded ParsedValue)
    avg_spo2_source VARCHAR(20),
    avg_spo2_parsed_value INTEGER,
    lowest_spo2_value INTEGER,
    highest_spo2_value INTEGER,
    -- Average SpO2 HR Sleep (embedded ParsedValue)
    avg_spo2_hr_source VARCHAR(20),
    avg_spo2_hr_parsed_value INTEGER,
    -- Average Respiration Value (embedded ParsedValue)
    avg_respiration_source VARCHAR(20),
    avg_respiration_parsed_value INTEGER,
    -- Lowest Respiration Value (embedded ParsedValue)
    lowest_respiration_source VARCHAR(20),
    lowest_respiration_parsed_value INTEGER,
    -- Highest Respiration Value (embedded ParsedValue)
    highest_respiration_source VARCHAR(20),
    highest_respiration_parsed_value INTEGER,
    awake_count INTEGER NOT NULL DEFAULT 0,
    -- Average Sleep Stress (embedded ParsedValue)
    avg_sleep_stress_source VARCHAR(20),
    avg_sleep_stress_parsed_value INTEGER,
    age_group VARCHAR(20),
    -- Average Heart Rate (embedded ParsedValue)
    avg_heart_rate_source VARCHAR(20),
    avg_heart_rate_parsed_value INTEGER,
    sleep_score_feedback VARCHAR(50),
    sleep_score_insight VARCHAR(50),
    sleep_score_personalized_insight VARCHAR(100),
    sleep_version INTEGER,
    breathing_disruption_severity VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Sleep scores table
CREATE TABLE IF NOT EXISTS sleep_scores (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    -- Overall score (embedded)
    overall_value INTEGER,
    overall_qualifier_key VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep score details table
CREATE TABLE IF NOT EXISTS sleep_score_details (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_scores_id VARCHAR NOT NULL,
    score_type VARCHAR(50) NOT NULL, -- totalDuration, stress, awakeCount, etc.
    value INTEGER,
    qualifier_key VARCHAR(50),
    -- Optimal Start (embedded ParsedValue)
    optimal_start_source VARCHAR(20),
    optimal_start_parsed_value INTEGER,
    -- Optimal End (embedded ParsedValue)
    optimal_end_source VARCHAR(20),
    optimal_end_parsed_value INTEGER,
    ideal_start_in_seconds DECIMAL(10,2),
    ideal_end_in_seconds DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_scores_id) REFERENCES sleep_scores(id) ON DELETE CASCADE
);

-- Sleep needs table
CREATE TABLE IF NOT EXISTS sleep_needs (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    user_profile_pk BIGINT NOT NULL,
    calendar_date DATE NOT NULL,
    device_id BIGINT NOT NULL,
    timestamp_gmt TIMESTAMP WITH TIME ZONE,
    baseline INTEGER NOT NULL,
    actual INTEGER NOT NULL,
    feedback VARCHAR(50),
    training_feedback VARCHAR(50),
    sleep_history_adjustment VARCHAR(50),
    hrv_adjustment VARCHAR(50),
    nap_adjustment VARCHAR(50),
    displayed_for_the_day BOOLEAN NOT NULL DEFAULT FALSE,
    preferred_activity_tracker BOOLEAN NOT NULL DEFAULT FALSE,
    need_type VARCHAR(20) NOT NULL, -- 'current' or 'next'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep movements table
CREATE TABLE IF NOT EXISTS sleep_movements (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    start_gmt TIMESTAMP WITH TIME ZONE NOT NULL,
    end_gmt TIMESTAMP WITH TIME ZONE NOT NULL,
    activity_level DECIMAL(15,12) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep levels table
CREATE TABLE IF NOT EXISTS sleep_levels (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    start_gmt TIMESTAMP WITH TIME ZONE NOT NULL,
    end_gmt TIMESTAMP WITH TIME ZONE NOT NULL,
    -- Activity Level (embedded ParsedValue)
    activity_level_source VARCHAR(20),
    activity_level_parsed_value INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep restless moments table
CREATE TABLE IF NOT EXISTS sleep_restless_moments (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    start_gmt TIMESTAMP WITH TIME ZONE NOT NULL,
    end_gmt TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep SpO2 summary table
CREATE TABLE IF NOT EXISTS sleep_spo2_summary (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    user_profile_pk BIGINT NOT NULL,
    device_id BIGINT NOT NULL,
    sleep_measurement_start_gmt TIMESTAMP WITH TIME ZONE,
    sleep_measurement_end_gmt TIMESTAMP WITH TIME ZONE,
    alert_threshold_value INTEGER,
    number_of_events_below_threshold INTEGER,
    duration_of_events_below_threshold INTEGER,
    -- Average SPO2 (embedded ParsedValue)
    average_spo2_source VARCHAR(20),
    average_spo2_parsed_value INTEGER,
    -- Average SpO2 HR (embedded ParsedValue)
    average_spo2_hr_source VARCHAR(20),
    average_spo2_hr_parsed_value INTEGER,
    lowest_spo2 INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep SpO2 data table
CREATE TABLE IF NOT EXISTS sleep_spo2_data (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    user_profile_pk BIGINT NOT NULL,
    epoch_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    device_id BIGINT NOT NULL,
    calendar_date DATE NOT NULL,
    epoch_duration INTEGER NOT NULL,
    spo2_reading INTEGER NOT NULL,
    reading_confidence INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep respiration data table
CREATE TABLE IF NOT EXISTS sleep_respiration_data (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    user_profile_pk BIGINT NOT NULL,
    epoch_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    device_id BIGINT NOT NULL,
    calendar_date DATE NOT NULL,
    epoch_duration INTEGER NOT NULL,
    respiration_reading INTEGER NOT NULL,
    reading_confidence INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep heart rate table
CREATE TABLE IF NOT EXISTS sleep_heart_rate (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    user_profile_pk BIGINT NOT NULL,
    epoch_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    device_id BIGINT NOT NULL,
    calendar_date DATE NOT NULL,
    epoch_duration INTEGER NOT NULL,
    heart_rate_reading INTEGER NOT NULL,
    reading_confidence INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep stress table
CREATE TABLE IF NOT EXISTS sleep_stress (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    user_profile_pk BIGINT NOT NULL,
    epoch_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    device_id BIGINT NOT NULL,
    calendar_date DATE NOT NULL,
    epoch_duration INTEGER NOT NULL,
    stress_reading INTEGER NOT NULL,
    reading_confidence INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep body battery table
CREATE TABLE IF NOT EXISTS sleep_body_battery (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    user_profile_pk BIGINT NOT NULL,
    epoch_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    device_id BIGINT NOT NULL,
    calendar_date DATE NOT NULL,
    epoch_duration INTEGER NOT NULL,
    body_battery_reading INTEGER NOT NULL,
    reading_confidence INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep HRV data table
CREATE TABLE IF NOT EXISTS sleep_hrv_data (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    user_profile_pk BIGINT NOT NULL,
    epoch_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    device_id BIGINT NOT NULL,
    calendar_date DATE NOT NULL,
    epoch_duration INTEGER NOT NULL,
    hrv_reading INTEGER NOT NULL,
    reading_confidence INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Sleep breathing disruption table
CREATE TABLE IF NOT EXISTS sleep_breathing_disruption (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    sleep_id BIGINT NOT NULL,
    user_profile_pk BIGINT NOT NULL,
    epoch_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    device_id BIGINT NOT NULL,
    calendar_date DATE NOT NULL,
    epoch_duration INTEGER NOT NULL,
    breathing_disruption_reading INTEGER NOT NULL,
    reading_confidence INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sleep_id) REFERENCES daily_sleep(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_daily_sleep_user_profile ON daily_sleep(user_profile_pk);
CREATE INDEX IF NOT EXISTS idx_daily_sleep_calendar_date ON daily_sleep(calendar_date);
CREATE INDEX IF NOT EXISTS idx_daily_sleep_user_date ON daily_sleep(user_profile_pk, calendar_date);

CREATE INDEX IF NOT EXISTS idx_sleep_scores_sleep_id ON sleep_scores(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_score_details_scores_id ON sleep_score_details(sleep_scores_id);
CREATE INDEX IF NOT EXISTS idx_sleep_needs_sleep_id ON sleep_needs(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_needs_user_date ON sleep_needs(user_profile_pk, calendar_date);

CREATE INDEX IF NOT EXISTS idx_sleep_movements_sleep_id ON sleep_movements(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_movements_time ON sleep_movements(start_gmt, end_gmt);

CREATE INDEX IF NOT EXISTS idx_sleep_levels_sleep_id ON sleep_levels(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_levels_time ON sleep_levels(start_gmt, end_gmt);

CREATE INDEX IF NOT EXISTS idx_sleep_restless_sleep_id ON sleep_restless_moments(sleep_id);

CREATE INDEX IF NOT EXISTS idx_sleep_spo2_summary_sleep_id ON sleep_spo2_summary(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_spo2_data_sleep_id ON sleep_spo2_data(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_spo2_data_timestamp ON sleep_spo2_data(epoch_timestamp);

CREATE INDEX IF NOT EXISTS idx_sleep_respiration_sleep_id ON sleep_respiration_data(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_respiration_timestamp ON sleep_respiration_data(epoch_timestamp);

CREATE INDEX IF NOT EXISTS idx_sleep_heart_rate_sleep_id ON sleep_heart_rate(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_heart_rate_timestamp ON sleep_heart_rate(epoch_timestamp);

CREATE INDEX IF NOT EXISTS idx_sleep_stress_sleep_id ON sleep_stress(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_stress_timestamp ON sleep_stress(epoch_timestamp);

CREATE INDEX IF NOT EXISTS idx_sleep_body_battery_sleep_id ON sleep_body_battery(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_body_battery_timestamp ON sleep_body_battery(epoch_timestamp);

CREATE INDEX IF NOT EXISTS idx_sleep_hrv_sleep_id ON sleep_hrv_data(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_hrv_timestamp ON sleep_hrv_data(epoch_timestamp);

CREATE INDEX IF NOT EXISTS idx_sleep_breathing_sleep_id ON sleep_breathing_disruption(sleep_id);
CREATE INDEX IF NOT EXISTS idx_sleep_breathing_timestamp ON sleep_breathing_disruption(epoch_timestamp);

-- Add comments to tables for documentation
COMMENT ON TABLE daily_sleep IS 'Main sleep data from Garmin devices';
COMMENT ON TABLE sleep_scores IS 'Sleep quality scoring metrics';
COMMENT ON TABLE sleep_score_details IS 'Detailed breakdown of sleep scores by type';
COMMENT ON TABLE sleep_needs IS 'Sleep need recommendations and feedback';
COMMENT ON TABLE sleep_movements IS 'Movement data during sleep periods';
COMMENT ON TABLE sleep_levels IS 'Sleep level classifications over time';
COMMENT ON TABLE sleep_restless_moments IS 'Periods of restless sleep';
COMMENT ON TABLE sleep_spo2_summary IS 'Summary of SpO2 measurements during sleep';
COMMENT ON TABLE sleep_spo2_data IS 'Individual SpO2 readings during sleep';
COMMENT ON TABLE sleep_respiration_data IS 'Respiration measurements during sleep';
COMMENT ON TABLE sleep_heart_rate IS 'Heart rate measurements during sleep';
COMMENT ON TABLE sleep_stress IS 'Stress level measurements during sleep';
COMMENT ON TABLE sleep_body_battery IS 'Body battery measurements during sleep';
COMMENT ON TABLE sleep_hrv_data IS 'Heart rate variability data during sleep';
COMMENT ON TABLE sleep_breathing_disruption IS 'Breathing disruption events during sleep';

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_daily_sleep_updated_at BEFORE UPDATE ON daily_sleep FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sleep_scores_updated_at BEFORE UPDATE ON sleep_scores FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sleep_score_details_updated_at BEFORE UPDATE ON sleep_score_details FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sleep_needs_updated_at BEFORE UPDATE ON sleep_needs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();