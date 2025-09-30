-- ----------------------------
-- Sequence structure for activity_detail_metrics_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."activity_detail_metrics_id_seq";
CREATE SEQUENCE "public"."activity_detail_metrics_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Table structure for activities
-- ----------------------------
DROP TABLE IF EXISTS "public"."activities";
CREATE TABLE "public"."activities" (
  "activity_id" int8 NOT NULL,
  "activity_name" varchar(255) COLLATE "pg_catalog"."default",
  "start_time_local" timestamp(6),
  "start_time_gmt" timestamp(6),
  "end_time_gmt" timestamp(6),
  "distance" float8,
  "duration" float8,
  "elapsed_duration" float8,
  "moving_duration" float8,
  "elevation_gain" float8,
  "elevation_loss" float8,
  "average_speed" float8,
  "max_speed" float8,
  "start_latitude" float8,
  "start_longitude" float8,
  "end_latitude" float8,
  "end_longitude" float8,
  "has_polyline" bool,
  "has_images" bool,
  "owner_id" int8,
  "owner_display_name" varchar(255) COLLATE "pg_catalog"."default",
  "owner_full_name" varchar(100) COLLATE "pg_catalog"."default",
  "owner_profile_image_url_small" varchar(500) COLLATE "pg_catalog"."default",
  "owner_profile_image_url_medium" varchar(500) COLLATE "pg_catalog"."default",
  "owner_profile_image_url_large" varchar(500) COLLATE "pg_catalog"."default",
  "calories" float8,
  "bmr_calories" float8,
  "average_hr" float8,
  "max_hr" float8,
  "average_running_cadence" float8,
  "max_running_cadence" float8,
  "steps" int4,
  "user_pro" bool,
  "has_video" bool,
  "time_zone_id" int4,
  "begin_timestamp" int8,
  "sport_type_id" int4,
  "avg_power" float8,
  "max_power" float8,
  "aerobic_training_effect" float8,
  "anaerobic_training_effect" float8,
  "norm_power" float8,
  "avg_vertical_oscillation" float8,
  "avg_ground_contact_time" float8,
  "avg_stride_length" float8,
  "vo2_max_value" float8,
  "avg_vertical_ratio" float8,
  "device_id" int8,
  "min_temperature" float8,
  "max_temperature" float8,
  "min_elevation" float8,
  "max_elevation" float8,
  "max_double_cadence" float8,
  "max_vertical_speed" float8,
  "manufacturer" varchar(50) COLLATE "pg_catalog"."default",
  "location_name" varchar(100) COLLATE "pg_catalog"."default",
  "lap_count" int4,
  "water_estimated" float8,
  "training_effect_label" varchar(50) COLLATE "pg_catalog"."default",
  "min_activity_lap_duration" float8,
  "aerobic_training_effect_message" varchar(100) COLLATE "pg_catalog"."default",
  "anaerobic_training_effect_message" varchar(100) COLLATE "pg_catalog"."default",
  "has_splits" bool,
  "moderate_intensity_minutes" int4,
  "vigorous_intensity_minutes" int4,
  "avg_grade_adjusted_speed" float8,
  "difference_body_battery" int4,
  "has_heat_map" bool,
  "fastest_split_1000" float8,
  "fastest_split_1609" float8,
  "hr_time_in_zone_1" float8,
  "hr_time_in_zone_2" float8,
  "hr_time_in_zone_3" float8,
  "hr_time_in_zone_4" float8,
  "hr_time_in_zone_5" float8,
  "power_time_in_zone_1" float8,
  "power_time_in_zone_2" float8,
  "power_time_in_zone_3" float8,
  "power_time_in_zone_4" float8,
  "power_time_in_zone_5" float8,
  "qualifying_dive" bool,
  "parent" bool,
  "pr" bool,
  "favorite" bool,
  "purposeful" bool,
  "deco_dive" bool,
  "manual_activity" bool,
  "auto_calc_calories" bool,
  "elevation_corrected" bool,
  "atp_activity" bool,
  "activity_type_id" int4,
  "event_type_id" int4,
  "privacy_type_id" int4
)
;

