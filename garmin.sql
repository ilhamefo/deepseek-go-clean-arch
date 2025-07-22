-- Tabel utama untuk aktivitas
CREATE TABLE activities (
    activity_id BIGINT PRIMARY KEY,
    activity_name VARCHAR(255),
    start_time_local TIMESTAMP,
    start_time_gmt TIMESTAMP,
    end_time_gmt TIMESTAMP,
    distance DECIMAL(10, 5),
    duration DECIMAL(10, 5),
    elapsed_duration DECIMAL(10, 5),
    moving_duration DECIMAL(10, 5),
    elevation_gain DECIMAL(8, 2),
    elevation_loss DECIMAL(8, 2),
    average_speed DECIMAL(8, 5),
    max_speed DECIMAL(8, 5),
    start_latitude DECIMAL(15, 12),
    start_longitude DECIMAL(15, 12),
    end_latitude DECIMAL(15, 12),
    end_longitude DECIMAL(15, 12),
    has_polyline BOOLEAN,
    has_images BOOLEAN,
    owner_id BIGINT,
    owner_display_name VARCHAR(255),
    owner_full_name VARCHAR(100),
    owner_profile_image_url_small VARCHAR(500),
    owner_profile_image_url_medium VARCHAR(500),
    owner_profile_image_url_large VARCHAR(500),
    calories DECIMAL(8, 2),
    bmr_calories DECIMAL(8, 2),
    average_hr DECIMAL(5, 2),
    max_hr DECIMAL(5, 2),
    average_running_cadence DECIMAL(8, 4),
    max_running_cadence DECIMAL(8, 2),
    steps INTEGER,
    user_pro BOOLEAN,
    has_video BOOLEAN,
    time_zone_id INTEGER,
    begin_timestamp BIGINT,
    sport_type_id INTEGER,
    avg_power DECIMAL(8, 2),
    max_power DECIMAL(8, 2),
    aerobic_training_effect DECIMAL(8, 6),
    anaerobic_training_effect DECIMAL(8, 2),
    norm_power DECIMAL(8, 2),
    avg_vertical_oscillation DECIMAL(8, 6),
    avg_ground_contact_time DECIMAL(8, 6),
    avg_stride_length DECIMAL(8, 6),
    vo2_max_value DECIMAL(5, 2),
    avg_vertical_ratio DECIMAL(8, 6),
    device_id BIGINT,
    min_temperature DECIMAL(5, 2),
    max_temperature DECIMAL(5, 2),
    min_elevation DECIMAL(8, 6),
    max_elevation DECIMAL(8, 2),
    max_double_cadence DECIMAL(8, 2),
    max_vertical_speed FLOAT,
    manufacturer VARCHAR(50),
    location_name VARCHAR(100),
    lap_count INTEGER,
    water_estimated DECIMAL(8, 2),
    training_effect_label VARCHAR(50),
    min_activity_lap_duration DECIMAL(8, 6),
    aerobic_training_effect_message VARCHAR(100),
    anaerobic_training_effect_message VARCHAR(100),
    has_splits BOOLEAN,
    moderate_intensity_minutes INTEGER,
    vigorous_intensity_minutes INTEGER,
    avg_grade_adjusted_speed DECIMAL(8, 6),
    difference_body_battery INTEGER,
    has_heat_map BOOLEAN,
    fastest_split_1000 DECIMAL(8, 6),
    fastest_split_1609 DECIMAL(8, 6),
    hr_time_in_zone_1 DECIMAL(8, 3),
    hr_time_in_zone_2 DECIMAL(8, 3),
    hr_time_in_zone_3 DECIMAL(8, 3),
    hr_time_in_zone_4 DECIMAL(8, 3),
    hr_time_in_zone_5 DECIMAL(8, 3),
    power_time_in_zone_1 DECIMAL(8, 3),
    power_time_in_zone_2 DECIMAL(8, 3),
    power_time_in_zone_3 DECIMAL(8, 3),
    power_time_in_zone_4 DECIMAL(8, 3),
    power_time_in_zone_5 DECIMAL(8, 3),
    qualifying_dive BOOLEAN,
    parent BOOLEAN,
    pr BOOLEAN,
    favorite BOOLEAN,
    purposeful BOOLEAN,
    deco_dive BOOLEAN,
    manual_activity BOOLEAN,
    auto_calc_calories BOOLEAN,
    elevation_corrected BOOLEAN,
    atp_activity BOOLEAN
);

-- Tabel untuk tipe aktivitas
CREATE TABLE activity_types (
    type_id INTEGER PRIMARY KEY,
    type_key VARCHAR(50),
    parent_type_id INTEGER,
    is_hidden BOOLEAN,
    restricted BOOLEAN,
    trimmable BOOLEAN
);

-- Tabel untuk tipe event
CREATE TABLE event_types (
    type_id INTEGER PRIMARY KEY,
    type_key VARCHAR(50),
    sort_order INTEGER
);

-- Tabel untuk privacy settings
CREATE TABLE privacy_settings (
    type_id INTEGER PRIMARY KEY,
    type_key VARCHAR(50)
);

-- Tabel untuk user roles
CREATE TABLE user_roles (
    id INTEGER PRIMARY KEY,
    activity_id BIGINT,
    role_name VARCHAR(100),
    FOREIGN KEY (activity_id) REFERENCES activities(activity_id)
);

-- Tabel untuk split summaries
CREATE TABLE split_summaries (
    id INTEGER PRIMARY KEY,
    activity_id BIGINT,
    no_of_splits INTEGER,
    total_ascent DECIMAL(8, 2),
    duration DECIMAL(8, 2),
    split_type VARCHAR(50),
    num_climb_sends INTEGER,
    max_elevation_gain DECIMAL(8, 2),
    average_elevation_gain DECIMAL(8, 2),
    max_distance INTEGER,
    distance DECIMAL(10, 6),
    average_speed DECIMAL(8, 6),
    max_speed DECIMAL(8, 6),
    num_falls INTEGER,
    elevation_loss DECIMAL(8, 2),
    FOREIGN KEY (activity_id) REFERENCES activities(activity_id)
);

-- Menambahkan foreign key references ke tabel utama
ALTER TABLE activities 
ADD COLUMN activity_type_id INTEGER,
ADD COLUMN event_type_id INTEGER,
ADD COLUMN privacy_type_id INTEGER,
ADD FOREIGN KEY (activity_type_id) REFERENCES activity_types(type_id),
ADD FOREIGN KEY (event_type_id) REFERENCES event_types(type_id),
ADD FOREIGN KEY (privacy_type_id) REFERENCES privacy_settings(type_id);