-- ============================================================
-- Garmin Activity Details - PostgreSQL Schema Design
-- Source: activity-details.json
-- Extends existing garmin.sql schema
-- ============================================================

-- NOTE: Table 'activities' already exists in garmin.sql
-- This schema adds detailed metrics tables to complement it

-- ============================================================
-- TABLE 1: metric_units
-- Master data for measurement units
-- ============================================================
CREATE TABLE IF NOT EXISTS metric_units (
    unit_id INT PRIMARY KEY,
    unit_key VARCHAR(50) NOT NULL UNIQUE,
    factor NUMERIC(10, 2) NOT NULL
);

INSERT INTO metric_units (unit_id, unit_key, factor) VALUES
    (1,   'meter',          100.0),
    (5,   'centimeter',     1.0),
    (6,   'dimensionless',  1.0),
    (10,  'watt',           1.0),
    (20,  'mps',            0.1),
    (40,  'second',         1000.0),
    (41,  'ms',             1.0),
    (60,  'dd',             1.0),
    (92,  'stepsPerMinute', 1.0),
    (100, 'bpm',            1.0),
    (120, 'gmt',            0.0),
    (227, 'celcius',        1.0)
ON CONFLICT (unit_id) DO NOTHING;

COMMENT ON TABLE metric_units IS 'Master data for measurement units (from metricDescriptors.unit)';