-- ----------------------------
-- Table structure for activity_detail_metrics
-- ----------------------------
DROP TABLE IF EXISTS "public"."activity_detail_metrics";
CREATE TABLE "public"."activity_detail_metrics" (
  "id" int4 NOT NULL DEFAULT nextval('activity_detail_metrics_id_seq'::regclass),
  "activity_id" int8 NOT NULL,
  "measurement_index" int4 NOT NULL,
  "metrics" jsonb,
  "created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for activity_splits
-- ----------------------------
DROP TABLE IF EXISTS "public"."activity_splits";
CREATE TABLE "public"."activity_splits" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "activity_id" int8 NOT NULL,
  "start_time_gmt" varchar(50) COLLATE "pg_catalog"."default",
  "start_latitude" float8,
  "start_longitude" float8,
  "end_latitude" float8,
  "end_longitude" float8,
  "distance" float8,
  "duration" float8,
  "moving_duration" float8,
  "elapsed_duration" float8,
  "elevation_gain" float8,
  "elevation_loss" float8,
  "max_elevation" float8,
  "min_elevation" float8,
  "average_speed" float8,
  "average_moving_speed" float8,
  "max_speed" float8,
  "calories" float8,
  "bmr_calories" float8,
  "average_hr" float8,
  "max_hr" float8,
  "average_run_cadence" float8,
  "max_run_cadence" float8,
  "average_temperature" float8,
  "max_temperature" float8,
  "min_temperature" float8,
  "average_power" float8,
  "max_power" float8,
  "min_power" float8,
  "normalized_power" float8,
  "total_work" float8,
  "ground_contact_time" float8,
  "stride_length" float8,
  "vertical_oscillation" float8,
  "vertical_ratio" float8,
  "max_vertical_speed" float8,
  "avg_grade_adjusted_speed" float8,
  "lap_index" int4,
  "intensity_type" varchar(20) COLLATE "pg_catalog"."default",
  "message_index" int4
)
;

-- ----------------------------
-- Table structure for activity_types
-- ----------------------------
DROP TABLE IF EXISTS "public"."activity_types";
CREATE TABLE "public"."activity_types" (
  "type_id" int4 NOT NULL,
  "type_key" varchar(50) COLLATE "pg_catalog"."default",
  "parent_type_id" int4,
  "is_hidden" bool,
  "restricted" bool,
  "trimmable" bool
)
;

-- ----------------------------
-- Table structure for available_training_days
-- ----------------------------
DROP TABLE IF EXISTS "public"."available_training_days";
CREATE TABLE "public"."available_training_days" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "training_day" varchar(20) COLLATE "pg_catalog"."default",
  "is_preferred_long_day" bool DEFAULT false,
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for event_types
-- ----------------------------
DROP TABLE IF EXISTS "public"."event_types";
CREATE TABLE "public"."event_types" (
  "type_id" int4 NOT NULL,
  "type_key" varchar(50) COLLATE "pg_catalog"."default",
  "sort_order" int4
)
;

-- ----------------------------
-- Table structure for heart_rate_details
-- ----------------------------
DROP TABLE IF EXISTS "public"."heart_rate_details";
CREATE TABLE "public"."heart_rate_details" (
  "heartrate" int2,
  "timestamp" int8,
  "user_profile_pk" int8 NOT NULL,
  "calendar_date" date NOT NULL
)
;

-- ----------------------------
-- Table structure for heart_rates
-- ----------------------------
DROP TABLE IF EXISTS "public"."heart_rates";
CREATE TABLE "public"."heart_rates" (
  "user_profile_pk" int8 NOT NULL,
  "calendar_date" date,
  "start_timestamp_gmt" timestamp(6),
  "end_timestamp_gmt" timestamp(6),
  "start_timestamp_local" timestamp(6),
  "end_timestamp_local" timestamp(6),
  "max_heart_rate" int4,
  "min_heart_rate" int4,
  "resting_heart_rate" int4,
  "last_seven_days_avg_resting_heart_rate" int4,
  "created_at" timestamp(6),
  "updated_at" timestamp(6),
  "id" uuid NOT NULL DEFAULT uuid_generate_v4()
)
;

