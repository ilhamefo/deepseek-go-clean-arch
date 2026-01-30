package metric

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Request metrics
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active HTTP connections",
		},
	)

	// Database metrics
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"database", "operation", "table"},
	)

	DatabaseConnectionsActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
		[]string{"database"},
	)

	DatabaseConnectionsIdle = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections_idle",
			Help: "Number of idle database connections",
		},
		[]string{"database"},
	)

	// Excel generation metrics
	ExcelGenerationDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "excel_generation_duration_seconds",
			Help:    "Excel file generation duration",
			Buckets: []float64{1, 5, 10, 30, 60, 120, 300, 600},
		},
	)

	ExcelRowsProcessed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "excel_rows_processed_total",
			Help: "Total number of Excel rows processed",
		},
	)

	ExcelFilesGenerated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "excel_files_generated_total",
			Help: "Total number of Excel files generated",
		},
	)

	// Redis metrics
	RedisCacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_hits_total",
			Help: "Total number of Redis cache hits",
		},
	)

	RedisCacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_misses_total",
			Help: "Total number of Redis cache misses",
		},
	)

	RedisOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_operation_duration_seconds",
			Help:    "Redis operation duration",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1},
		},
		[]string{"operation"},
	)

	// Auth metrics
	AuthLoginAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_login_attempts_total",
			Help: "Total number of login attempts",
		},
		[]string{"method", "status"}, // method: google_oauth, jwt; status: success, failure
	)

	ActiveSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sessions",
			Help: "Number of active user sessions",
		},
	)

	JWTTokensIssued = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "jwt_tokens_issued_total",
			Help: "Total number of JWT tokens issued",
		},
	)

	// Garmin sync metrics
	GarminSyncDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "garmin_sync_duration_seconds",
			Help:    "Garmin data sync duration",
			Buckets: []float64{1, 5, 10, 30, 60, 120},
		},
	)

	GarminSyncTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "garmin_sync_total",
			Help: "Total number of Garmin sync operations",
		},
		[]string{"status"}, // success, failure
	)
)

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	status := fmt.Sprintf("%dxx", statusCode/100)

	HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration.Seconds())
	HTTPRequestTotal.WithLabelValues(method, path, status).Inc()
}

// RecordDatabaseQuery records database query metrics
func RecordDatabaseQuery(database, operation, table string, duration time.Duration) {
	DatabaseQueryDuration.WithLabelValues(database, operation, table).Observe(duration.Seconds())
}

// RecordExcelGeneration records Excel generation metrics
func RecordExcelGeneration(duration time.Duration, rowCount int) {
	ExcelGenerationDuration.Observe(duration.Seconds())
	ExcelRowsProcessed.Add(float64(rowCount))
	ExcelFilesGenerated.Inc()
}

// RecordRedisOperation records Redis operation metrics
func RecordRedisOperation(operation string, duration time.Duration, hit bool) {
	RedisOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())

	if hit {
		RedisCacheHits.Inc()
	} else {
		RedisCacheMisses.Inc()
	}
}

// RecordLoginAttempt records authentication attempts
func RecordLoginAttempt(method string, success bool) {
	status := "failure"
	if success {
		status = "success"
	}
	AuthLoginAttempts.WithLabelValues(method, status).Inc()

	if success {
		JWTTokensIssued.Inc()
	}
}

// RecordGarminSync records Garmin sync operations
func RecordGarminSync(duration time.Duration, success bool) {
	GarminSyncDuration.Observe(duration.Seconds())

	status := "failure"
	if success {
		status = "success"
	}
	GarminSyncTotal.WithLabelValues(status).Inc()
}