-- ============================================================
-- TABLE 2: metric_descriptors
-- Defines available metrics for each activity
-- Maps metric index to metric name and unit
-- ============================================================
CREATE TABLE IF NOT EXISTS metric_descriptors (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES activities(activity_id) ON DELETE CASCADE,
    metrics_index SMALLINT NOT NULL,
    metric_key VARCHAR(100) NOT NULL,
    unit_id INT NOT NULL REFERENCES metric_units(unit_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_activity_metric_index UNIQUE (activity_id, metrics_index)
);

CREATE INDEX idx_md_activity ON metric_descriptors(activity_id);
CREATE INDEX idx_md_key ON metric_descriptors(metric_key);

COMMENT ON TABLE metric_descriptors IS 'Schema definition from metricDescriptors[] in activity-details.json';
COMMENT ON COLUMN metric_descriptors.metrics_index IS 'Position in metrics array (0-22)';
COMMENT ON COLUMN metric_descriptors.metric_key IS 'Metric key from JSON (e.g., directHeartRate, sumDuration)';

-- ============================================================
-- TABLE 3: activity_metrics_timeseries
-- Expanded time-series data points (from activityDetailMetrics[])
-- Alternative to storing metrics as JSONB in activity_detail_metrics
-- ============================================================
CREATE TABLE IF NOT EXISTS activity_metrics_timeseries (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES activities(activity_id) ON DELETE CASCADE,
    sequence INT NOT NULL,
    
    -- Metric index 0: sumDuration (second) - cumulative
    sum_duration NUMERIC(12, 2),
    
    -- Metric index 1: directPower (watt)
    direct_power NUMERIC(10, 2),
    
    -- Metric index 2: directGradeAdjustedSpeed (mps)
    direct_grade_adjusted_speed NUMERIC(8, 4),
    
    -- Metric index 3: directAirTemperature (celsius)
    direct_air_temperature NUMERIC(5, 1),
    
    -- Metric index 4: directHeartRate (bpm)
    direct_heart_rate SMALLINT,
    
    -- Metric index 5: sumAccumulatedPower (watt) - cumulative
    sum_accumulated_power NUMERIC(12, 2),
    
    -- Metric index 6: directFractionalCadence (stepsPerMinute)
    direct_fractional_cadence NUMERIC(6, 2),
    
    -- Metric index 7: directBodyBattery (0-100)
    direct_body_battery SMALLINT CHECK (direct_body_battery IS NULL OR direct_body_battery BETWEEN 0 AND 100),
    
    -- Metric index 8: directElevation (meter)
    direct_elevation NUMERIC(8, 2),
    
    -- Metric index 9: directRunCadence (stepsPerMinute)
    direct_run_cadence SMALLINT,
    
    -- Metric index 10: directDoubleCadence (stepsPerMinute)
    direct_double_cadence SMALLINT,
    
    -- Metric index 11: directSpeed (mps)
    direct_speed NUMERIC(8, 4),
    
    -- Metric index 12: sumMovingDuration (second) - cumulative
    sum_moving_duration NUMERIC(12, 2),
    
    -- Metric index 13: sumDistance (meter) - cumulative
    sum_distance NUMERIC(12, 2),
    
    -- Metric index 14: sumElapsedDuration (second) - cumulative
    sum_elapsed_duration NUMERIC(12, 2),
    
    -- Metric index 15: directTimestamp (unix epoch milliseconds)
    direct_timestamp BIGINT,
    
    -- Metric index 16: directLongitude (decimal degrees)
    direct_longitude DOUBLE PRECISION CHECK (direct_longitude IS NULL OR direct_longitude BETWEEN -180 AND 180),
    
    -- Metric index 17: directVerticalOscillation (centimeter)
    direct_vertical_oscillation NUMERIC(6, 2),
    
    -- Metric index 18: directLatitude (decimal degrees)
    direct_latitude DOUBLE PRECISION CHECK (direct_latitude IS NULL OR direct_latitude BETWEEN -90 AND 90),
    
    -- Metric index 19: directVerticalRatio (dimensionless)
    direct_vertical_ratio NUMERIC(6, 2),
    
    -- Metric index 20: directStrideLength (centimeter)
    direct_stride_length NUMERIC(6, 2),
    
    -- Metric index 21: directVerticalSpeed (mps)
    direct_vertical_speed NUMERIC(8, 4),
    
    -- Metric index 22: directGroundContactTime (ms)
    direct_ground_contact_time NUMERIC(6, 2),
    
    -- Computed columns
    recorded_at TIMESTAMPTZ GENERATED ALWAYS AS (
        CASE WHEN direct_timestamp IS NOT NULL 
        THEN to_timestamp(direct_timestamp / 1000.0)
        ELSE NULL END
    ) STORED,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_activity_sequence UNIQUE (activity_id, sequence)
);

-- Performance indexes
CREATE INDEX idx_amts_activity ON activity_metrics_timeseries(activity_id);
CREATE INDEX idx_amts_timestamp ON activity_metrics_timeseries(direct_timestamp) 
    WHERE direct_timestamp IS NOT NULL;
CREATE INDEX idx_amts_recorded_at ON activity_metrics_timeseries(recorded_at) 
    WHERE recorded_at IS NOT NULL;
CREATE INDEX idx_amts_location ON activity_metrics_timeseries(direct_latitude, direct_longitude) 
    WHERE direct_latitude IS NOT NULL AND direct_longitude IS NOT NULL;
CREATE INDEX idx_amts_sequence ON activity_metrics_timeseries(activity_id, sequence);
CREATE INDEX idx_amts_heart_rate ON activity_metrics_timeseries(activity_id, direct_heart_rate)
    WHERE direct_heart_rate IS NOT NULL;

COMMENT ON TABLE activity_metrics_timeseries IS 'Expanded time-series metrics from activityDetailMetrics[] (~4000+ rows per activity)';
COMMENT ON COLUMN activity_metrics_timeseries.sequence IS 'Measurement sequence number (0-based index from JSON array)';
COMMENT ON COLUMN activity_metrics_timeseries.recorded_at IS 'Human-readable timestamp converted from direct_timestamp';

-- ============================================================
-- TABLE 4: geo_polylines
-- GPS route data for activities
-- ============================================================
CREATE TABLE IF NOT EXISTS geo_polylines (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL UNIQUE REFERENCES activities(activity_id) ON DELETE CASCADE,
    start_point JSONB,
    end_point JSONB,
    min_latitude DOUBLE PRECISION CHECK (min_latitude IS NULL OR min_latitude BETWEEN -90 AND 90),
    max_latitude DOUBLE PRECISION CHECK (max_latitude IS NULL OR max_latitude BETWEEN -90 AND 90),
    min_longitude DOUBLE PRECISION CHECK (min_longitude IS NULL OR min_longitude BETWEEN -180 AND 180),
    max_longitude DOUBLE PRECISION CHECK (max_longitude IS NULL OR max_longitude BETWEEN -180 AND 180),
    polyline JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gp_bbox ON geo_polylines(min_latitude, max_latitude, min_longitude, max_longitude)
    WHERE min_latitude IS NOT NULL;

COMMENT ON TABLE geo_polylines IS 'GPS polyline from geoPolylineDTO in activity-details.json';
COMMENT ON COLUMN geo_polylines.polyline IS 'Array of GPS coordinate points';

-- ============================================================
-- TABLE 5: heart_rate_timeseries
-- Dedicated heart rate data (from heartRateDTOs if present)
-- ============================================================
CREATE TABLE IF NOT EXISTS heart_rate_timeseries (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES activities(activity_id) ON DELETE CASCADE,
    timestamp_ms BIGINT NOT NULL,
    heart_rate SMALLINT NOT NULL CHECK (heart_rate > 0 AND heart_rate < 300),
    recorded_at TIMESTAMPTZ GENERATED ALWAYS AS (to_timestamp(timestamp_ms / 1000.0)) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_hrts_activity ON heart_rate_timeseries(activity_id, timestamp_ms);
CREATE INDEX idx_hrts_recorded_at ON heart_rate_timeseries(recorded_at);

COMMENT ON TABLE heart_rate_timeseries IS 'Dedicated HR data from heartRateDTOs[] if provided separately';

-- ============================================================
-- TABLE 6: activity_details_summary
-- Stores top-level metadata from activity-details.json
-- ============================================================
CREATE TABLE IF NOT EXISTS activity_details_summary (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL UNIQUE REFERENCES activities(activity_id) ON DELETE CASCADE,
    measurement_count INT NOT NULL,
    metrics_count INT NOT NULL,
    total_metrics_count INT NOT NULL,
    details_available BOOLEAN NOT NULL DEFAULT TRUE,
    pending_data JSONB,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ads_fetched ON activity_details_summary(fetched_at DESC);

COMMENT ON TABLE activity_details_summary IS 'Top-level metadata from activity-details.json response';
COMMENT ON COLUMN activity_details_summary.measurement_count IS 'Number of metric types (e.g., 23)';
COMMENT ON COLUMN activity_details_summary.metrics_count IS 'Number of data points (e.g., 4005)';

-- ============================================================
-- VIEW: v_activity_metrics_summary
-- Aggregated statistics combining activities and metrics_timeseries
-- ============================================================
CREATE OR REPLACE VIEW v_activity_metrics_summary AS
SELECT 
    a.activity_id,
    a.activity_name,
    a.start_time_local,
    a.sport_type_id,
    ads.measurement_count,
    ads.metrics_count,
    COUNT(amt.id) as actual_data_points,
    MIN(amt.recorded_at) as ts_start_time,
    MAX(amt.recorded_at) as ts_end_time,
    EXTRACT(EPOCH FROM (MAX(amt.recorded_at) - MIN(amt.recorded_at))) as ts_duration_seconds,
    
    -- Distance & Duration (from timeseries)
    MAX(amt.sum_distance) as ts_total_distance_m,
    MAX(amt.sum_duration) as ts_total_duration_s,
    MAX(amt.sum_moving_duration) as ts_moving_duration_s,
    
    -- Heart Rate
    ROUND(AVG(amt.direct_heart_rate), 0) as ts_avg_heart_rate,
    MAX(amt.direct_heart_rate) as ts_max_heart_rate,
    MIN(amt.direct_heart_rate) as ts_min_heart_rate,
    
    -- Speed
    ROUND(AVG(amt.direct_speed), 2) as ts_avg_speed_mps,
    MAX(amt.direct_speed) as ts_max_speed_mps,
    ROUND(AVG(amt.direct_grade_adjusted_speed), 2) as ts_avg_grade_adj_speed,
    
    -- Elevation
    ROUND(AVG(amt.direct_elevation), 1) as ts_avg_elevation_m,
    MAX(amt.direct_elevation) as ts_max_elevation_m,
    MIN(amt.direct_elevation) as ts_min_elevation_m,
    
    -- Cadence
    ROUND(AVG(amt.direct_run_cadence), 1) as ts_avg_cadence,
    MAX(amt.direct_run_cadence) as ts_max_cadence,
    
    -- Power
    ROUND(AVG(amt.direct_power), 1) as ts_avg_power,
    MAX(amt.direct_power) as ts_max_power,
    
    -- Running Dynamics
    ROUND(AVG(amt.direct_vertical_oscillation), 2) as ts_avg_vertical_osc_cm,
    ROUND(AVG(amt.direct_ground_contact_time), 2) as ts_avg_gct_ms,
    ROUND(AVG(amt.direct_stride_length), 2) as ts_avg_stride_length_cm
    
FROM activities a
LEFT JOIN activity_details_summary ads ON a.activity_id = ads.activity_id
LEFT JOIN activity_metrics_timeseries amt ON a.activity_id = amt.activity_id
GROUP BY a.activity_id, a.activity_name, a.start_time_local, a.sport_type_id, 
         ads.measurement_count, ads.metrics_count;

COMMENT ON VIEW v_activity_metrics_summary IS 'Combined summary from activities table and timeseries metrics';

-- ============================================================
-- FUNCTION: Auto-update updated_at timestamp
-- ============================================================
CREATE OR REPLACE FUNCTION fn_update_details_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_activity_details_updated
    BEFORE UPDATE ON activity_details_summary
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_details_timestamp();

-- ============================================================
-- Sample Queries
-- ============================================================

-- Query 1: Get activity with detailed metrics summary
-- SELECT * FROM v_activity_metrics_summary WHERE activity_id = 22019284616;

-- Query 2: GPS tracking points
-- SELECT 
--     sequence,
--     direct_latitude,
--     direct_longitude,
--     direct_elevation,
--     direct_speed,
--     recorded_at
-- FROM activity_metrics_timeseries 
-- WHERE activity_id = 22019284616 
--   AND direct_latitude IS NOT NULL
-- ORDER BY sequence;

-- Query 3: Heart rate zones analysis
-- SELECT 
--     CASE 
--         WHEN direct_heart_rate < 100 THEN 'Zone 1: Rest'
--         WHEN direct_heart_rate < 130 THEN 'Zone 2: Easy'
--         WHEN direct_heart_rate < 150 THEN 'Zone 3: Moderate'
--         WHEN direct_heart_rate < 170 THEN 'Zone 4: Hard'
--         ELSE 'Zone 5: Maximum'
--     END as hr_zone,
--     COUNT(*) as sample_count,
--     ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 1) as percentage,
--     ROUND(AVG(direct_heart_rate), 0) as avg_hr,
--     MIN(direct_heart_rate) as min_hr,
--     MAX(direct_heart_rate) as max_hr
-- FROM activity_metrics_timeseries 
-- WHERE activity_id = 22019284616 
--   AND direct_heart_rate IS NOT NULL
-- GROUP BY hr_zone
-- ORDER BY avg_hr;

-- Query 4: Pace analysis (speed over distance)
-- SELECT 
--     FLOOR(sequence / 100) as segment,
--     ROUND(AVG(direct_speed), 2) as avg_speed_mps,
--     ROUND(AVG(1000.0 / NULLIF(direct_speed, 0) / 60), 2) as avg_pace_min_per_km,
--     ROUND(AVG(direct_heart_rate), 0) as avg_hr,
--     ROUND(AVG(direct_elevation), 1) as avg_elevation
-- FROM activity_metrics_timeseries
-- WHERE activity_id = 22019284616 
--   AND direct_speed > 0
-- GROUP BY segment
-- ORDER BY segment;

-- Query 5: Elevation profile with gains/losses
-- SELECT 
--     sequence,
--     direct_elevation,
--     direct_elevation - LAG(direct_elevation) OVER (ORDER BY sequence) as elevation_change,
--     CASE 
--         WHEN direct_elevation - LAG(direct_elevation) OVER (ORDER BY sequence) > 1 THEN 'Climbing'
--         WHEN direct_elevation - LAG(direct_elevation) OVER (ORDER BY sequence) < -1 THEN 'Descending'
--         ELSE 'Flat'
--     END as terrain_type,
--     recorded_at
-- FROM activity_metrics_timeseries
-- WHERE activity_id = 22019284616 
--   AND direct_elevation IS NOT NULL
-- ORDER BY sequence;

-- Query 6: Compare activities table with timeseries data
-- SELECT 
--     a.activity_id,
--     a.activity_name,
--     a.distance as summary_distance_m,
--     MAX(amt.sum_distance) as ts_distance_m,
--     a.duration as summary_duration_s,
--     MAX(amt.sum_duration) as ts_duration_s,
--     a.average_hr as summary_avg_hr,
--     ROUND(AVG(amt.direct_heart_rate), 0) as ts_avg_hr
-- FROM activities a
-- LEFT JOIN activity_metrics_timeseries amt ON a.activity_id = amt.activity_id
-- WHERE a.activity_id = 22019284616
-- GROUP BY a.activity_id, a.activity_name, a.distance, a.duration, a.average_hr;

-- ============================================================
-- TABLE 2: metric_units
-- Master data for measurement units
-- ============================================================
CREATE TABLE IF NOT EXISTS metric_units (
    unit_id INT PRIMARY KEY,
    unit_key VARCHAR(50) NOT NULL UNIQUE,
    factor NUMERIC(10, 2) NOT NULL
);

INSERT INTO metric_units (unit_id, unit_key, factor) VALUES
    (1,   'meter',          100.0),
    (5,   'centimeter',     1.0),
    (6,   'dimensionless',  1.0),
    (10,  'watt',           1.0),
    (20,  'mps',            0.1),
    (40,  'second',         1000.0),
    (41,  'ms',             1.0),
    (60,  'dd',             1.0),
    (92,  'stepsPerMinute', 1.0),
    (100, 'bpm',            1.0),
    (120, 'gmt',            0.0),
    (227, 'celcius',        1.0)
ON CONFLICT (unit_id) DO NOTHING;

COMMENT ON TABLE metric_units IS 'Master data for measurement units';

-- ============================================================
-- TABLE 3: metric_descriptors
-- Defines available metrics for each activity
-- Maps metric index to metric name and unit
-- ============================================================
CREATE TABLE IF NOT EXISTS metric_descriptors (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES activities(activity_id) ON DELETE CASCADE,
    metrics_index SMALLINT NOT NULL,
    metric_key VARCHAR(100) NOT NULL,
    unit_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_activity_metric_index UNIQUE (activity_id, metrics_index)
);

CREATE INDEX idx_md_activity ON metric_descriptors(activity_id);
CREATE INDEX idx_md_key ON metric_descriptors(metric_key);

COMMENT ON TABLE metric_descriptors IS 'Schema definition: maps metric index to name and unit for each activity';
COMMENT ON COLUMN metric_descriptors.metrics_index IS 'Position in metrics array (0-22)';
COMMENT ON COLUMN metric_descriptors.metric_key IS 'Metric name (e.g., directHeartRate, sumDuration)';

-- ============================================================
-- TABLE 4: activity_metrics
-- Time-series data points
-- Stores all 23 metrics for each measurement point
-- ============================================================
CREATE TABLE IF NOT EXISTS activity_metrics (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES activities(activity_id) ON DELETE CASCADE,
    sequence INT NOT NULL,
    
    -- Metric index 0: sumDuration (second)
    sum_duration NUMERIC(12, 2),
    
    -- Metric index 1: directPower (watt)
    direct_power NUMERIC(10, 2),
    
    -- Metric index 2: directGradeAdjustedSpeed (mps)
    direct_grade_adjusted_speed NUMERIC(8, 4),
    
    -- Metric index 3: directAirTemperature (celsius)
    direct_air_temperature NUMERIC(5, 1),
    
    -- Metric index 4: directHeartRate (bpm)
    direct_heart_rate SMALLINT,
    
    -- Metric index 5: sumAccumulatedPower (watt)
    sum_accumulated_power NUMERIC(12, 2),
    
    -- Metric index 6: directFractionalCadence (stepsPerMinute)
    direct_fractional_cadence NUMERIC(6, 2),
    
    -- Metric index 7: directBodyBattery (dimensionless, 0-100)
    direct_body_battery SMALLINT,
    
    -- Metric index 8: directElevation (meter)
    direct_elevation NUMERIC(8, 2),
    
    -- Metric index 9: directRunCadence (stepsPerMinute)
    direct_run_cadence SMALLINT,
    
    -- Metric index 10: directDoubleCadence (stepsPerMinute)
    direct_double_cadence SMALLINT,
    
    -- Metric index 11: directSpeed (mps)
    direct_speed NUMERIC(8, 4),
    
    -- Metric index 12: sumMovingDuration (second)
    sum_moving_duration NUMERIC(12, 2),
    
    -- Metric index 13: sumDistance (meter)
    sum_distance NUMERIC(12, 2),
    
    -- Metric index 14: sumElapsedDuration (second)
    sum_elapsed_duration NUMERIC(12, 2),
    
    -- Metric index 15: directTimestamp (gmt epoch milliseconds)
    direct_timestamp BIGINT,
    
    -- Metric index 16: directLongitude (decimal degrees)
    direct_longitude DOUBLE PRECISION,
    
    -- Metric index 17: directVerticalOscillation (centimeter)
    direct_vertical_oscillation NUMERIC(6, 2),
    
    -- Metric index 18: directLatitude (decimal degrees)
    direct_latitude DOUBLE PRECISION,
    
    -- Metric index 19: directVerticalRatio (dimensionless)
    direct_vertical_ratio NUMERIC(6, 2),
    
    -- Metric index 20: directStrideLength (centimeter)
    direct_stride_length NUMERIC(6, 2),
    
    -- Metric index 21: directVerticalSpeed (mps)
    direct_vertical_speed NUMERIC(8, 4),
    
    -- Metric index 22: directGroundContactTime (ms)
    direct_ground_contact_time NUMERIC(6, 2),
    
    -- Computed column
    recorded_at TIMESTAMPTZ GENERATED ALWAYS AS (
        to_timestamp(direct_timestamp / 1000.0)
    ) STORED,
    
    CONSTRAINT uq_activity_sequence UNIQUE (activity_id, sequence)
);

CREATE INDEX idx_am_activity ON activity_metrics(activity_id);
CREATE INDEX idx_am_timestamp ON activity_metrics(direct_timestamp) WHERE direct_timestamp IS NOT NULL;
CREATE INDEX idx_am_recorded_at ON activity_metrics(recorded_at) WHERE recorded_at IS NOT NULL;
CREATE INDEX idx_am_location ON activity_metrics(direct_latitude, direct_longitude) 
    WHERE direct_latitude IS NOT NULL AND direct_longitude IS NOT NULL;
CREATE INDEX idx_am_activity_seq ON activity_metrics(activity_id, sequence);

COMMENT ON TABLE activity_metrics IS 'Time-series metrics data (~4000 rows per activity)';
COMMENT ON COLUMN activity_metrics.sequence IS 'Measurement sequence number (0-based index)';
COMMENT ON COLUMN activity_metrics.recorded_at IS 'Timestamp converted from direct_timestamp milliseconds';

-- ============================================================
-- TABLE 5: geo_polylines
-- GPS route data for activities
-- ============================================================
CREATE TABLE IF NOT EXISTS geo_polylines (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL UNIQUE REFERENCES activities(activity_id) ON DELETE CASCADE,
    start_point JSONB,
    end_point JSONB,
    min_latitude DOUBLE PRECISION,
    max_latitude DOUBLE PRECISION,
    min_longitude DOUBLE PRECISION,
    max_longitude DOUBLE PRECISION,
    polyline JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gp_bbox ON geo_polylines(min_latitude, max_latitude, min_longitude, max_longitude)
    WHERE min_latitude IS NOT NULL;

COMMENT ON TABLE geo_polylines IS 'GPS polyline and bounding box for activity routes';
COMMENT ON COLUMN geo_polylines.polyline IS 'Array of GPS coordinate points';

-- ============================================================
-- TABLE 6: heart_rate_data
-- Separate heart rate measurements (if provided)
-- ============================================================
CREATE TABLE IF NOT EXISTS heart_rate_data (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES activities(activity_id) ON DELETE CASCADE,
    timestamp BIGINT NOT NULL,
    heart_rate SMALLINT NOT NULL CHECK (heart_rate > 0 AND heart_rate < 300),
    recorded_at TIMESTAMPTZ GENERATED ALWAYS AS (to_timestamp(timestamp / 1000.0)) STORED
);

CREATE INDEX idx_hr_activity ON heart_rate_data(activity_id, timestamp);
CREATE INDEX idx_hr_recorded_at ON heart_rate_data(recorded_at);

COMMENT ON TABLE heart_rate_data IS 'Dedicated heart rate data points (from heartRateDTOs if present)';

-- ============================================================
-- VIEW: activity_summary
-- Quick access to activity statistics
-- ============================================================
CREATE OR REPLACE VIEW v_activity_summary AS
SELECT 
    a.activity_id,
    a.created_at as activity_date,
    a.measurement_count,
    a.metrics_count,
    COUNT(am.id) as actual_data_points,
    MIN(am.recorded_at) as start_time,
    MAX(am.recorded_at) as end_time,
    EXTRACT(EPOCH FROM (MAX(am.recorded_at) - MIN(am.recorded_at))) as duration_seconds,
    MAX(am.sum_distance) as total_distance_meters,
    ROUND(AVG(am.direct_heart_rate), 0) as avg_heart_rate_bpm,
    MAX(am.direct_heart_rate) as max_heart_rate_bpm,
    MIN(am.direct_heart_rate) as min_heart_rate_bpm,
    ROUND(AVG(am.direct_speed), 2) as avg_speed_mps,
    MAX(am.direct_speed) as max_speed_mps,
    ROUND(AVG(am.direct_elevation), 1) as avg_elevation_m,
    MAX(am.direct_elevation) as max_elevation_m,
    MIN(am.direct_elevation) as min_elevation_m
FROM activities a
LEFT JOIN activity_metrics am ON a.activity_id = am.activity_id
GROUP BY a.activity_id, a.created_at, a.measurement_count, a.metrics_count;

COMMENT ON VIEW v_activity_summary IS 'Aggregated statistics for each activity';

-- ============================================================
-- FUNCTION: Auto-update updated_at timestamp
-- ============================================================
CREATE OR REPLACE FUNCTION fn_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_activities_updated
    BEFORE UPDATE ON activities
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_timestamp();

-- ============================================================
-- Sample queries
-- ============================================================

-- Query 1: Get activity summary
-- SELECT * FROM v_activity_summary WHERE activity_id = 22019284616;

-- Query 2: Get GPS route points
-- SELECT sequence, direct_latitude, direct_longitude, recorded_at
-- FROM activity_metrics 
-- WHERE activity_id = 22019284616 
--   AND direct_latitude IS NOT NULL
-- ORDER BY sequence;

-- Query 3: Heart rate analysis
-- SELECT 
--     CASE 
--         WHEN direct_heart_rate < 100 THEN 'Rest'
--         WHEN direct_heart_rate < 130 THEN 'Low'
--         WHEN direct_heart_rate < 150 THEN 'Medium'
--         WHEN direct_heart_rate < 170 THEN 'High'
--         ELSE 'Max'
--     END as hr_zone,
--     COUNT(*) as sample_count,
--     ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 1) as percentage,
--     ROUND(AVG(direct_heart_rate), 0) as avg_hr,
--     MIN(direct_heart_rate) as min_hr,
--     MAX(direct_heart_rate) as max_hr
-- FROM activity_metrics_timeseries 
-- WHERE activity_id = 22019284616 
--   AND direct_heart_rate IS NOT NULL
-- GROUP BY hr_zone
-- ORDER BY avg_hr;

-- Query 4: Pace analysis (speed over distance)
-- SELECT 
--     FLOOR(sequence / 100) as segment,
--     ROUND(AVG(direct_speed), 2) as avg_speed_mps,
--     ROUND(AVG(1000.0 / NULLIF(direct_speed, 0) / 60), 2) as avg_pace_min_per_km,
--     ROUND(AVG(direct_heart_rate), 0) as avg_hr,
--     ROUND(AVG(direct_elevation), 1) as avg_elevation
-- FROM activity_metrics_timeseries
-- WHERE activity_id = 22019284616 
--   AND direct_speed > 0
-- GROUP BY segment
-- ORDER BY segment;

-- Query 5: Elevation profile with gains/losses
-- SELECT 
--     sequence,
--     direct_elevation,
--     direct_elevation - LAG(direct_elevation) OVER (ORDER BY sequence) as elevation_change,
--     CASE 
--         WHEN direct_elevation - LAG(direct_elevation) OVER (ORDER BY sequence) > 1 THEN 'Climbing'
--         WHEN direct_elevation - LAG(direct_elevation) OVER (ORDER BY sequence) < -1 THEN 'Descending'
--         ELSE 'Flat'
--     END as terrain_type,
--     recorded_at
-- FROM activity_metrics_timeseries
-- WHERE activity_id = 22019284616 
--   AND direct_elevation IS NOT NULL
-- ORDER BY sequence;

-- Query 6: Compare activities table with timeseries data
-- SELECT 
--     a.activity_id,
--     a.activity_name,
--     a.distance as summary_distance_m,
--     MAX(amt.sum_distance) as ts_distance_m,
--     a.duration as summary_duration_s,
--     MAX(amt.sum_duration) as ts_duration_s,
--     a.average_hr as summary_avg_hr,
--     ROUND(AVG(amt.direct_heart_rate), 0) as ts_avg_hr
-- FROM activities a
-- LEFT JOIN activity_metrics_timeseries amt ON a.activity_id = amt.activity_id
-- WHERE a.activity_id = 22019284616
-- GROUP BY a.activity_id, a.activity_name, a.distance, a.duration, a.average_hr;