-- ----------------------------
-- Table structure for hydration_containers
-- ----------------------------
DROP TABLE IF EXISTS "public"."hydration_containers";
CREATE TABLE "public"."hydration_containers" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "name" varchar(100) COLLATE "pg_catalog"."default",
  "volume" int4,
  "unit" varchar(20) COLLATE "pg_catalog"."default",
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for privacy_settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."privacy_settings";
CREATE TABLE "public"."privacy_settings" (
  "type_id" int4 NOT NULL,
  "type_key" varchar(50) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for split_summaries
-- ----------------------------
DROP TABLE IF EXISTS "public"."split_summaries";
CREATE TABLE "public"."split_summaries" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "activity_id" int8,
  "no_of_splits" int4,
  "total_ascent" float8,
  "duration" float8,
  "split_type" varchar(50) COLLATE "pg_catalog"."default",
  "num_climb_sends" int4,
  "max_elevation_gain" float8,
  "average_elevation_gain" float8,
  "max_distance" int4,
  "distance" float8,
  "average_speed" float8,
  "max_speed" float8,
  "num_falls" int4,
  "elevation_loss" float8
)
;

-- ----------------------------
-- Table structure for step_details
-- ----------------------------
DROP TABLE IF EXISTS "public"."step_details";
CREATE TABLE "public"."step_details" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "start_gmt" timestamp(6),
  "end_gmt" timestamp(6),
  "steps" int4,
  "pushes" int4,
  "primary_activity_level" varchar(50) COLLATE "pg_catalog"."default",
  "activity_level_constant" bool,
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for user_data
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_data";
CREATE TABLE "public"."user_data" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "gender" varchar(10) COLLATE "pg_catalog"."default",
  "time_format" varchar(50) COLLATE "pg_catalog"."default",
  "birth_date" date,
  "measurement_system" varchar(20) COLLATE "pg_catalog"."default",
  "activity_level" varchar(50) COLLATE "pg_catalog"."default",
  "handedness" varchar(10) COLLATE "pg_catalog"."default",
  "intensity_minutes_calc_method" varchar(20) COLLATE "pg_catalog"."default",
  "moderate_intensity_minutes_hr_zone" int4,
  "vigorous_intensity_minutes_hr_zone" int4,
  "hydration_measurement_unit" varchar(20) COLLATE "pg_catalog"."default",
  "hydration_auto_goal_enabled" bool,
  "firstbeat_max_stress_score" numeric(8,2),
  "firstbeat_cycling_lt_timestamp" timestamp(6),
  "firstbeat_running_lt_timestamp" timestamp(6),
  "threshold_heart_rate_auto_detected" bool,
  "ftp_auto_detected" bool,
  "training_status_paused_date" date,
  "golf_distance_unit" varchar(20) COLLATE "pg_catalog"."default",
  "golf_elevation_unit" varchar(20) COLLATE "pg_catalog"."default",
  "golf_speed_unit" varchar(20) COLLATE "pg_catalog"."default",
  "external_bottom_time" timestamp(6),
  "virtual_caddie_data_source" varchar(50) COLLATE "pg_catalog"."default",
  "number_dives_automatically" bool,
  "dive_number" int4,
  "lactate_threshold_speed" numeric(8,4),
  "lactate_threshold_heart_rate" numeric(5,2),
  "weight" float4,
  "height" float4,
  "vo2_max_running" numeric(5,2),
  "vo2_max_cycling" numeric(5,2),
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for user_heart_rate_formats
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_heart_rate_formats";
CREATE TABLE "public"."user_heart_rate_formats" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "format_id" int4,
  "format_key" varchar(50) COLLATE "pg_catalog"."default",
  "min_fraction" int4,
  "max_fraction" int4,
  "grouping_used" bool,
  "display_format" text COLLATE "pg_catalog"."default",
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for user_power_formats
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_power_formats";
CREATE TABLE "public"."user_power_formats" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "format_id" int4,
  "format_key" varchar(50) COLLATE "pg_catalog"."default",
  "min_fraction" int4,
  "max_fraction" int4,
  "grouping_used" bool,
  "display_format" text COLLATE "pg_catalog"."default",
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_roles";
CREATE TABLE "public"."user_roles" (
  "activity_id" int8,
  "role_name" varchar(100) COLLATE "pg_catalog"."default",
  "id" uuid NOT NULL DEFAULT uuid_generate_v4()
)
;

-- ----------------------------
-- Table structure for user_settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_settings";
CREATE TABLE "public"."user_settings" (
  "id" int8 NOT NULL,
  "connect_date" timestamp(6),
  "source_type" varchar(50) COLLATE "pg_catalog"."default",
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for user_sleep
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_sleep";
CREATE TABLE "public"."user_sleep" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "sleep_time" int4,
  "default_sleep_time" bool,
  "wake_time" int4,
  "default_wake_time" bool,
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for user_sleep_windows
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_sleep_windows";
CREATE TABLE "public"."user_sleep_windows" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "sleep_window_frequency" varchar(20) COLLATE "pg_catalog"."default",
  "start_sleep_time_seconds_from_midnight" int4,
  "end_sleep_time_seconds_from_midnight" int4,
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for user_available_training_days
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_available_training_days";
CREATE TABLE "public"."user_available_training_days" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "days" VARCHAR(50)[],
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

ALTER TABLE "public"."user_available_training_days" 
  ADD CONSTRAINT "user_available_training_days_unique_fields" UNIQUE ("user_profile_pk");
  
-- ----------------------------
-- Table structure for user_preferred_long_training_days
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_preferred_long_training_days";
CREATE TABLE "public"."user_preferred_long_training_days" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "days" VARCHAR(50)[],
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

ALTER TABLE "public"."user_preferred_long_training_days" 
  ADD CONSTRAINT "user_preferred_long_training_days_unique_fields" UNIQUE ("user_profile_pk");

-- ----------------------------
-- Table structure for weather_locations
-- ----------------------------
DROP TABLE IF EXISTS "public"."weather_locations";
CREATE TABLE "public"."weather_locations" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user_profile_pk" int8,
  "use_fixed_location" bool,
  "latitude" numeric(15,12),
  "longitude" numeric(15,12),
  "location_name" varchar(100) COLLATE "pg_catalog"."default",
  "iso_country_code" varchar(10) COLLATE "pg_catalog"."default",
  "postal_code" varchar(20) COLLATE "pg_catalog"."default",
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Function structure for uuid_generate_v1
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v1"();
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v1"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v1'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v1mc
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v1mc"();
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v1mc"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v1mc'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v3
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v3"("namespace" uuid, "name" text);
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v3"("namespace" uuid, "name" text)
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v3'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v4
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v4"();
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v4"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v4'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v5
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v5"("namespace" uuid, "name" text);
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v5"("namespace" uuid, "name" text)
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v5'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_nil
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_nil"();
CREATE OR REPLACE FUNCTION "public"."uuid_nil"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_nil'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_dns
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_dns"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_dns"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_dns'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_oid
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_oid"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_oid"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_oid'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_url
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_url"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_url"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_url'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_x500
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_x500"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_x500"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_x500'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."activity_detail_metrics_id_seq"
OWNED BY "public"."activity_detail_metrics"."id";
SELECT setval('"public"."activity_detail_metrics_id_seq"', 1, false);

-- ----------------------------
-- Primary Key structure for table activities
-- ----------------------------
ALTER TABLE "public"."activities" ADD CONSTRAINT "activities_pkey" PRIMARY KEY ("activity_id");

-- ----------------------------
-- Uniques structure for table activity_detail_metrics
-- ----------------------------
ALTER TABLE "public"."activity_detail_metrics" ADD CONSTRAINT "activity_detail_metrics_activity_id_measurement_index_key" UNIQUE ("activity_id", "measurement_index");

-- ----------------------------
-- Primary Key structure for table activity_detail_metrics
-- ----------------------------
ALTER TABLE "public"."activity_detail_metrics" ADD CONSTRAINT "activity_detail_metrics_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table activity_splits
-- ----------------------------
ALTER TABLE "public"."activity_splits" ADD CONSTRAINT "activity_splits_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table activity_types
-- ----------------------------
ALTER TABLE "public"."activity_types" ADD CONSTRAINT "activity_types_pkey" PRIMARY KEY ("type_id");

-- ----------------------------
-- Primary Key structure for table available_training_days
-- ----------------------------
ALTER TABLE "public"."available_training_days" ADD CONSTRAINT "available_training_days_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table event_types
-- ----------------------------
ALTER TABLE "public"."event_types" ADD CONSTRAINT "event_types_pkey" PRIMARY KEY ("type_id");

-- ----------------------------
-- Primary Key structure for table heart_rates
-- ----------------------------
ALTER TABLE "public"."heart_rates" ADD CONSTRAINT "heart_rates_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table hydration_containers
-- ----------------------------
ALTER TABLE "public"."hydration_containers" ADD CONSTRAINT "hydration_containers_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table privacy_settings
-- ----------------------------
ALTER TABLE "public"."privacy_settings" ADD CONSTRAINT "privacy_settings_pkey" PRIMARY KEY ("type_id");

-- ----------------------------
-- Primary Key structure for table split_summaries
-- ----------------------------
ALTER TABLE "public"."split_summaries" ADD CONSTRAINT "split_summaries_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table step_details
-- ----------------------------
ALTER TABLE "public"."step_details" ADD CONSTRAINT "step_details_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table user_data
-- ----------------------------
ALTER TABLE "public"."user_data" ADD CONSTRAINT "user_data_fk_unique" UNIQUE ("user_profile_pk");

-- ----------------------------
-- Primary Key structure for table user_data
-- ----------------------------
ALTER TABLE "public"."user_data" ADD CONSTRAINT "user_data_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table user_heart_rate_formats
-- ----------------------------
CREATE INDEX "idx_heart_rate_formats_user_profile_pk" ON "public"."user_heart_rate_formats" USING btree (
  "user_profile_pk" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table user_heart_rate_formats
-- ----------------------------
ALTER TABLE "public"."user_heart_rate_formats" ADD CONSTRAINT "unique_user_profile_pk_heart_rate_format" UNIQUE ("user_profile_pk");

-- ----------------------------
-- Primary Key structure for table user_heart_rate_formats
-- ----------------------------
ALTER TABLE "public"."user_heart_rate_formats" ADD CONSTRAINT "user_heart_rate_formats_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table user_power_formats
-- ----------------------------
CREATE INDEX "idx_user_profile_pk" ON "public"."user_power_formats" USING btree (
  "user_profile_pk" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table user_power_formats
-- ----------------------------
ALTER TABLE "public"."user_power_formats" ADD CONSTRAINT "unique_user_profile_pk_power_format" UNIQUE ("user_profile_pk");

-- ----------------------------
-- Primary Key structure for table user_power_formats
-- ----------------------------
ALTER TABLE "public"."user_power_formats" ADD CONSTRAINT "user_power_formats_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table user_roles
-- ----------------------------
ALTER TABLE "public"."user_roles" ADD CONSTRAINT "user_roles_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table user_settings
-- ----------------------------
ALTER TABLE "public"."user_settings" ADD CONSTRAINT "user_profiles_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table user_sleep
-- ----------------------------
ALTER TABLE "public"."user_sleep" ADD CONSTRAINT "id_user_profile_pk" UNIQUE ("user_profile_pk");

-- ----------------------------
-- Primary Key structure for table user_sleep
-- ----------------------------
ALTER TABLE "public"."user_sleep" ADD CONSTRAINT "user_sleep_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table user_sleep_windows
-- ----------------------------
ALTER TABLE "public"."user_sleep_windows" ADD CONSTRAINT "id_user_profile_pk_sleep_windows_freq" UNIQUE ("user_profile_pk", "sleep_window_frequency");

-- ----------------------------
-- Primary Key structure for table user_sleep_windows
-- ----------------------------
ALTER TABLE "public"."user_sleep_windows" ADD CONSTRAINT "user_sleep_windows_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table weather_locations
-- ----------------------------
ALTER TABLE "public"."weather_locations" ADD CONSTRAINT "weather_locations_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table activities
-- ----------------------------
ALTER TABLE "public"."activities" ADD CONSTRAINT "activities_activity_type_id_fkey" FOREIGN KEY ("activity_type_id") REFERENCES "public"."activity_types" ("type_id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."activities" ADD CONSTRAINT "activities_event_type_id_fkey" FOREIGN KEY ("event_type_id") REFERENCES "public"."event_types" ("type_id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."activities" ADD CONSTRAINT "activities_privacy_type_id_fkey" FOREIGN KEY ("privacy_type_id") REFERENCES "public"."privacy_settings" ("type_id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table activity_detail_metrics
-- ----------------------------
ALTER TABLE "public"."activity_detail_metrics" ADD CONSTRAINT "activity_detail_metrics_activity_id_fkey" FOREIGN KEY ("activity_id") REFERENCES "public"."activities" ("activity_id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table activity_splits
-- ----------------------------
ALTER TABLE "public"."activity_splits" ADD CONSTRAINT "activity_splits_activity_id_fkey" FOREIGN KEY ("activity_id") REFERENCES "public"."activities" ("activity_id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table available_training_days
-- ----------------------------
ALTER TABLE "public"."available_training_days" ADD CONSTRAINT "available_training_days_user_profile_pk_fkey" FOREIGN KEY ("user_profile_pk") REFERENCES "public"."user_settings" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table hydration_containers
-- ----------------------------
ALTER TABLE "public"."hydration_containers" ADD CONSTRAINT "hydration_containers_user_profile_pk_fkey" FOREIGN KEY ("user_profile_pk") REFERENCES "public"."user_settings" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table split_summaries
-- ----------------------------
ALTER TABLE "public"."split_summaries" ADD CONSTRAINT "split_summaries_activity_id_fkey" FOREIGN KEY ("activity_id") REFERENCES "public"."activities" ("activity_id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table step_details
-- ----------------------------
ALTER TABLE "public"."step_details" ADD CONSTRAINT "step_details_user_profile_pk_fkey" FOREIGN KEY ("user_profile_pk") REFERENCES "public"."user_settings" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table user_data
-- ----------------------------
ALTER TABLE "public"."user_data" ADD CONSTRAINT "user_data_user_setting_pk_fkey" FOREIGN KEY ("user_profile_pk") REFERENCES "public"."user_settings" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table user_heart_rate_formats
-- ----------------------------
ALTER TABLE "public"."user_heart_rate_formats" ADD CONSTRAINT "user_heart_rate_formats_user_profile_pk_fkey" FOREIGN KEY ("user_profile_pk") REFERENCES "public"."user_settings" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table user_power_formats
-- ----------------------------
ALTER TABLE "public"."user_power_formats" ADD CONSTRAINT "user_power_formats_user_profile_pk_fkey" FOREIGN KEY ("user_profile_pk") REFERENCES "public"."user_settings" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table user_roles
-- ----------------------------
ALTER TABLE "public"."user_roles" ADD CONSTRAINT "user_roles_activity_id_fkey" FOREIGN KEY ("activity_id") REFERENCES "public"."activities" ("activity_id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table user_sleep
-- ----------------------------
ALTER TABLE "public"."user_sleep" ADD CONSTRAINT "user_sleep_user_settings_pk_fkey" FOREIGN KEY ("user_profile_pk") REFERENCES "public"."user_settings" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table user_sleep_windows
-- ----------------------------
ALTER TABLE "public"."user_sleep_windows" ADD CONSTRAINT "user_sleep_windows_user_profile_pk_fkey" FOREIGN KEY ("user_profile_pk") REFERENCES "public"."user_settings" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table weather_locations
-- ----------------------------
ALTER TABLE "public"."weather_locations" ADD CONSTRAINT "weather_locations_user_profile_pk_fkey" FOREIGN KEY ("user_profile_pk") REFERENCES "public"."user_settings" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
